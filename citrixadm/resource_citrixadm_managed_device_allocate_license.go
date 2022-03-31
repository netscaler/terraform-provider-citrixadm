package citrixadm

import (
	"context"
	"log"

	"terraform-provider-citrixadm/service"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceManagedDeviceAllocateLicense() *schema.Resource {
	return &schema.Resource{
		Description:   "Configuration for Managed Device resource",
		CreateContext: resourceManagedDeviceAllocateLicenseCreate,
		ReadContext:   resourceManagedDeviceAllocateLicenseRead,
		UpdateContext: resourceManagedDeviceAllocateLicenseCreate,
		DeleteContext: resourceManagedDeviceAllocateLicenseDelete,
		Schema: map[string]*schema.Schema{
			"managed_device_id": {
				Description: "ID of the managed device to which license is to be allocated",
				Type:        schema.TypeString,
				Required:    true,
			},
			"license_edition": {
				Description: "Edition of instance",
				Type:        schema.TypeString,
				Required:    true,
			},
			"plt_bw_config": {
				Description: "Platinum Bandwidth configured",
				Type:        schema.TypeInt,
				Required:    true,
			},
		},
	}
}

func resourceManagedDeviceAllocateLicenseCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceManagedDeviceAllocateLicenseCreate")
	var diags diag.Diagnostics

	c := m.(*service.NitroClient)

	endpoint := "managed_device"
	data := make(map[string]interface{})

	data["id"] = d.Get("managed_device_id").(string)
	data["license_edition"] = d.Get("license_edition").(string)
	data["plt_bw_config"] = d.Get("plt_bw_config").(int)

	_, err := c.AddResourceWithActionParams(endpoint, data, "allocate_license")

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(d.Get("managed_device_id").(string))

	return diags

	// FIXME: As of now, there is API problem that the immediate READ (GET) after license allocation gives empty value for license_edition and plt_bw_config. Hence skipping READ for now.
	// return resourceManagedDeviceAllocateLicenseRead(ctx, d, m)

}

func resourceManagedDeviceAllocateLicenseRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceManagedDeviceAllocateLicenseRead")
	var diags diag.Diagnostics
	c := m.(*service.NitroClient)

	resourceID := d.Id()
	endpoint := "managed_device"

	returnData, err := c.GetResource(endpoint, resourceID)
	if err != nil {
		return diag.FromErr(err)
	}

	getResponseData := returnData[endpoint].([]interface{})[0].(map[string]interface{})

	d.Set("managed_device_id", getResponseData["id"].(string))
	d.Set("license_edition", getResponseData["license_edition"].(string))
	d.Set("plt_bw_config", getResponseData["plt_bw_config"].(string))

	return diags
}

func resourceManagedDeviceAllocateLicenseDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceManagedDeviceAllocateLicenseDelete")
	var diags diag.Diagnostics

	c := m.(*service.NitroClient)

	endpoint := "managed_device"

	data := make(map[string]interface{})
	data["id"] = d.Id()
	data["license_edition"] = d.Get("license_edition").(string)
	data["plt_bw_config"] = 0

	_, err := c.AddResourceWithActionParams(endpoint, data, "allocate_license")

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
