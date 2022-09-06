package citrixadm

import (
	"context"
	"fmt"
	"log"

	"terraform-provider-citrixadm/service"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCertStore() *schema.Resource {
	return &schema.Resource{
		Description: "Get a Certstore Instance Id by name",
		ReadContext: dataSourceCertStoreRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Name of the Certstore instance",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func dataSourceCertStoreRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In dataSourceCertStoreRead")
	var diags diag.Diagnostics
	c := m.(*service.NitroClient)

	endpoint := "cert_store"

	certstoreName := d.Get("name").(string)

	returnData, err := c.GetAllResource(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}

	// Find the correct resource with the given name and store deployment_id from the object
	for _, v := range returnData[endpoint].([]interface{}) {
		CertStore := v.(map[string]interface{})
		if CertStore["name"].(string) == certstoreName {
			d.SetId(CertStore["id"].(string))
			d.Set("name", CertStore["name"].(string))
			break
		}
	}
	if d.Id() == "" {
		return diag.FromErr(fmt.Errorf("ApiGW Certstore with name %s not found", certstoreName))
	}
	return diags
}
