package librenms

const (
	AlertFiring  = 1
	AlertCleared = 0
)

type AddDeviceResponse struct {
	Devices []struct {
		DeviceID int         `json:"device_id"`
		Hostname string      `json:"hostname"`
		Serial   interface{} `json:"serial"`
	} `json:"devices"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

// LibreAlert represents an API alert from LibreNMS.
// Transport configuration should look like this:
/*
Transport Headers: Content-type=application/json
Transport Body: ```
{"device_id": {{ $device_id }},
"timestamp": "{{ $timestamp }}",
"subject": "{{ $title }}",
"host": "{{ $hostname }}",
"sysName": "{{ $sysName }}",
"location": "{{ $location }}",
"ip": "{{ $ip }}",
"state": {{ $state }},
"severity": "{{ $severity }}",
"id": "{{ $id }}",
"uid": "{{ $uid }}",
"runbook": "{{ $proc }}"
}
```
*/
type LibreAlert struct {
	DeviceID  int    `json:"device_id"`
	Host      string `json:"host"`
	ID        string `json:"id"`
	IP        string `json:"ip"`
	Location  string `json:"location"`
	Runbook   string `json:"runbook"`
	Severity  string `json:"severity"`
	State     int    `json:"state"`
	Subject   string `json:"subject"`
	SysName   string `json:"sysName"`
	Timestamp string `json:"timestamp"`
	UID       string `json:"uid"`
}
type LibreDevice struct {
	Community           *string     `json:"community"`
	DeviceID            int         `json:"device_id"`
	DisableNotify       int         `json:"disable_notify"`
	Disabled            int         `json:"disabled"`
	Display             *string     `json:"display"`
	Features            interface{} `json:"features"`
	Hardware            *string     `json:"hardware"`
	Hostname            string      `json:"hostname"`
	IP                  string      `json:"ip"`
	Ignore              int         `json:"ignore"`
	IgnoreStatus        int         `json:"ignore_status"`
	Lat                 *float32    `json:"lat"`
	Lng                 *float32    `json:"lng"`
	Location            string      `json:"location"`
	LocationID          int         `json:"location_id"`
	MaxDepth            int         `json:"max_depth"`
	Notes               *string     `json:"notes"`
	Os                  string      `json:"os"`
	PollerGroup         int         `json:"poller_group"`
	Port                int         `json:"port"`
	PortAssociationMode int         `json:"port_association_mode"`
	Purpose             *string     `json:"purpose"`
	Serial              *string     `json:"serial"`
	Status              bool        `json:"status"`
	StatusReason        string      `json:"status_reason"`
	SysContact          *string     `json:"sysContact"`
	SysDescr            *string     `json:"sysDescr"`
	SysName             *string     `json:"sysName"`
	SysObjectID         *string     `json:"sysObjectID"`
	Type                *string     `json:"type"`
	Uptime              int         `json:"uptime"`
	Version             *string     `json:"version"`
}
type LibreDeviceResponse struct {
	Count   int           `json:"count"`
	Devices []LibreDevice `json:"devices"`
	Status  string        `json:"status"`
}

// LibreErrorResponse can be used to decode the
// body when a non-2xx status is returned
type LibreErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type IP struct {
	ContextName   string `json:"context_name"`
	Ipv4Address   string `json:"ipv4_address"`
	Ipv4AddressID int    `json:"ipv4_address_id"`
	Ipv4NetworkID string `json:"ipv4_network_id"`
	Ipv4Prefixlen int    `json:"ipv4_prefixlen"`
	PortID        int    `json:"port_id"`
}
type IPResponse struct {
	Count       int    `json:"count"`
	IPAddresses []IP   `json:"ip_addresses"`
	Status      string `json:"status"`
}

type Port struct {
	DeviceID           int     `json:"device_id"`
	IfAdminStatus      string  `json:"ifAdminStatus"`
	IfAdminStatusPrev  string  `json:"ifAdminStatus_prev"`
	IfAlias            string  `json:"ifAlias"`
	IfConnectorPresent string  `json:"ifConnectorPresent"`
	IfDescr            string  `json:"ifDescr"`
	IfDuplex           *string `json:"ifDuplex"`
	IfMtu              int     `json:"ifMtu"`
	IfName             string  `json:"ifName"`
	IfOperStatus       string  `json:"ifOperStatus"`
	IfOperStatusPrev   string  `json:"ifOperStatus_prev"`
	IfSpeed            *int    `json:"ifSpeed"`
	IfSpeedPrev        int     `json:"ifSpeed_prev"`
	IfType             string  `json:"ifType"`
	IfPhysAddress      *string `json:"ifPhysAddress"`
	IfVlan             *string `json:"ifVlan"`
	IfTrunk            *string `json:"ifTrunk"`
	PortDescrCircuit   *string `json:"port_descr_circuit"`
	PortDescrDescr     *string `json:"port_descr_descr"`
	PortDescrNotes     *string `json:"port_descr_notes"`
	PortDescrSpeed     *int    `json:"port_descr_speed"`
	PortDescrType      *string `json:"port_descr_type"`
	PortID             int     `json:"port_id"`
	PortName           *string `json:"portName"`
}

func (p Port) GetSpeed() int {
	if p.IfSpeed == nil {
		return 0
	}
	return *p.IfSpeed
}

func (p Port) GetPhysAddress() string {
	var mac string
	if p.IfPhysAddress != nil {
		mac = *p.IfPhysAddress
	}
	return mac
}

func (p Port) GetDuplex() string {
	if p.IfDuplex == nil {
		return "auto"
	}
	return *p.IfDuplex
}

// PortResponse is returned when querying for a given port
type PortResponse struct {
	Count  int    `json:"count"`
	Ports  []Port `json:"port"`
	Status string `json:"status"`
}

type PortSearchResponse struct {
	Ports  []Port `json:"ports"`
	Status string `json:"status"`
}
