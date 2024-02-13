package service

import (
	"fmt"
	"net"
	"os"

	"golang.org/x/exp/slog"

	"github.com/rsapc/webhooks/librenms"
	"github.com/rsapc/webhooks/models"
	"github.com/rsapc/webhooks/netbox"
)

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

	err = s.netbox.UpdateCustomFieldOnModel(model, modelID, "monitoring_id", devid)
	if err != nil {
		s.logger.Error(err.Error())
		s.netbox.AddJournalEntry(model, modelID, netbox.WarningLevel, fmt.Sprintf("failed to add monitoring_id: %d", devid))
		return err
	} else {
		msg := fmt.Sprintf("added monitoring_id %d to %s %d", devid, model, modelID)
		s.netbox.AddJournalEntry(model, modelID, netbox.SuccessLevel, msg)
	}
	return nil
}
