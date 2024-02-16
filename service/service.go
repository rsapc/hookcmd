package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"

	"golang.org/x/exp/slog"

	"github.com/rsapc/hookcmd/librenms"
	"github.com/rsapc/hookcmd/models"
	"github.com/rsapc/hookcmd/netbox"
)

var ErrUnimplemented = errors.New("method has not been implemented")

type Service struct {
	getenv   func(string) string
	logger   models.Logger
	netbox   *netbox.Client
	librenms *librenms.Client
}

// NewService creates a new instance of the service.
//
//	getenv: a function to return envvars.  If nil
//	        gets set to os.GetEnv
func NewService(getenv func(string) string, logger models.Logger) *Service {
	s := &Service{}
	if getenv == nil {
		s.getenv = os.Getenv
	} else {
		s.getenv = getenv
	}
	if logger == nil {
		s.logger = slog.Default()
	} else {
		s.logger = logger
	}
	s.netbox = netbox.NewClient(s.getenv("NETBOX_URL"), s.getenv("NETBOX_TOKEN"), s.logger)
	s.librenms = librenms.NewClient(s.getenv("LIBRENMS_URL"), s.getenv("LIBRENMS_TOKEN"), s.logger)
	return s
}

func (s *Service) IPdnsUpdate(addr string) error {
	ip := netbox.IPfromCIDR(addr)
	addrs, err := net.LookupAddr(ip)
	if err != nil {
		slog.Error("Could not find address", "err", err)
		os.Exit(1)
	}

	if len(addrs) > 0 {
		err = s.netbox.SetIPDNS(ip, addrs[0])
		if err != nil {
			s.logger.Error(fmt.Sprintf("Failed to update Netbox IPAddress record: %v", err))
			return fmt.Errorf("failed to update Netbox IPAddress record")
		}
	}
	return nil
}

// AddToLibreNMS adds the IP to libre and updates Netbox
func (s *Service) AddToLibreNMS(addr string, model string, modelID int64) error {
	ip := netbox.IPfromCIDR(addr)
	devid, err := s.librenms.AddDevice(ip)
	if err != nil {
		s.netbox.AddJournalEntry(model, modelID, netbox.WarningLevel, err.Error())
		return err
	}
	if err = s.netbox.AddJournalEntry(model, modelID, netbox.InfoLevel, fmt.Sprintf("added device to LibreNMS.  id=%d", devid)); err != nil {
		s.logger.Error(fmt.Sprintf("could not add journal entry: %v", err), "service", "service")
	}
	return s.netbox.SetMonitoringID(model, modelID, devid)
}

// DeviceDown will set the device status in Netbox based on the
// state of the alert.  Payload is expected to be the JSON of
// the alert.
func (s *Service) DeviceDown(payload string) error {
	var alert librenms.LibreAlert
	err := json.Unmarshal(([]byte)(payload), &alert)
	if err != nil {
		s.logger.Error(fmt.Sprintf("could not decode alert payload: %v", err), "service", "service")
		return err
	}

	objectType, objectID, err := s.netbox.FindMonitoredObject(alert.DeviceID)
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}

	var journalEntry = fmt.Sprintf(`%s
	
	%s status updated as of %s

	%s`, alert.Subject, alert.SysName, alert.Timestamp, alert.Runbook)

	data := make(map[string]interface{})
	data["status"] = "offline"

	switch alert.State {
	case librenms.AlertFiring:
		// if this is the first occurance of the alert
		if alert.ID == alert.UID {
			err = s.netbox.UpdateObject(objectType, objectID, data)
			if err != nil {
				return err
			}
			return s.netbox.AddJournalEntry(objectType, objectID, netbox.DangerLevel, journalEntry)
		}
	case librenms.AlertCleared:
		data["status"] = "active"
		err = s.netbox.UpdateObject(objectType, objectID, data)
		if err != nil {
			return err
		}
		return s.netbox.AddJournalEntry(objectType, objectID, netbox.SuccessLevel, journalEntry)
	}
	return nil
}

func (s *Service) GetDeviceInfo(deviceID int) error {
	netboxType, netboxID, err := s.netbox.FindMonitoredObject(deviceID)
	if err != nil {
		s.logger.Error("could not find netbox device", "device_id", deviceID, "err", err)
		return err
	}
	device, err := s.librenms.GetDevice(deviceID)
	if err != nil {
		return err
	}
	return s.updateDeviceInfo(device, netboxType, netboxID)
}

