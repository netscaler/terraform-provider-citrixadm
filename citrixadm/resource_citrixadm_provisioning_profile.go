package citrixadm

import (
	"context"
	"encoding/json"
	"log"

	"terraform-provider-citrixadm/service"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceProvisioningProfile() *schema.Resource {
	return &schema.Resource{
		Description:   "Configuration for Device profile for  Citrix ADC(MPX/VPX/CPX/Gateway) instances  resource",
		CreateContext: resourceProvisioningProfileCreate,
		ReadContext:   resourceProvisioningProfileRead,
		// UpdateContext: resourceProvisioningProfileUpdate, // TODO: For now, UPDATE was using different API and was giving error. Hence going with forced re-creation.
		// DeleteContext: resourceProvisioningProfileDelete, // TODO: For now, DELETE is effective as prescribed by NITRO documentation
		DeleteContext: schema.NoopContext,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the ProvisioningProfile",
				Required:    true,
				ForceNew:    true,
			},
			"instance_type": {
				Type:        schema.TypeString,
				Description: "Only NetScaler is supported as of now",
				Required:    true,
				ForceNew:    true,
			},
			"site_id": {
				Type:        schema.TypeString,
				Description: "Reference to MAS site which has location info where instance has to be provisioned",
				Required:    true,
				ForceNew:    true,
			},
			"platform_type": {
				Type:        schema.TypeString,
				Description: "Platform type",
				Optional:    true,
				ForceNew:    true,
			},
			"deployment_details": {
				Type:        schema.TypeString,
				Description: "Deployment Details",
				Optional:    true,
				ForceNew:    true,
			},
			"mas_registration_details": {
				Type:        schema.TypeList,
				Description: "MAS Registration Details",
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mas_agent_id": {
							Type:        schema.TypeString,
							Description: "Reference to MAS Agent that has to be used in order to add and manage provisioned instance in MAS",
							Required:    true,
						},
						"access_profile_id": {
							Type:        schema.TypeString,
							Description: "Reference to Instance/Device Access Profile to be set for instance being provisioned",
							Optional:    true,
						},
					},
				},
			},
			"instance_capacity_details": {
				Type:        schema.TypeList,
				Description: "Instance Capacity Details",
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"config_job_templates": {
							Type:        schema.TypeList,
							Description: "Config Job Templates",
							Required:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func getProvisioningProfilePayload(d *schema.ResourceData) interface{} {
	data := make(map[string]interface{})

	data["name"] = d.Get("name").(string)
	data["instance_type"] = d.Get("instance_type").(string)
	data["site_id"] = d.Get("site_id").(string)
	if v, ok := d.GetOk("platform_type"); ok {
		data["platform_type"] = v.(string)
	}
	if v, ok := d.GetOk("deployment_details"); ok {
		// json unmarshall deployment_details
		var deploymentDetails interface{}
		err := json.Unmarshal([]byte(v.(string)), &deploymentDetails)
		if err != nil {
			log.Printf("[DEBUG] Error unmarshalling deployment_details: %s", err)
		}
		data["deployment_details"] = deploymentDetails
	}
	if v, ok := d.GetOk("mas_registration_details"); ok {
		data["mas_registration_details"] = v.([]interface{})[0].(map[string]interface{})
	}
	if v, ok := d.GetOk("instance_capacity_details"); ok {
		data["instance_capacity_details"] = v.([]interface{})[0].(map[string]interface{})
	}
	return data
}

func resourceProvisioningProfileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceProvisioningProfileCreate")

	c := m.(*service.NitroClient)

	endpoint := "provisioning_profiles"

	returnData, err := c.AddResource(endpoint, getProvisioningProfilePayload(d))

	if err != nil {
		return diag.Errorf("error creating provisioning_profiles: %s", err.Error())
	}

	resourceID := returnData[endpoint].([]interface{})[0].(map[string]interface{})["id"].(string)
	log.Printf("id %s", resourceID)

	d.SetId(resourceID)
	return resourceProvisioningProfileRead(ctx, d, m)
}

func resourceProvisioningProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceProvisioningProfileRead")
	var diags diag.Diagnostics
	c := m.(*service.NitroClient)

	resourceID := d.Id()
	endpoint := "provisioning_profiles"

	returnData, err := c.GetResource(endpoint, resourceID)
	if err != nil {
		return diag.Errorf("error reading ns_device_profile %s: %s", resourceID, err.Error())
	}

	getResponseData := returnData[endpoint].([]interface{})[0].(map[string]interface{})

	log.Println("getResponseData", getResponseData)

	d.Set("name", getResponseData["name"].(string))
	d.Set("instance_type", getResponseData["instance_type"].(string))
	d.Set("site_id", getResponseData["site_id"].(string))
	if v, ok := getResponseData["platform_type"]; ok {
		d.Set("platform_type", v.(string))
	}
	if v, ok := getResponseData["deployment_details"]; ok {
		d.Set("deployment_details", v.(map[string]interface{}))
	}
	if v, ok := getResponseData["mas_registration_details"]; ok {
		d.Set("mas_registration_details", flattenProvisioningMasRegistrationDetails(v.(map[string]interface{})))
	}
	if v, ok := getResponseData["instance_capacity_details"]; ok {
		d.Set("instance_capacity_details", flattenProvisioningInstanceCapacityDetails(v.(map[string]interface{})))
	}
	return diags
}

func flattenProvisioningMasRegistrationDetails(masRegistrationDetails map[string]interface{}) []interface{} {
	s := make(map[string]interface{})
	s["mas_agent_id"] = masRegistrationDetails["mas_agent_id"].(string)
	if v, ok := masRegistrationDetails["access_profile_id"]; ok {
		s["access_profile_id"] = v.(string)
	}
	return []interface{}{s}
}

func flattenProvisioningInstanceCapacityDetails(instanceCapacityDetails map[string]interface{}) []interface{} {
	s := make(map[string]interface{})
	if v, ok := instanceCapacityDetails["config_job_templates"]; ok {
		s["config_job_templates"] = v.([]interface{})
	}
	return []interface{}{s}
}

// TODO: For now, UPDATE was using different API and was giving error. Hence going with forced re-creation.
// func resourceProvisioningProfileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
// 	log.Printf("In resourceProvisioningProfileUpdate")
// 	c := m.(*service.NitroClient)

// 	resourceID := d.Id()
// 	endpoint := "ns_device_profile"

// 	_, err := c.UpdateResource(endpoint, getProvisioningProfilePayload(d), resourceID)

// 	if err != nil {
// 		return diag.FromErr(err)
// 	}

// 	return resourceProvisioningProfileRead(ctx, d, m)

// }

// TODO: For now, DELETE is effective as prescribed by NITRO documentation
// func resourceProvisioningProfileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
// 	log.Printf("In resourceProvisioningProfileDelete")
// 	var diags diag.Diagnostics

// 	c := m.(*service.NitroClient)

// 	endpoint := "provisioning_profiles"
// 	resourceID := d.Id()
// 	_, err := c.DeleteResource(endpoint, resourceID)

// 	if err != nil {
// 		return diag.FromErr(err)
// 	}

// 	d.SetId("")

// 	return diags
// }
