# webhooks
Command-line programs that function as netbox webhooks via https://github.com/adnanh/webhook 


## Commands

| Command | Params | Description |
| ------- | ------ | ----------- |
addLibreDevice | * {IP}<br/> * { netbox model }<br/> * { netbox model ID } | Adds the device/VM to LibreNMS.  The model and ID are used for posting the response back to Netbox.
ipdnsupdate | * {IP} | Does a DNS PTR lookup on the address and sets the `dns_name` field on the IP in Netbox.
devicedown | * {alert payload} | Sets the Netbox status to `Offline` when LibreNMS detects that it is down.
updatebyip |  * {IP}  (-x to return html) | Finds the IP in LibreNMS and updates the corresponding device in Netbox
updatedevice |  {monitoring_id} (-x to return hmtl)  |  Updates Netbox for the given LibreNMS ID
libreMissingReport | -o output | Generates a CSV of netbox devices that are not in LibreNMS