// updateDeviceInfo takes a Device object from LibreNMS and updates the corresponding
// fields in Netbox.  It will also make Journal Entries in Netbox with what has been
// done.
func (s *Service) updateDeviceInfo(device librenms.LibreDevice, netboxType string, netboxID int64) error {
	nbdev, err := s.netbox.GetDeviceOrVMbyType(netboxType, netboxID)
	if err != nil {
		return err
	}
	data, err := s.updateNetboxDevice(device, nbdev)
	if err != nil {
		s.netbox.AddJournalEntry(netboxType, netboxID, netbox.WarningLevel, "could not update device:\n\n%s", err.Error())
		return err
	}
	s.netbox.AddJournalEntry(netboxType, netboxID, netbox.SuccessLevel, "device updated with values from LibreNMS\n\nUpdate Data:\n%s", data)
	s.logger.Info("successfully updated device from LibreNMS", "deviceType", netboxType, "ID", netboxID)
	return nil
}

func (s *Service) updateNetboxDevice(device librenms.LibreDevice, nbdev netbox.DeviceOrVM) (string, error) {
	cf := make(map[string]interface{})
	data := make(map[string]interface{})
	if nbdev.CustomFields.MonitoringID == nil {
		cf["monitoring_id"] = device.DeviceID
		data["custom_fields"] = cf
	}

	if nbdev.Description == "" {
		var description *string
		if device.SysDescr != nil {
			description = device.SysDescr
		}
		if device.Purpose != nil {
			description = device.Purpose
		} else if device.Hardware != nil {
			description = device.Hardware
		}
		if description != nil {
			data["description"] = *description
		}
	}

	if device.Serial != nil && *device.Serial != "" {
		data["serial"] = *device.Serial
	}

	if device.Lat != nil {
		data["latitude"] = *device.Lat
	}
	if device.Lng != nil {
		data["longitude"] = *device.Lng
	}

	d, _ := json.Marshal(data)
	return string(d), s.netbox.UpdateObjectByURL(nbdev.URL, data)
}

// FindDevice searches for a device by IP
func (s *Service) FindDevice(addr string) error {
	ip := netbox.IPfromCIDR(addr)
	ipInfo, err := s.netbox.SearchIP(ip)
	if err != nil {
		return err
	}
	if ipInfo.Count != 1 {
		s.logger.Error(fmt.Sprintf("invalid number of netbox IPs found: %d", ipInfo.Count))
		return fmt.Errorf("invalid netbox device count: %d", ipInfo.Count)
	}
	device, err := s.librenms.GetDeviceByIP(ip)
	if err != nil {
		return err
	}
	nbdev, err := s.netbox.GetDeviceOrVM(ipInfo.Results[0].URL)
	if err != nil {
		return err
	}
	_, err = s.updateNetboxDevice(device, nbdev)
	if err != nil {
		errMsg := fmt.Sprintf("could not update Netbox device: %v", err)
		s.logger.Error(errMsg, "url", nbdev.URL)
		return err
	}
	return err
}

// MissingFromLibre generates a report of Netbox devices that are not
// in LibreNMS
func (s *Service) MissingFromLibre(out io.Writer) error {
	devices, err := s.netbox.SearchDeviceAndVM(
		"status=active",
		"has_primary_ip=true",
		"cf_monitoring_id__lte=0")
	if err != nil {
		s.logger.Error("could not get list of netbox devices", "err", err)
		return err
	}
	io.WriteString(out, "Name,IP\n")
	for _, device := range devices {
		port, err := s.librenms.FindPortForIP(device.PrimaryIP.Address)
		if err != nil {
			if errors.Is(err, librenms.ErrNotFound) {
				io.WriteString(out, fmt.Sprintf("%s,%s\n", device.Name, netbox.IPfromCIDR(device.PrimaryIP.Address)))
				continue
			}
			return err
		}
		libreDev, err := s.librenms.GetDevice(port.DeviceID)
		if err != nil {
			return err
		}
		s.updateNetboxDevice(libreDev, device)
	}
	return nil
}
