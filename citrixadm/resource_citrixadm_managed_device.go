package citrixadm

import (
	"context"
	"errors"
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

			// "std_bw_config": {
			// 	Description: "Standard Bandwidth running",
			// 	Type:        schema.TypeInt,
			// 	Optional:    true,
			// },
			// "description": {
			// 	Description: "Description of managed device",
			// 	Type:        schema.TypeString,
			// 	Optional:    true,
			// },
			// "vcpu_config": {
			// 	Description: "Number of vCPU allocated for the device",
			// 	Type:        schema.TypeInt,
			// 	Optional:    true,
			// },
			// "ent_bw_config": {
			// 	Description: "Enterprise Bandwidth configured",
			// 	Type:        schema.TypeInt,
			// 	Optional:    true,
			// },
			// "instance_config": {
			// 	Description: "Instance license running",
			// 	Type:        schema.TypeInt,
			// 	Optional:    true,
			// },
			// "is_managed": {
			// 	Description: "Is Managed",
			// 	Type:        schema.TypeBool,
			// 	Optional:    true,
			// },
			// "servicepackage": {
			// 	Description: "Service Package Name of the device",
			// 	Type:        schema.TypeString,
			// 	Optional:    true,
			// },
			// "license_edition": {
			// 	Description: "Edition of instance",
			// 	Type:        schema.TypeString,
			// 	Optional:    true,
			// },
			// "isolation_policy": {
			// 	Description: "Isolation Policy of the Device",
			// 	Type:        schema.TypeString,
			// 	Optional:    true,
			// },
			// "plt_bw_config": {
			// 	Description: "Platinum Bandwidth configured",
			// 	Type:        schema.TypeInt,
			// 	Optional:    true,
			// },
			// "peer_device_ip": {
			// 	Description: "Peer Device IP address for instance of type BLX ADC.",
			// 	Type:        schema.TypeString,
			// 	Optional:    true,
			// },
			// "file_name": {
			// 	Description: "File name which contains comma separated instances to be  discovered",
			// 	Type:        schema.TypeString,
			// 	Optional:    true,
			// },
			// "file_location_path": {
			// 	Description: "File Location on Client for upload/download",
			// 	Type:        schema.TypeString,
			// 	Optional:    true,
			// },
			// "peer_host_device_ip": {
			// 	Description: "Peer Host Device IP Address for instance of type BLX ADC.",
			// 	Type:        schema.TypeString,
			// 	Optional:    true,
			// },
			// "device_host_ip": {
			// 	Description: "Device Host IP Address for instance of type BLX ADC.",
			// 	Type:        schema.TypeString,
			// 	Optional:    true,
			// },
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

func getManagedDevicePayload(d *schema.ResourceData) interface{} {
	data := make(map[string]interface{})

	data["ip_address"] = d.Get("ip_address").(string)

	data["profile_name"] = d.Get("profile_name").(string)

	data["datacenter_id"] = d.Get("datacenter_id").(string)
	data["agent_id"] = d.Get("agent_id").(string)

	// if v, ok := d.GetOk("std_bw_config"); ok {
	// 	data["std_bw_config"] = v.(int)
	// }

	// if v, ok := d.GetOk("description"); ok {
	// 	data["description"] = v.(string)
	// }

	// if v, ok := d.GetOk("vcpu_config"); ok {
	// 	data["vcpu_config"] = v.(int)
	// }
	// if v, ok := d.GetOk("ent_bw_config"); ok {
	// 	data["ent_bw_config"] = v.(int)
	// }

	// if v, ok := d.GetOk("instance_config"); ok {
	// 	data["instance_config"] = v.(string)
	// }
	// if v, ok := d.GetOk("is_managed"); ok {
	// 	data["is_managed"] = v.(bool)
	// }

	// if v, ok := d.GetOk("servicepackage"); ok {
	// 	data["servicepackage"] = v.(string)
	// }

	// if v, ok := d.GetOk("license_edition"); ok {
	// 	data["license_edition"] = v.(string)
	// }
	// if v, ok := d.GetOk("isolation_policy"); ok {
	// 	data["isolation_policy"] = v.(string)
	// }
	// if v, ok := d.GetOk("plt_bw_config"); ok {
	// 	data["plt_bw_config"] = v.(int)
	// }
	// if v, ok := d.GetOk("peer_device_ip"); ok {
	// 	data["peer_device_ip"] = v.(string)
	// }
	// if v, ok := d.GetOk("file_name"); ok {
	// 	data["file_name"] = v.(string)
	// }
	// if v, ok := d.GetOk("file_location_path"); ok {
	// 	data["file_location_path"] = v.(string)
	// }
	// if v, ok := d.GetOk("peer_host_device_ip"); ok {
	// 	data["peer_host_device_ip"] = v.(string)
	// }
	// if v, ok := d.GetOk("device_host_ip"); ok {
	// 	data["device_host_ip"] = v.(string)
	// }
	if v, ok := d.GetOk("entity_tag"); ok {
		data["entity_tag"] = v.([]interface{})
	}

	// var payload []interface{}
	// payload = append(payload, data)

	return data
}

func resourceManagedDeviceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceManagedDeviceCreate")

	c := m.(*service.NitroClient)

	endpoint := "managed_device"

	returnData, err := c.AddResource(endpoint, getManagedDevicePayload(d))

	if err != nil {
		return diag.Errorf("unable to create Managed Device: %s", err.Error())
	}

	activityStatusID := returnData[endpoint].([]interface{})[0].(map[string]interface{})["act_id"].(string)

	// Wait for activity to complete
	log.Printf("Waiting for activity to complete")
	err = c.WaitForActivityCompletion(activityStatusID, time.Duration(c.ActivityTimeout)*time.Second)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceID, err := getManagedDeviceID(c, d.Get("ip_address").(string))

	if err != nil {
		return diag.FromErr(errors.New("Failed to find resource ID"))
	}

	d.SetId(resourceID)
	return resourceManagedDeviceRead(ctx, d, m)
}

