package service

import (
	"regexp"

	"github.com/rsapc/hookcmd/librenms"
	"github.com/rsapc/hookcmd/netbox"
)

var ifRegex = regexp.MustCompile(`^(?P<intf>[^\. ]+)\.?(?P<vlan>\d+)?`)
var parentIdx = ifRegex.SubexpIndex("intf")
var vlanIdx = ifRegex.SubexpIndex("vlan")

// GetINterfaceTypeFromIfType takes an SNMP ifType value
// and returns a corresponding Netbox interface type (or
// an approximate equivalent)
//
// If typ is virtual the parent interface will be returned if it can be determined
func GetInterfaceTypeFromIfType(ifType string, ifName string) (typ string, parent string) {
	var nbType string
	switch ifType {
	case "ethernetCsmacd":
		nbType = "1000base-t"
	case "ieee8023adLag":
	case "ds1":
		nbType = "t1"
	case "ds3":
		nbType = "t3"
	case "l2vlan":
		nbType = "virtual"
	case "l3ipvlan":
		nbType = "virtual"
	case "gpon":
		fallthrough
	case "aluGponOnu":
		fallthrough
	case "aluGponPhysicalUni":
		nbType = "xgs-pon"
	case "bridge":
		nbType = "bridge"
	case "other":
		nbType = "other"
	case "propVirtual":
		nbType = "virtual"
	case "softwareLoopback":
		fallthrough
	default:
		nbType = "other"
	}

	if ifRegex.MatchString(ifName) {
		matches := ifRegex.FindStringSubmatch(ifName)
		if matches[vlanIdx] != "" {
			parent = matches[parentIdx]
			nbType = "virtual"
		}
	}
	return nbType, parent
}

// GetUpdatedInterface compares a Netbox interface to a libreNMS port.  If there are changes the
// edit interface for Netbox is returned along with true to indicate changes should be made.
func GetUpdatedInterface(intf netbox.Interface, port librenms.Port) (*netbox.InterfaceEdit, bool) {
	update := false
	ifUpd := &netbox.InterfaceEdit{Device: &intf.Device.ID}
	if intf.Description != port.IfAlias {
		ifUpd.Description = port.IfAlias
		update = true
	}
	if intf.GetSpeed() != port.GetSpeed()/1000 {
		update = ifUpd.SetSpeed(port.GetSpeed()/1000) || update
	}
	if intf.GetDuplex() != port.GetDuplex() {
		update = ifUpd.SetDuplex(port.IfDuplex) || update
	}
	portAddr := port.GetPhysAddress()
	intfAddr := intf.GetMacAddress()
	if intfAddr != portAddr {
		update = ifUpd.SetMac(portAddr) || update
	}

	return ifUpd, update
}
