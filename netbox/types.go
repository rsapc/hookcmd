package netbox

import "fmt"

type JournalLevel int

// JournalLevels
const (
	Undefined JournalLevel = iota
	InfoLevel
	SuccessLevel
	DangerLevel
	WarningLevel
)

type MonitoredObject struct {
	ID         int64  `json:"id"`
	URL        string `json:"url"`
	ObjectType string `json:"-"`
}

type MonitoringSearchResults struct {
	Count    int               `json:"count"`
	Next     interface{}       `json:"next"`
	Previous interface{}       `json:"previous"`
	Results  []MonitoredObject `json:"results"`
}

// getObjectType returns the full netbox object type for the given model.
// For example, given the type of "device" will return "dcim.device"
func getObjectType(aModel string) string {
	var group string
	switch aModel {
	case "location":
		fallthrough
	case "device":
		group = "dcim"
	case "virtualmachine":
		group = "virtualization"
	case "ipaddress":
		group = "ipam"
	default:
		return "Invalid"
	}
	return fmt.Sprintf("%s.%s", group, aModel)
}

func getJournalLevel(level JournalLevel) string {
	switch level {
	case Undefined:
		return ""
	case InfoLevel:
		return "info"
	case SuccessLevel:
		return "success"
	case WarningLevel:
		return "warning"
	case DangerLevel:
		return "danger"
	}
	return ""
}

type DeviceOrVM struct {
	AssetTag     *string `json:"asset_tag"`
	Comments     string  `json:"comments"`
	Created      string  `json:"created"`
	CustomFields struct {
		MonitoringID *int `json:"monitoring_id"`
	} `json:"custom_fields"`
	Description string        `json:"description"`
	DeviceRole  DisplayIDName `json:"device_role"`
	Display     string        `json:"display"`
	ID          int           `json:"id"`
	LastUpdated string        `json:"last_updated"`
	Latitude    *float64      `json:"latitude"`
	Longitude   *float64      `json:"longitude"`
	Name        string        `json:"name"`
	PrimaryIP   PrimaryI      `json:"primary_ip"`
	PrimaryIp4  PrimaryI      `json:"primary_ip4"`
	Rack        struct {
		Display string `json:"display"`
		ID      int    `json:"id"`
		Name    string `json:"name"`
		URL     string `json:"url"`
	} `json:"rack"`
	Role   DisplayIDName `json:"role"`
	Serial string        `json:"serial"`
	Site   DisplayIDName `json:"site"`
	Status LabelValue    `json:"status"`
	URL    string        `json:"url"`
}
type DisplayIDName struct {
	Display string `json:"display"`
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Slug    string `json:"slug"`
	URL     string `json:"url"`
}
type LabelValue struct {
	Label string `json:"label"`
	Value string `json:"value"`
}
type PrimaryI struct {
	Address string `json:"address"`
	Display string `json:"display"`
	Family  int    `json:"family"`
	ID      int    `json:"id"`
	URL     string `json:"url"`
}
