package citrixadm

import (
	"context"
	"fmt"
	"log"

	"terraform-provider-citrixadm/service"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceApiGwDefinition() *schema.Resource {
	return &schema.Resource{
		Description: "Get an API Definition Id by name",
		ReadContext: dataSourceApiGwDefinitionRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Name of an API Definition instance",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func dataSourceApiGwDefinitionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In dataSourceApiGwDefinitionRead")
	var diags diag.Diagnostics
	c := m.(*service.NitroClient)

	endpoint := "apidefs"

	DefinitionName := d.Get("name").(string)

	returnData, err := c.GetAllResource(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}

	// Find the correct resource with the given name and store Definition_id from the object
	for _, v := range returnData[endpoint].([]interface{}) {
		apigwDefinition := v.(map[string]interface{})
		if apigwDefinition["name"].(string) == DefinitionName {
			d.SetId(apigwDefinition["id"].(string))
			d.Set("name", apigwDefinition["name"].(string))
			break
		}
	}
	if d.Id() == "" {
		return diag.FromErr(fmt.Errorf("ApiGW Definition with name %s not found", DefinitionName))
	}
	return diags
}
