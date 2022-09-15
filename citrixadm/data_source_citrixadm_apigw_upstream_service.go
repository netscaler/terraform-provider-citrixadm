package citrixadm

import (
	"context"
	"fmt"
	"log"

	"terraform-provider-citrixadm/service"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceApiGwUpstreamService() *schema.Resource {
	return &schema.Resource{
		Description: "Get a Upstream Service Id by name",
		ReadContext: dataSourceApiGwUpstreamServiceRead,
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

func dataSourceApiGwUpstreamServiceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In dataSourceApiGwUpstreamServiceRead")
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
	// Find the correct resource with the given name and store id from the object
	for _, v := range returnData[endpoint].([]interface{}) {
		apigwUpstreamService := v.(map[string]interface{})
		if apigwUpstreamService["name"].(string) == upstreamServiceName {
			d.SetId(apigwUpstreamService["id"].(string))
			d.Set("name", apigwUpstreamService["name"].(string))
			break
		}
	}
	if d.Id() == "" {
		return diag.FromErr(fmt.Errorf("ApiGW Upstream Service with name %s not found", upstreamServiceName))
	}
	return diags
}