func getManagedDeviceID(c *service.NitroClient, ipAddress string) (string, error) {
	endpoint := "managed_device"
	returnData, err := c.GetAllResource(endpoint)
	if err != nil {
		return "", err
	}

	for _, v := range returnData[endpoint].([]interface{}) {
		if v.(map[string]interface{})["ip_address"].(string) == ipAddress {
			return v.(map[string]interface{})["id"].(string), nil
		}
	}
	return "", errors.New("Failed to find managed device resource ID with IP: " + ipAddress)
}

func resourceManagedDeviceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceManagedDeviceRead")
	var diags diag.Diagnostics
	c := m.(*service.NitroClient)

	resourceID := d.Id()
	endpoint := "managed_device"

	returnData, err := c.GetResource(endpoint, resourceID)
	if err != nil {
		return diag.FromErr(err)
	}

	getResponseData := returnData[endpoint].([]interface{})[0].(map[string]interface{})

	d.Set("agent_id", getResponseData["agent_id"].(string))
	d.Set("ip_address", getResponseData["ip_address"].(string))
	d.Set("profile_name", getResponseData["profile_name"].(string))
	d.Set("datacenter_id", getResponseData["datacenter_id"].(string))
	d.Set("entity_tag", getResponseData["entity_tag"].([]interface{}))

	return diags
}

func resourceManagedDeviceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceManagedDeviceUpdate")
	c := m.(*service.NitroClient)

	resourceID := d.Id()
	endpoint := "managed_device"

	_, err := c.UpdateResource(endpoint, getManagedDevicePayload(d), resourceID)

	if err != nil {
		return diag.FromErr(err)
	}

	// if license_edition or plt_bw_config has changed then allocate_license
	if d.HasChange("license_edition") || d.HasChange("plt_bw_config") {
		data := make(map[string]interface{})
		data["id"] = resourceID
		data["license_edition"] = d.Get("license_edition").(string)
		data["plt_bw_config"] = d.Get("plt_bw_config").(int)
		var payload []interface{}
		payload = append(payload, data)

		_, err = c.AddResourceWithActionParams(endpoint, payload, "allocate_license")
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceManagedDeviceRead(ctx, d, m)
}

func resourceManagedDeviceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceManagedDeviceDelete")
	var diags diag.Diagnostics

	c := m.(*service.NitroClient)

	endpoint := "managed_device"
	resourceID := d.Id()

	_, err := c.DeleteResource(endpoint, resourceID)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
