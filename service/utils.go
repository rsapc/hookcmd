package service

import (
	"github.com/rsapc/hookcmd/librenms"
	"github.com/rsapc/hookcmd/netbox"
)

// GetINterfaceTypeFromIfType takes an SNMP ifType value
// and returns a corresponding Netbox interface type (or
// an approximate equivalent)
func GetInterfaceTypeFromIfType(ifType string) string {
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
	return nbType
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
		update = update || ifUpd.SetSpeed(port.GetSpeed()/1000)
	}
	if intf.GetDuplex() != port.GetDuplex() {
		update = update || ifUpd.SetDuplex(port.IfDuplex)
	}
	if intf.GetMacAddress() != port.GetPhysAddress() {
		update = update || ifUpd.SetMac(port.IfPhysAddress)
	}

	return ifUpd, update
}
