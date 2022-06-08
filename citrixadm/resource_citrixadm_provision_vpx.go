package citrixadm

import (
	"context"
	"errors"
	"log"

	"terraform-provider-citrixadm/service"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceProvisionVpx() *schema.Resource {
	return &schema.Resource{
		Description:   "Configuration for Managed Device resource",
		CreateContext: resourceProvisionVpxCreate,
		ReadContext:   resourceProvisionVpxRead,
		// UpdateContext: resourceProvisionVpxUpdate,
		// DeleteContext: resourceProvisionVpxDelete,
		DeleteContext: schema.NoopContext,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the instance to be provisioned.",
				Required:    true,
				ForceNew:    true,
			},
			"provisioning_profile_id": {
				Type:        schema.TypeString,
				Description: "Provisioning Profile Id used to provision instance.",
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func getProvisionVpxPayload(d *schema.ResourceData) interface{} {
	data := make(map[string]interface{})
	data["name"] = d.Get("name").(string)
	data["provisioning_profile_id"] = d.Get("provisioning_profile_id").(string)
	return data
}

func resourceProvisionVpxCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceProvisionVpxCreate")

	c := m.(*service.NitroClient)

	endpoint := "instances"

	returnData, err := c.AddResource(endpoint, getProvisionVpxPayload(d))

	if err != nil {
		return diag.Errorf("unable to Provisoin VPX: %s", err.Error())
	}

	// if bodyResource, ok := service.URLResourceToBodyResource[endpoint]; ok {
	// 	endpoint = bodyResource
	// }

	jobID := returnData["instance"].(map[string]interface{})["job_id"].(string)

	// Wait for the job to complete
	log.Printf("Waiting for the job to complete")
	err = c.WaitForProvisioningJobCompletion(jobID)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceID, err := getProvisionVpxID(c, d.Get("name").(string))

	if err != nil {
		return diag.FromErr(errors.New("Failed to find resource ID"))
	}

	d.SetId(resourceID)
	return resourceProvisionVpxRead(ctx, d, m)
}

func getProvisionVpxID(c *service.NitroClient, name string) (string, error) {
	endpoint := "instances"
	returnData, err := c.GetAllResource(endpoint)
	if err != nil {
		return "", err
	}

	for _, v := range returnData[endpoint].([]interface{}) {
		if v.(map[string]interface{})["name"].(string) == name {
			return v.(map[string]interface{})["id"].(string), nil
		}
	}
	return "", errors.New("Failed to find instance resource ID with Host name: " + name)
}

func resourceProvisionVpxRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceProvisionVpxRead")
	var diags diag.Diagnostics
	c := m.(*service.NitroClient)

	resourceID := d.Id()
	endpoint := "instances"

	returnData, err := c.GetResource(endpoint, resourceID)
	if err != nil {
		return diag.FromErr(err)
	}

	getResponseData := returnData[endpoint].([]interface{})[0].(map[string]interface{})

	d.Set("name", getResponseData["name"].(string))
	d.Set("provisioning_profile_id", getResponseData["provisioning_profile_id"].(string))

	return diags
}

// func resourceProvisionVpxDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
// 	log.Printf("In resourceProvisionVpxDelete")
// 	var diags diag.Diagnostics

// 	c := m.(*service.NitroClient)

// 	endpoint := "instances"
// 	resourceID := d.Id()

// 	returnData, err := c.DeleteResource(endpoint, resourceID)

// 	// if bodyResource, ok := service.URLResourceToBodyResource[endpoint]; ok {
// 	// 	endpoint = bodyResource
// 	// }

// 	jobID := returnData["instance"].(map[string]interface{})["job_id"].(string)

// 	// Wait for the delete job to complete
// 	log.Printf("Waiting for the delete job to complete")
// 	err = c.WaitForProvisioningJobCompletion(jobID)

// 	if err != nil {
// 		return diag.FromErr(err)
// 	}

// 	d.SetId("")

// 	return diags
// }
