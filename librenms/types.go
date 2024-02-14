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
