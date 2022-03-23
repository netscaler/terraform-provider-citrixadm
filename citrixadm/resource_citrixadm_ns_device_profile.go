package citrixadm

import (
	"context"
	"encoding/json"
	"log"

	"terraform-provider-citrixadm/service"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNsDeviceProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNsDeviceProfileCreate,
		ReadContext:   resourceNsDeviceProfileRead,
		DeleteContext: resourceNsDeviceProfileDelete,
		Schema: map[string]*schema.Schema{

			//       name:
			//         description: "Profile Name"
			//         type: string
			//         format: string
			//         minLength: 1
			//         maxLength: 128
			"name": {
				// This description is used by the documentation generator and the language server.
				Description: "Profile Name",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},

			//       svm_ns_comm:
			//         description: "Communication protocol (http or https) with Instances"
			//         type: string
			//         format: string
			//         minLength: 1
			//         maxLength: 10
			"svm_ns_comm": {
				// This description is used by the documentation generator and the language server.
				Description: "Communication protocol (http or https) with Instances",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},

			//       use_global_setting_for_communication_with_ns:
			//         description: "True, if the communication with Instance needs to be global and not device specific"
			//         type: boolean
			//         format: boolean
			"use_global_setting_for_communication_with_ns": {
				// This description is used by the documentation generator and the language server.
				Description: "True, if the communication with Instance needs to be global and not device specific",
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
			},
			//       type:
			//         description: "Profile Type, This must be with in specified supported instance types: blx,sdvanvw,ns,nssdx,cbwanopt,cpx"
			//         type: string
			//         format: string
			//         minLength: 1
			//         maxLength: 128
			"type": {
				// This description is used by the documentation generator and the language server.
				Description: "Profile Type, This must be with in specified supported instance types: blx,sdvanvw,ns,nssdx,cbwanopt,cpx",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			//       snmpsecurityname:
			//         description: "SNMP v3 security name for this profile"
			//         type: string
			//         format: string
			//         maxLength: 31
			"snmpsecurityname": {
				// This description is used by the documentation generator and the language server.
				Description: "SNMP v3 security name for this profile",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			//       snmpauthprotocol:
			//         description: "SNMP v3 auth protocol for this profile"
			//         type: string
			//         format: string
			"snmpauthprotocol": {
				// This description is used by the documentation generator and the language server.
				Description: "SNMP v3 auth protocol for this profile",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			//       ssl_private_key:
			//         description: "SSL Private Key for key based authentication"
			//         type: string
			//         format: password
			"ssl_private_key": {
				// This description is used by the documentation generator and the language server.
				Description: "SSL Private Key for key based authentication",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			//       ssl_cert:
			//         description: "SSL Certificate for certificate based authentication"
			//         type: string
			//         format: string
			"ssl_cert": {
				Description: "SSL Certificate for certificate based authentication",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			//       http_port:
			//         description: "HTTP port to connect to the device"
			//         type: integer
			//         format: int32
			"http_port": {
				Description: "HTTP port to connect to the device",
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
			},
			//       ns_profile_name:
			//         description: "Profile Name, This is one of the already created Citrix ADC profiles"
			//         type: string
			//         format: string
			"ns_profile_name": {
				Description: "Profile Name, This is one of the already created Citrix ADC profiles",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			//       ssh_port:
			//         description: "SSH port to connect to the device"
			//         type: string
			//         format: string
			"ssh_port": {
				Description: "SSH port to connect to the device",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			//       password:
			//         description: "Instance credentials.Password for this profile"
			//         type: string
			//         format: password
			//         minLength: 1
			//         maxLength: 127
			"password": {
				Description: "Instance credentials.Password for this profile",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				// sensitive
				Sensitive: true,
			},

			//       snmpsecuritylevel:
			//         description: "SNMP v3 security level for this profile"
			//         type: string
			//         format: string
			"snmpsecuritylevel": {
				Description: "SNMP v3 security level for this profile",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			//       snmpcommunity:
			//         description: "SNMP community for this profile"
			//         type: string
			//         format: string
			//         maxLength: 31
			"snmpcommunity": {
				Description: "SNMP community for this profile",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			//       passphrase:
			//         description: "Passphrase with which private key is encrypted"
			//         type: string
			//         format: password
			"passphrase": {
				Description: "Passphrase with which private key is encrypted",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Sensitive:   true,
			},
			//       snmpprivprotocol:
			//         description: "SNMP v3 priv protocol for this profile"
			//         type: string
			//         format: string
			"snmpprivprotocol": {
				Description: "SNMP v3 priv protocol for this profile",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			//       https_port:
			//         description: "HTTPS port to connect to the device"
			//         type: integer
			//         format: int32
			"https_port": {
				Description: "HTTPS port to connect to the device",
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
			},
			//       username:
			//         description: "Instance credentials.Username provided in the profile will be used to contact the instance"
			//         type: string
			//         format: string
			//         minLength: 1
			//         maxLength: 127
			"username": {
				Description: "Instance credentials.Username provided in the profile will be used to contact the instance",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			//       host_password:
			//         description: "Host Password for this profile.Used for BLX form factor of ADC"
			//         type: string
			//         format: password
			//         minLength: 1
			//         maxLength: 127
			"host_password": {
				Description: "Host Password for this profile.Used for BLX form factor of ADC",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			//       max_wait_time_reboot:
			//         description: "Max waiting time to reboot Citrix ADC"
			//         type: string
			//         format: string
			"max_wait_time_reboot": {
				Description: "Max waiting time to reboot Citrix ADC",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			//       snmpprivpassword:
			//         description: "SNMP v3 priv password for this profile"
			//         type: string
			//         format: password
			//         minLength: 8
			//         maxLength: 31
			"snmpprivpassword": {
				Description: "SNMP v3 priv password for this profile",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Sensitive:   true,
			},
			//       snmpversion:
			//         description: "SNMP version for this profile"
			//         type: string
			//         format: string
			"snmpversion": {
				Description: "SNMP version for this profile",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			//       cb_profile_name:
			//         description: "Profile Name, This is one of the already created Citrix SD-WAN profiles"
			//         type: string
			//         format: string
			"cb_profile_name": {
				Description: "Profile Name, This is one of the already created Citrix SD-WAN profiles",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			//       snmpauthpassword:
			//         description: "SNMP v3 auth password for this profile"
			//         type: string
			//         format: password
			//         minLength: 8
			//         maxLength: 31
			"snmpauthpassword": {
				Description: "SNMP v3 auth password for this profile",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Sensitive:   true,
			},
			//       host_username:
			//         description: "Host User Name for this profile.Used for BLX form factor of ADC"
			//         type: string
			//         format: string
			//         minLength: 1
			//         maxLength: 127
			"host_username": {
				Description: "Host User Name for this profile.Used for BLX form factor of ADC",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceNsDeviceProfileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	log.Printf("In resourceNsDeviceProfileCreate")
	var diags diag.Diagnostics

	c := m.(*service.NitroClient)

	data := make(map[string]interface{})

	if v, ok := d.GetOk("name"); ok {
		data["name"] = v.(string)
	}

	if v, ok := d.GetOk("svm_ns_comm"); ok {
		data["svm_ns_comm"] = v.(string)
	}
	if v, ok := d.GetOk("use_global_setting_for_communication_with_ns"); ok {
		data["use_global_setting_for_communication_with_ns"] = v.(bool)
	}
	if v, ok := d.GetOk("type"); ok {
		data["type"] = v.(string)
	}
	if v, ok := d.GetOk("snmpsecurityname"); ok {
		data["snmpsecurityname"] = v.(string)
	}
	if v, ok := d.GetOk("snmpauthprotocol"); ok {
		data["snmpauthprotocol"] = v.(string)
	}
	if v, ok := d.GetOk("ssl_private_key"); ok {
		data["ssl_private_key"] = v.(string)
	}
	if v, ok := d.GetOk("ssl_cert"); ok {
		data["ssl_cert"] = v.(string)
	}
	if v, ok := d.GetOk("http_port"); ok {
		data["http_port"] = v.(int)
	}
	if v, ok := d.GetOk("ns_profile_name"); ok {
		data["ns_profile_name"] = v.(string)
	}
	if v, ok := d.GetOk("ssh_port"); ok {
		data["ssh_port"] = v.(int)
	}
	if v, ok := d.GetOk("password"); ok {
		data["password"] = v.(string)
	}
	if v, ok := d.GetOk("snmpsecuritylevel"); ok {
		data["snmpsecuritylevel"] = v.(string)
	}
	if v, ok := d.GetOk("snmpcommunity"); ok {
		data["snmpcommunity"] = v.(string)
	}
	if v, ok := d.GetOk("passphrase"); ok {
		data["passphrase"] = v.(string)
	}
	if v, ok := d.GetOk("snmpprivprotocol"); ok {
		data["snmpprivprotocol"] = v.(string)
	}
	if v, ok := d.GetOk("https_port"); ok {
		data["https_port"] = v.(int)
	}
	if v, ok := d.GetOk("username"); ok {
		data["username"] = v.(string)
	}
	if v, ok := d.GetOk("host_password"); ok {
		data["host_password"] = v.(string)
	}
	if v, ok := d.GetOk("max_wait_time_reboot"); ok {
		data["max_wait_time_reboot"] = v.(string)
	}
	if v, ok := d.GetOk("snmpprivpassword"); ok {
		data["snmpprivpassword"] = v.(string)
	}
	if v, ok := d.GetOk("snmpversion"); ok {
		data["snmpversion"] = v.(string)
	}
	if v, ok := d.GetOk("cb_profile_name"); ok {
		data["cb_profile_name"] = v.(string)
	}
	if v, ok := d.GetOk("snmpauthpassword"); ok {
		data["snmpauthpassword"] = v.(string)
	}
	if v, ok := d.GetOk("host_username"); ok {
		data["host_username"] = v.(string)
	}

	var payload []interface{}
	payload = append(payload, data)

	endpoint := "ns_device_profile"

	n := service.NitroRequestParams{
		Resource: endpoint,

		// ResourcePath:       https://adm.cloud.com/massvc/{{customerid}}/nitro/v2/config/ns_device_profile
		// ResourcePath:       fmt.Sprintf("%s/massvc/%s/nitro/v2/config/%s", c.host, c.customerid, endpoint),
		ResourceData:       payload,
		Method:             "POST",
		SuccessStatusCodes: []int{200, 201},
	}

	body, err := c.MakeNitroRequest(n)
	if err != nil {
		return diag.FromErr(err)
	}
	var returnData map[string]interface{}

	err = json.Unmarshal(body, &returnData)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("return data %v", returnData)

	resourceID := returnData[endpoint].([]interface{})[0].(map[string]interface{})["id"].(string)
	log.Printf("id %s", resourceID)

	d.SetId(resourceID)

	resourceNsDeviceProfileRead(ctx, d, m)

	return diags
}

func resourceNsDeviceProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	log.Printf("In resourceNsDeviceProfileRead")
	var diags diag.Diagnostics
	c := m.(*service.NitroClient)

	resourceID := d.Id()
	endpoint := "ns_device_profile"

	n := service.NitroRequestParams{
		// ResourcePath:       fmt.Sprintf("massvc/%s/nitro/v2/config/%s", c.customerid, endpoint),
		Method: "GET",
		// Headers:            map[string]string{},
		Resource:           endpoint,
		ResourceID:         resourceID,
		ResourceData:       d,
		SuccessStatusCodes: []int{200},
	}

	body, err := c.MakeNitroRequest(n)
	if err != nil {
		return diag.FromErr(err)
	}
	var returnData map[string]interface{}

	err = json.Unmarshal(body, &returnData)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("return data %v", returnData)
	// {
	// 	"ns_device_profile": [
	// 		{
	// 			"type": "ns",
	// 			"is_default": "true",
	// 			"svm_ns_comm": "https",
	// 			"tenant_id": "",
	// 			"id": "cbd6cf21-b654-498d-9475-89c12f7b27c0",
	// 			"is_backup": "false",
	// 			"name": "ns_nsroot_profile",
	// 			"use_global_setting_for_communication_with_ns": "true",
	// 			"snmpsecuritylevel": "",
	// 			"act_id": "",
	// 			"ssl_cert": "",
	// 			"host_username": "",
	// 			"sync_operation": "true",
	// 			"ns_profile_name": "",
	// 			"cb_profile_name": "",
	// 			"snmpsecurityname": "",
	// 			"https_port": "443",
	// 			"username": "nsroot",
	// 			"snmpauthprotocol": "",
	// 			"ssh_port": "22",
	// 			"snmpversion": "",
	// 			"max_wait_time_reboot": "1800",
	// 			"http_port": "80",
	// 			"snmpprivprotocol": "",
	// 			"snmpcommunity": ""
	// 		}
	// 	]
	// }
	getResponseData := returnData[endpoint].([]interface{})[0].(map[string]interface{})

	log.Println("getResponseData", getResponseData)

	// d.Set("act_id", getResponseData["act_id"].(string))
	d.Set("cb_profile_name", getResponseData["cb_profile_name"].(string))
	d.Set("host_username", getResponseData["host_username"].(string))
	d.Set("http_port", getResponseData["http_port"].(string)) // FIXME: Ask George: panic: interface conversion: interface {} is string, not int
	d.Set("https_port", getResponseData["https_port"].(string))
	// d.Set("id", getResponseData["id"].(string))
	// d.Set("is_backup", getResponseData["is_backup"].(string))
	// d.Set("is_default", getResponseData["is_default"].(string))
	d.Set("max_wait_time_reboot", getResponseData["max_wait_time_reboot"].(string))
	d.Set("name", getResponseData["name"].(string))
	d.Set("ns_profile_name", getResponseData["ns_profile_name"].(string))
	d.Set("snmpauthprotocol", getResponseData["snmpauthprotocol"].(string))
	d.Set("snmpcommunity", getResponseData["snmpcommunity"].(string))
	d.Set("snmpprivprotocol", getResponseData["snmpprivprotocol"].(string))
	d.Set("snmpsecuritylevel", getResponseData["snmpsecuritylevel"].(string))
	d.Set("snmpsecurityname", getResponseData["snmpsecurityname"].(string))
	d.Set("snmpversion", getResponseData["snmpversion"].(string))
	d.Set("ssh_port", getResponseData["ssh_port"].(string))
	d.Set("ssl_cert", getResponseData["ssl_cert"].(string))
	d.Set("svm_ns_comm", getResponseData["svm_ns_comm"].(string))
	// d.Set("sync_operation", getResponseData["sync_operation"].(string))
	// d.Set("tenant_id", getResponseData["tenant_id"].(string))
	d.Set("type", getResponseData["type"].(string))
	d.Set("use_global_setting_for_communication_with_ns", getResponseData["use_global_setting_for_communication_with_ns"].(string)) // FIXME: panic: interface conversion: interface {} is string, not bool
	d.Set("username", getResponseData["username"].(string))

	return diags
}

func resourceNsDeviceProfileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceNsDeviceProfileDelete")
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c := m.(*service.NitroClient)

	endpoint := "ns_device_profile"
	resourceID := d.Id()

	n := service.NitroRequestParams{
		// ResourcePath:       fmt.Sprintf("adcaas/nitro/v1/config/endpoints/%s", environmentId),
		Method:             "DELETE",
		Resource:           endpoint,
		ResourceID:         resourceID,
		SuccessStatusCodes: []int{200, 204},
	}

	body, err := c.MakeNitroRequest(n)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("delete response %v", body)

	d.SetId("")

	return diags
}
