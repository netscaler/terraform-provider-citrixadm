package citrixadm

import (
	"context"
	"log"
	"time"

	"terraform-provider-citrixadm/service"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceStylebookConfigpack() *schema.Resource {
	return &schema.Resource{
		Description:   "Configuration for Stylebook Configpack resource",
		CreateContext: resourceStylebookConfigpackCreate,
		ReadContext:   resourceStylebookConfigpackRead,
		UpdateContext: resourceStylebookConfigpackUpdate,
		DeleteContext: resourceStylebookConfigpackDelete,
		Schema: map[string]*schema.Schema{
			// FIXME: George: terraform plan creates unnecessary diffs
			// https://github.com/hashicorp/terraform-plugin-sdk/issues/477
			// FIXME: George: not implementing upgrade, should it be ForceNew = true?
			"parameters": {
				Description: "A JSON dictionary containing the values for the Parameters of the StyleBook, where the key of each item in the dictionary is the name of the parameter and the value is the value of the parameter (note that the value can be an arbitrary JSON object depending on the type of the parameter (refer to the StyleBook schema).",
				Type:        schema.TypeMap,
				Required:    true,
				Elem: &schema.Schema{
					Type:     schema.TypeString,
					Computed: true,
				},
				// DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				// 	return old == new
				// },
			},
			"stylebook": {
				Description: "The StyleBook to use for the managed device.",
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "The name of the StyleBook.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"namespace": {
							Description: "The namespace of the StyleBook.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"version": {
							Description: "The version of the StyleBook.",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
			},
			"targets": {
				Description: "A dictionary specifying the devices to which the configpack is applied. The key of each item in the dictionary is the device's IP address and the value is a dictionary that contains one item which corresponds to the devices's ID in ADM in the form 'id':'<id-value>'",
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Description: "The device's ID in ADM",
							Type:        schema.TypeString,
							Required:    true,
						},
						"instance_ip": {
							Description: "The device's IP address",
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func getStylebookConfigpackPayload(d *schema.ResourceData) interface{} {
	data := make(map[string]interface{})

	data["parameters"] = d.Get("parameters").(map[string]interface{})

	data["stylebook"] = d.Get("stylebook").([]interface{})[0].(map[string]interface{})

	if v, ok := d.GetOk("targets"); ok {
		data["targets"] = v.([]interface{})
	}

	if v, ok := d.GetOk("entity_tag"); ok {
		data["entity_tag"] = v.([]interface{})
	}

	return data
}

func resourceStylebookConfigpackCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceStylebookConfigpackCreate")

	c := m.(*service.NitroClient)

	endpoint := "configpacks"

	returnData, err := c.AddResource(endpoint, getStylebookConfigpackPayload(d))

	if err != nil {
		return diag.Errorf("unable to create Configpack: %s", err.Error())
	}

	// returnData
	// {
	// 	"job": {
	// 		"job_id": "1629958816"
	// 	}
	// }

	jobID := returnData["job"].(map[string]interface{})["job_id"].(string)

	// Wait for the job to complete
	log.Printf("Waiting for the job to complete")
	err = c.WaitForStylebookJobCompletion(jobID, time.Duration(c.ActivityTimeout)*time.Second)
	if err != nil {
		return diag.FromErr(err)
	}

	// JobID itself is the resource ID
	d.SetId(jobID)
	return resourceStylebookConfigpackRead(ctx, d, m)

}

func resourceStylebookConfigpackRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceStylebookConfigpackRead")
	var diags diag.Diagnostics
	c := m.(*service.NitroClient)

	resourceID := d.Id()
	endpoint := "configpacks"

	returnData, err := c.GetResource(endpoint, resourceID)
	if err != nil {
		return diag.FromErr(err)
	}

	if bodyResource, ok := service.URLResourceToBodyResource[endpoint]; ok {
		endpoint = bodyResource
	}
	getResponseData := returnData[endpoint].(map[string]interface{})

	// Update the state with the returned data
	d.Set("parameters", getResponseData["parameters"])
	d.Set("stylebook", getResponseData["stylebook"])
	d.Set("targets", getResponseData["targets"])

	return diags
}

func resourceStylebookConfigpackUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceStylebookConfigpackUpdate")
	c := m.(*service.NitroClient)

	resourceID := d.Id()
	endpoint := "configpacks"

	returnData, err := c.UpdateResource(endpoint, getStylebookConfigpackPayload(d), resourceID)

	if err != nil {
		return diag.FromErr(err)
	}

	jobID := returnData["job"].(map[string]interface{})["job_id"].(string)

	// Wait for the job to complete
	log.Printf("Waiting for the job to complete")
	err = c.WaitForStylebookJobCompletion(jobID, time.Duration(c.ActivityTimeout)*time.Second)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceStylebookConfigpackRead(ctx, d, m)
}

func resourceStylebookConfigpackDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceStylebookConfigpackDelete")
	var diags diag.Diagnostics

	c := m.(*service.NitroClient)

	endpoint := "configpacks"
	resourceID := d.Id()

	returnData, err := c.DeleteResource(endpoint, resourceID)

	if err != nil {
		return diag.FromErr(err)
	}

	jobID := returnData["job"].(map[string]interface{})["job_id"].(string)

	// Wait for the job to complete
	log.Printf("Waiting for the job to complete")
	err = c.WaitForStylebookJobCompletion(jobID, time.Duration(c.ActivityTimeout)*time.Second)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
