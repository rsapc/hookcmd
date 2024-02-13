package netbox

import "fmt"

type model string
type JournalLevel int

// JournalLevels
const (
	Undefined JournalLevel = iota
	InfoLevel
	SuccessLevel
	DangerLevel
	WarningLevel
)

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
