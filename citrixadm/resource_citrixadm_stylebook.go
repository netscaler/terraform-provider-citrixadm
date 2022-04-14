package citrixadm

import (
	"context"
	"fmt"
	"log"

	"terraform-provider-citrixadm/service"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceStylebook() *schema.Resource {
	return &schema.Resource{
		Description:   "Uploads (Imports) a new StyleBook",
		CreateContext: resourceStylebookCreate,
		ReadContext:   resourceStylebookRead,
		UpdateContext: resourceStylebookUpdate,
		DeleteContext: resourceStylebookDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Name of the StyleBook",
				Type:        schema.TypeString,
				Required:    true,
			},
			"namespace": {
				Description: "Namespace of the StyleBook",
				Type:        schema.TypeString,
				Required:    true,
			},
			"version": {
				Description: "Version of the StyleBook",
				Type:        schema.TypeString,
				Required:    true,
			},
			"source": {
				Description: "The YAML contents of the StyleBook",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func resourceStylebookCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceStylebookCreate")

	c := m.(*service.NitroClient)

	endpoint := "stylebooks"

	data := make(map[string]interface{})
	data["source"] = d.Get("source").(string)

	// POST: https://cocoa.adm.cloud.com/stylebook/nitro/v2/config/stylebooks/actions/import?mode=async
	_, err := c.AddResourceWithActionParams(endpoint, data, "import")

	if err != nil {
		return diag.Errorf("unable to upload Stylebook: %s", err.Error())
	}

	// ID is {namespace}/{version}/{name}
	id := fmt.Sprintf("%s/%s/%s", d.Get("namespace").(string), d.Get("version").(string), d.Get("name").(string))
	d.SetId(id)
	return resourceStylebookRead(ctx, d, m)
}

func resourceStylebookRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceStylebookRead")
	var diags diag.Diagnostics
	c := m.(*service.NitroClient)

	resourceID := d.Id()
	endpoint := "stylebooks"

	// GET URL: https://cocoa.adm.cloud.com/stylebook/nitro/v2/config/stylebooks/{namespace}/{version}/{name}
	returnData, err := c.GetResource(endpoint, resourceID)
	if err != nil {
		return diag.FromErr(err)
	}

	if bodyResource, ok := service.URLResourceToBodyResource[endpoint]; ok {
		endpoint = bodyResource
	}
	getResponseData := returnData[endpoint].(map[string]interface{})

	d.Set("name", getResponseData["name"].(string))
	d.Set("namespace", getResponseData["namespace"].(string))
	d.Set("version", getResponseData["version"].(string))

	return diags
}

func resourceStylebookUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceStylebookUpdate")
	c := m.(*service.NitroClient)

	resourceID := d.Id()
	endpoint := "stylebooks"

	data := make(map[string]interface{})
	data["source"] = d.Get("source").(string)

	updateResource := fmt.Sprintf("%s/%s", endpoint, resourceID)
	_, err := c.AddResourceWithActionParams(updateResource, data, "update")

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceStylebookRead(ctx, d, m)
}

func resourceStylebookDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceStylebookDelete")
	var diags diag.Diagnostics

	c := m.(*service.NitroClient)

	endpoint := "stylebooks"
	resourceID := d.Id()

	// DELETE: https://cocoa.adm.cloud.com/stylebook/nitro/v2/config/stylebooks/com.example.stylebooks/0.1/basic-lb-config
	_, err := c.DeleteResource(endpoint, resourceID)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
