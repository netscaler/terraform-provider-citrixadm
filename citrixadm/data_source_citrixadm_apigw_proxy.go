package citrixadm

import (
	"context"
	"fmt"
	"log"

	"terraform-provider-citrixadm/service"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceApiGwProxy() *schema.Resource {
	return &schema.Resource{
		Description: "Get an API Proxy Id by name",
		ReadContext: dataSourceApiGwProxyRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Name of the API Proxy instance",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func dataSourceApiGwProxyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In dataSourceApiGwProxyRead")
	var diags diag.Diagnostics
	c := m.(*service.NitroClient)

	endpoint := "apiproxies"

	apigwproxyName := d.Get("name").(string)

	returnData, err := c.GetAllResource(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}

	// Find the correct resource with the given name and store apigwproxy_id from the object
	for _, v := range returnData["apiproxies"].([]interface{}) {
		ApiGwProxy := v.(map[string]interface{})
		if ApiGwProxy["proxy_name"].(string) == apigwproxyName {
			d.SetId(ApiGwProxy["id"].(string))
			d.Set("name", ApiGwProxy["proxy_name"].(string))
			break
		}
	}
	if d.Id() == "" {
		return diag.FromErr(fmt.Errorf("ApiGW apigwproxy with name %s not found", apigwproxyName))
	}
	return diags
}
