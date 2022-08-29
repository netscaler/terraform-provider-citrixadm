package citrixadm

import (
	"context"
	"fmt"
	"log"

	"terraform-provider-citrixadm/service"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceApigwUpstreamService() *schema.Resource {
	return &schema.Resource{
		Description: "Get a Upstream Service Id by name",
		ReadContext: dataSourceApigwUpstreamServiceRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Name of the Upstream Service instance",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"deployment_id": {
				Description: "Deployment Id under which this upstreamService exist",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func dataSourceApigwUpstreamServiceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In dataSourceApigwUpstreamServiceRead")
	var diags diag.Diagnostics
	c := m.(*service.NitroClient)

	endpoint := "upstreamservices"
	parentName := "deployments"
	parentId := d.Get("deployment_id").(string)

	upstreamServiceName := d.Get("name").(string)

	returnData, err := c.GetAllChildResource(endpoint, parentName, parentId)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("Datatata %v",returnData)
	// Find the correct resource with the given name and store id from the object
	for _, v := range returnData[endpoint].([]interface{}) {
		ApigwUpstreamService := v.(map[string]interface{})
		if ApigwUpstreamService["name"].(string) == upstreamServiceName {
			d.SetId(ApigwUpstreamService["id"].(string))
			d.Set("name", ApigwUpstreamService["name"].(string))
			break
		}
	}
	if d.Id() == "" {
		return diag.FromErr(fmt.Errorf("ApiGW Upstream Service with name %s not found", upstreamServiceName))
	}
	return diags
}
