package citrixadm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"terraform-provider-citrixadm/service"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceManagedDevice() *schema.Resource {
	return &schema.Resource{
		Description:   "Configuration for Managed Device resource",
		CreateContext: resourceManagedDeviceCreate,
		ReadContext:   resourceManagedDeviceRead,
		UpdateContext: resourceManagedDeviceUpdate,
		DeleteContext: resourceManagedDeviceDelete,
		Schema: map[string]*schema.Schema{
			"ip_address": {
				Description: "IP Address for this managed device",
				Type:        schema.TypeString,
				Required:    true,
			},
			"profile_name": {
				Description: "Device Profile Name that is attached with this managed device",
				Type:        schema.TypeString,
				Required:    true,
			},
			"datacenter_id": {
				Description: "Datacenter Id is system generated key for data center",
				Type:        schema.TypeString,
				Required:    true,
			},
			"agent_id": {
				Description: "Agent Id",
				Type:        schema.TypeString,
				Required:    true,
			},

			"std_bw_config": {
				Description: "Standard Bandwidth running",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"description": {
				Description: "Description of managed device",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"vcpu_config": {
				Description: "Number of vCPU allocated for the device",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"ent_bw_config": {
				Description: "Enterprise Bandwidth configured",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"instance_config": {
				Description: "Instance license running",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"is_managed": {
				Description: "Is Managed",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"servicepackage": {
				Description: "Service Package Name of the device",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"license_edition": {
				Description: "Edition of instance",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"isolation_policy": {
				Description: "Isolation Policy of the Device",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"plt_bw_config": {
				Description: "Platinum Bandwidth configured",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"peer_device_ip": {
				Description: "Peer Device IP address for instance of type BLX ADC.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"file_name": {
				Description: "File name which contains comma separated instances to be  discovered",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"file_location_path": {
				Description: "File Location on Client for upload/download",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"peer_host_device_ip": {
				Description: "Peer Host Device IP Address for instance of type BLX ADC.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"device_host_ip": {
				Description: "Device Host IP Address for instance of type BLX ADC.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"entity_tag": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"prop_key": {
							Description: "Property key",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"prop_value": {
							Description: "Property value",
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func getManagedDevicePayload(d *schema.ResourceData) []interface{} {
	data := make(map[string]interface{})

	if v, ok := d.GetOk("ip_address"); ok {
		data["ip_address"] = v.(string)
	}

	if v, ok := d.GetOk("profile_name"); ok {
		data["profile_name"] = v.(string)
	}

	if v, ok := d.GetOk("datacenter_id"); ok {
		data["datacenter_id"] = v.(string)
	}
	if v, ok := d.GetOk("agent_id"); ok {
		data["agent_id"] = v.(string)
	}

	if v, ok := d.GetOk("std_bw_config"); ok {
		data["std_bw_config"] = v.(int)
	}

	if v, ok := d.GetOk("description"); ok {
		data["description"] = v.(string)
	}

	if v, ok := d.GetOk("vcpu_config"); ok {
		data["vcpu_config"] = v.(int)
	}
	if v, ok := d.GetOk("ent_bw_config"); ok {
		data["ent_bw_config"] = v.(int)
	}

	if v, ok := d.GetOk("instance_config"); ok {
		data["instance_config"] = v.(string)
	}
	if v, ok := d.GetOk("is_managed"); ok {
		data["is_managed"] = v.(bool)
	}

	if v, ok := d.GetOk("servicepackage"); ok {
		data["servicepackage"] = v.(string)
	}

	if v, ok := d.GetOk("license_edition"); ok {
		data["license_edition"] = v.(string)
	}
	if v, ok := d.GetOk("isolation_policy"); ok {
		data["isolation_policy"] = v.(string)
	}
	if v, ok := d.GetOk("plt_bw_config"); ok {
		data["plt_bw_config"] = v.(int)
	}
	if v, ok := d.GetOk("peer_device_ip"); ok {
		data["peer_device_ip"] = v.(string)
	}
	if v, ok := d.GetOk("file_name"); ok {
		data["file_name"] = v.(string)
	}
	if v, ok := d.GetOk("file_location_path"); ok {
		data["file_location_path"] = v.(string)
	}
	if v, ok := d.GetOk("peer_host_device_ip"); ok {
		data["peer_host_device_ip"] = v.(string)
	}
	if v, ok := d.GetOk("device_host_ip"); ok {
		data["device_host_ip"] = v.(string)
	}
	if v, ok := d.GetOk("entity_tag"); ok {
		data["entity_tag"] = v.([]interface{})
	}

	var payload []interface{}
	payload = append(payload, data)

	return payload

}
func resourceManagedDeviceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceManagedDeviceCreate")

	c := m.(*service.NitroClient)

	endpoint := "managed_device"

	n := service.NitroRequestParams{
		Resource:           endpoint,
		ResourcePath:       fmt.Sprintf("massvc/%s/nitro/v2/config/%s", c.CustomerID, endpoint),
		ResourceData:       getManagedDevicePayload(d),
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

	activityStatusID := returnData[endpoint].([]interface{})[0].(map[string]interface{})["act_id"].(string)

	// Wait for activity to complete
	log.Printf("Waiting for activity to complete")
	err = c.WaitForActivityCompletion(activityStatusID, time.Duration(c.ActivityTimeout)*time.Second)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceID := func() string {
		n := service.NitroRequestParams{
			Resource:           endpoint,
			Method:             "GET",
			ResourcePath:       fmt.Sprintf("massvc/%s/nitro/v2/config/%s", c.CustomerID, endpoint),
			SuccessStatusCodes: []int{200},
		}
		body, err := c.MakeNitroRequest(n)
		if err != nil {
			return ""
		}
		var returnData map[string]interface{}

		err = json.Unmarshal(body, &returnData)
		if err != nil {
			return ""
		}
		log.Printf("return data %v", returnData)

		for _, v := range returnData[endpoint].([]interface{}) {
			if v.(map[string]interface{})["ip_address"].(string) == d.Get("ip_address").(string) {
				return v.(map[string]interface{})["id"].(string)
			}
		}
		return ""
	}()

	if resourceID == "" {
		return diag.FromErr(errors.New("Failed to find resource ID"))
	}

	log.Printf("id %s", resourceID)

	d.SetId(resourceID)
	return resourceManagedDeviceRead(ctx, d, m)
}

func resourceManagedDeviceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceManagedDeviceRead")
	var diags diag.Diagnostics
	c := m.(*service.NitroClient)

	resourceID := d.Id()
	endpoint := "managed_device"

	n := service.NitroRequestParams{
		ResourcePath:       fmt.Sprintf("massvc/%s/nitro/v2/config/%s/%s", c.CustomerID, endpoint, resourceID),
		Method:             "GET",
		Resource:           endpoint,
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
	getResponseData := returnData[endpoint].([]interface{})[0].(map[string]interface{})

	log.Println("getResponseData", getResponseData)

	d.Set("is_managed", getResponseData["is_managed"].(string))
	d.Set("std_bw_config", getResponseData["std_bw_config"].(string))
	d.Set("description", getResponseData["description"].(string))
	d.Set("instance_config", getResponseData["instance_config"].(string))
	d.Set("vcpu_config", getResponseData["vcpu_config"].(string))
	d.Set("agent_id", getResponseData["agent_id"].(string))
	d.Set("servicepackage", getResponseData["servicepackage"].(string))
	d.Set("ip_address", getResponseData["ip_address"].(string))
	d.Set("plt_bw_config", getResponseData["plt_bw_config"].(string))
	d.Set("isolation_policy", getResponseData["isolation_policy"].(string))
	d.Set("profile_name", getResponseData["profile_name"].(string))
	d.Set("ent_bw_config", getResponseData["ent_bw_config"].(string))
	d.Set("datacenter_id", getResponseData["datacenter_id"].(string))
	d.Set("license_edition", getResponseData["license_edition"].(string))
	d.Set("template_interval", getResponseData["template_interval"].(string))
	d.Set("is_licensed", getResponseData["is_licensed"].(string))
	d.Set("contactperson", getResponseData["contactperson"].(string))
	d.Set("peer_host_device_ip", getResponseData["peer_host_device_ip"].(string))
	d.Set("device_host_ip", getResponseData["device_host_ip"].(string))
	d.Set("peer_device_ip", getResponseData["peer_device_ip"].(string))
	d.Set("file_location_path", getResponseData["file_location_path"].(string))
	d.Set("file_name", getResponseData["file_name"].(string))
	d.Set("entity_tag", getResponseData["entity_tag"].([]interface{}))

	return diags
}

func resourceManagedDeviceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceManagedDeviceUpdate")
	c := m.(*service.NitroClient)

	resourceID := d.Id()
	endpoint := "managed_device"

	n := service.NitroRequestParams{
		Resource:           endpoint,
		ResourcePath:       fmt.Sprintf("massvc/%s/nitro/v2/config/%s/%s", c.CustomerID, endpoint, resourceID),
		ResourceData:       getManagedDevicePayload(d),
		Method:             "PUT",
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
	return resourceManagedDeviceRead(ctx, d, m)
}

func resourceManagedDeviceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceManagedDeviceDelete")
	var diags diag.Diagnostics

	c := m.(*service.NitroClient)

	endpoint := "managed_device"
	resourceID := d.Id()

	n := service.NitroRequestParams{
		ResourcePath:       fmt.Sprintf("massvc/%s/nitro/v2/config/%s/%s", c.CustomerID, endpoint, resourceID),
		Method:             "DELETE",
		Resource:           endpoint,
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
