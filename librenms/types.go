package librenms

type AddDeviceResponse struct {
	Devices []struct {
		DeviceID int         `json:"device_id"`
		Hostname string      `json:"hostname"`
		Serial   interface{} `json:"serial"`
	} `json:"devices"`
	Message string `json:"message"`
	Status  string `json:"status"`
}
