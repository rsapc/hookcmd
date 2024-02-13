package librenms

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/rsapc/webhooks/models"
	"golang.org/x/exp/slog"
)

type Client struct {
	client  *resty.Client
	log     models.Logger
	baseURL string
	token   string
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

func (c *Client) buildURL(path string) string {
	return fmt.Sprintf("%s%s", c.baseURL, path)
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
		errMsg := fmt.Sprintf("invalid response from server: %d.  %s", resp.StatusCode(), obj.Message)
		c.log.Error(errMsg, "url", r.URL)
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
