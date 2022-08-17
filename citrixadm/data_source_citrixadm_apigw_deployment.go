package citrixadm

import (
	"context"
	"fmt"
	"log"

	"terraform-provider-citrixadm/service"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceApigwDeployment() *schema.Resource {
	return &schema.Resource{
		Description: "Get a Deployment Id by name",
		ReadContext: dataSourceApigwDeploymentRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Name of the Deployment instance",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func dataSourceApigwDeploymentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In dataSourceApigwDeploymentRead")
	var diags diag.Diagnostics
	c := m.(*service.NitroClient)

	endpoint := "deployments"

	deploymentName := d.Get("name").(string)

	returnData, err := c.GetAllResource(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}

	// Find the correct resource with the given name and store deployment_id from the object
	for _, v := range returnData[endpoint].([]interface{}) {
		apigwDeployment := v.(map[string]interface{})
		if apigwDeployment["name"].(string) == deploymentName {
			d.SetId(apigwDeployment["id"].(string))
			d.Set("name", apigwDeployment["name"].(string))
			break
		}
	}
	if d.Id() == "" {
		return diag.FromErr(fmt.Errorf("ApiGW Deployment with name %s not found", deploymentName))
	}
	return diags
}
