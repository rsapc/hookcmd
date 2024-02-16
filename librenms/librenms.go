package librenms

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/go-resty/resty/v2"
	"github.com/rsapc/hookcmd/models"
	"golang.org/x/exp/slog"
)

var ErrNotFound = errors.New("the request object was not found")

type Client struct {
	client  *resty.Client
	log     models.Logger
	baseURL string
	token   string
	ipList  *[]IP
	mux     sync.Mutex
}

// NewClient creates a new LibreNMS API client
func NewClient(url string, token string, logger models.Logger) *Client {
	c := &Client{log: logger}
	c.client = resty.New()
	c.client.SetRedirectPolicy(resty.FlexibleRedirectPolicy(5))

	c.baseURL = fmt.Sprintf("%s/api/v0", url)
	c.token = token
	if log, ok := logger.(*slog.Logger); ok {
		c.log = log.With("service", "librenms")
	}

	return c
}

func (c *Client) buildRequest() *resty.Request {
	return c.client.NewRequest().SetHeader("X-Auth-Token", c.token)
}

func (c *Client) buildURL(path string, args ...any) string {
	urlPath := fmt.Sprintf(path, args...)
	return fmt.Sprintf("%s%s", c.baseURL, urlPath)
}

// AddDevice adds the given IP to LibreNMS to monitor.  Returns the
// device ID assigned in LibreNMS
func (c *Client) AddDevice(ip string) (deviceID int, err error) {
	obj := AddDeviceResponse{}
	r := c.buildRequest().SetResult(&obj)
	data := make(map[string]interface{})
	data["hostname"] = ip
	data["ping_fallback"] = true
	r.SetBody(data)
	resp, err := r.Post(c.buildURL("/devices"))
	if err != nil {
		c.log.Error("Could not add device to LibreNMS: %v", err)
		return deviceID, err
	}
	if resp.IsError() {
		json.Unmarshal(resp.Body(), &obj)
		errMsg := fmt.Sprintf("invalid response from server: %d.  %s", resp.StatusCode(), obj.Message)
		c.log.Error(errMsg, "url", r.URL, "body", string(resp.Body()))
		return deviceID, fmt.Errorf("%s", errMsg)
	}
	if obj.Status == "ok" {
		c.log.Info(obj.Message)
		return obj.Devices[0].DeviceID, nil
	} else {
		errMsg := fmt.Sprintf("Invalid status [%s]: %s", obj.Status, obj.Message)
		c.log.Warn(errMsg, "url", r.URL)
		return deviceID, fmt.Errorf("%s", errMsg)
	}
}

// GetDevice returns the device with the corresponding ID
func (c *Client) GetDevice(deviceID int) (LibreDevice, error) {
	var device LibreDevice
	obj := LibreDeviceResponse{}
	r := c.buildRequest().SetResult(&obj)
	resp, err := r.Get(c.buildURL("/devices/%d", deviceID))
	if err != nil {
		c.log.Error("error getting device", "url", r.URL, "err", err)
		return device, err
	}
	if resp.IsError() {
		if resp.StatusCode() == 404 {
			return device, ErrNotFound
		}
		errObj, _ := GetLibreError(resp)
		c.log.Error("error status returned", "url", r.URL, "err", errObj.Message)
		return device, fmt.Errorf("error status returned %d: %s", resp.StatusCode(), errObj.Message)
	}
	if obj.Count == 0 {
		return device, ErrNotFound
	}
	if obj.Count > 1 {
		msg := fmt.Sprintf("too many devices found for id %d", deviceID)
		c.log.Error(msg)
		return device, fmt.Errorf(msg)
	}
	device = obj.Devices[0]
	return device, nil
}

func (c *Client) GetIPs() (ipList []IP, err error) {
	obj := IPResponse{}
	r := c.buildRequest().SetResult(&obj)
	resp, err := r.Get(c.buildURL("/resources/ip/addresses"))
	if err != nil {
		c.log.Error("error getting IP list", "url", r.URL, "err", err)
		return ipList, err
	}
	if resp.IsError() {
		if resp.StatusCode() == 404 {
			return ipList, ErrNotFound
		}
		errObj, _ := GetLibreError(resp)
		c.log.Error("error status returned", "url", r.URL, "err", errObj.Message)
		return ipList, fmt.Errorf("error status returned %d: %s", resp.StatusCode(), errObj.Message)
	}
	ipList = obj.IPAddresses
	return ipList, err
}

// GetPort returns the individual port request by ID
func (c *Client) GetPort(portID int) (port Port, err error) {
	obj := PortResponse{}
	r := c.buildRequest().SetResult(&obj)
	resp, err := r.Get(c.buildURL("/ports/%d", portID))
	if err != nil {
		c.log.Error("error getting port", "url", r.URL, "err", err)
		return port, err
	}
	if resp.IsError() {
		if resp.StatusCode() == 404 {
			return port, ErrNotFound
		}
		errObj, _ := GetLibreError(resp)
		c.log.Error("error status returned", "url", r.URL, "err", errObj.Message)
		return port, fmt.Errorf("error status returned %d: %s", resp.StatusCode(), errObj.Message)
	}
	if obj.Count == 0 {
		return port, ErrNotFound
	}
	if obj.Count > 1 {
		msg := fmt.Sprintf("too many ports found for id %d", portID)
		c.log.Error(msg)
		return port, fmt.Errorf(msg)
	}
	port = obj.Ports[0]
	return port, err
}

func (c *Client) LoadIPs() error {
	if c.ipList == nil {
		c.mux.Lock()
		ipList, err := c.GetIPs()
		if err != nil {
			c.log.Error("could not get IP list", "err", err)
			c.mux.Unlock()
			return err
		}
		c.ipList = &ipList
		c.mux.Unlock()
	}
	return nil
}

func (c *Client) findPortForIP(ip string) (port Port, err error) {
	found := false
	if err := c.LoadIPs(); err != nil {
		return port, err
	}
	for _, address := range *c.ipList {
		if ip == address.Ipv4Address {
			port, err = c.GetPort(address.PortID)
			if err != nil {
				return port, err
			}
			found = true
			break
		}
	}
	if !found {
		err = ErrNotFound
	}
	return port, err
}

func (c *Client) GetDeviceByIP(ip string) (device LibreDevice, err error) {
	port, err := c.findPortForIP(ip)
	if err != nil {
		c.log.Error("could not find device for IP", "err", err)
		return device, err
	}
	return c.GetDevice(port.DeviceID)
}
