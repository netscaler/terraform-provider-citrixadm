package citrixadm

import (
	"context"
	"fmt"
	"log"

	"terraform-provider-citrixadm/service"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceConfigJobTemplate() *schema.Resource {
	return &schema.Resource{
		Description: "Get a configuration job template by name",
		ReadContext: dataSourceConfigJobTemplateRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Configuration Job Template Name",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func dataSourceConfigJobTemplateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In dataSourceConfigJobTemplateRead")
	var diags diag.Diagnostics
	c := m.(*service.NitroClient)

	endpoint := "configuration_template"

	configJobTemplateName := d.Get("name").(string)

	returnData, err := c.GetAllResource(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}

	// Find the correct resource with the given name and store config_job_template_id from the object
	for _, v := range returnData[endpoint].([]interface{}) {
		configJob := v.(map[string]interface{})
		if configJob["name"].(string) == configJobTemplateName {
			d.SetId(configJob["id"].(string))
			d.Set("name", configJob["name"].(string))
			break
		}
	}
	if d.Id() == "" {
		return diag.FromErr(fmt.Errorf("Configuration Job Template with name %s not found", configJobTemplateName))
	}
	return diags
}
