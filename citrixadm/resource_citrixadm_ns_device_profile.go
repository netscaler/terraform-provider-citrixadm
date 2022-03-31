package citrixadm

import (
	"context"
	"log"

	"terraform-provider-citrixadm/service"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNsDeviceProfile() *schema.Resource {
	return &schema.Resource{
		Description:   "Configuration for Device profile for  Citrix ADC(MPX/VPX/CPX/Gateway) instances  resource",
		CreateContext: resourceNsDeviceProfileCreate,
		ReadContext:   resourceNsDeviceProfileRead,
		UpdateContext: resourceNsDeviceProfileUpdate,
		DeleteContext: resourceNsDeviceProfileDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Profile Name",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"svm_ns_comm": {
				Description: "Communication protocol (http or https) with Instances",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"use_global_setting_for_communication_with_ns": {
				Description: "True, if the communication with Instance needs to be global and not device specific",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"type": {
				Description: "Profile Type, This must be with in specified supported instance types: blx,sdvanvw,ns,nssdx,cbwanopt,cpx",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
			},
			"snmpsecurityname": {
				Description: "SNMP v3 security name for this profile",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
			},
			"snmpauthprotocol": {
				Description: "SNMP v3 auth protocol for this profile",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
			},
			"ssl_private_key": {
				Description: "SSL Private Key for key based authentication",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"ssl_cert": {
				Description: "SSL Certificate for certificate based authentication",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"http_port": {
				Description: "HTTP port to connect to the device",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
			"ns_profile_name": {
				Description: "Profile Name, This is one of the already created Citrix ADC profiles",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
			},
			"ssh_port": {
				Description: "SSH port to connect to the device",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"password": {
				Description: "Instance credentials.Password for this profile",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
			},
			"snmpsecuritylevel": {
				Description: "SNMP v3 security level for this profile",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
			},
			"snmpcommunity": {
				Description: "SNMP community for this profile",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
			},
			"passphrase": {
				Description: "Passphrase with which private key is encrypted",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Sensitive:   true,
			},
			"snmpprivprotocol": {
				Description: "SNMP v3 priv protocol for this profile",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
			},
			"https_port": {
				Description: "HTTPS port to connect to the device",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
			"username": {
				Description: "Instance credentials.Username provided in the profile will be used to contact the instance",
				Type:        schema.TypeString,
				Required:    true,
			},
			"host_password": {
				Description: "Host Password for this profile.Used for BLX form factor of ADC",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"max_wait_time_reboot": {
				Description: "Max waiting time to reboot Citrix ADC",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"snmpprivpassword": {
				Description: "SNMP v3 priv password for this profile",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Sensitive:   true,
			},
			"snmpversion": {
				Description: "SNMP version for this profile",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
			},
			"cb_profile_name": {
				Description: "Profile Name, This is one of the already created Citrix SD-WAN profiles",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"snmpauthpassword": {
				Description: "SNMP v3 auth password for this profile",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Sensitive:   true,
			},
			"host_username": {
				Description: "Host User Name for this profile.Used for BLX form factor of ADC",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func getNsDeviceProfilePayload(d *schema.ResourceData) interface{} {
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
		data["ssh_port"] = v.(string)
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

	// var payload []interface{}
	// payload = append(payload, data)

	return data

}

func resourceNsDeviceProfileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceNsDeviceProfileCreate")

	c := m.(*service.NitroClient)

	endpoint := "ns_device_profile"

	returnData, err := c.AddResource(endpoint, getNsDeviceProfilePayload(d))

	if err != nil {
		return diag.Errorf("error creating ns_device_profile: %s", err.Error())
	}

	resourceID := returnData[endpoint].([]interface{})[0].(map[string]interface{})["id"].(string)
	log.Printf("id %s", resourceID)

	d.SetId(resourceID)

	return resourceNsDeviceProfileRead(ctx, d, m)
}

func resourceNsDeviceProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceNsDeviceProfileRead")
	var diags diag.Diagnostics
	c := m.(*service.NitroClient)

	resourceID := d.Id()
	endpoint := "ns_device_profile"

	returnData, err := c.GetResource(endpoint, resourceID)
	if err != nil {
		return diag.Errorf("error reading ns_device_profile %s: %s", resourceID, err.Error())
	}

	getResponseData := returnData[endpoint].([]interface{})[0].(map[string]interface{})

	log.Println("getResponseData", getResponseData)

	d.Set("cb_profile_name", getResponseData["cb_profile_name"].(string))
	d.Set("host_username", getResponseData["host_username"].(string))
	d.Set("http_port", getResponseData["http_port"].(string))   // FIXME: API problem. Though API schema defines this as int, the GET response returns as string
	d.Set("https_port", getResponseData["https_port"].(string)) // FIXME: API problem. Though API schema defines this as int, the GET response returns as string
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
	d.Set("type", getResponseData["type"].(string))
	d.Set("use_global_setting_for_communication_with_ns", getResponseData["use_global_setting_for_communication_with_ns"].(string)) // FIXME: API problem. Though API schema defines this as bool, the GET response returns as string
	d.Set("username", getResponseData["username"].(string))

	return diags
}

func resourceNsDeviceProfileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceNsDeviceProfileUpdate")
	c := m.(*service.NitroClient)

	resourceID := d.Id()
	endpoint := "ns_device_profile"

	_, err := c.UpdateResource(endpoint, getManagedDevicePayload(d), resourceID)

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNsDeviceProfileRead(ctx, d, m)

}

func resourceNsDeviceProfileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceNsDeviceProfileDelete")
	var diags diag.Diagnostics

	c := m.(*service.NitroClient)

	endpoint := "ns_device_profile"
	resourceID := d.Id()
	err := c.DeleteResource(endpoint, resourceID)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
