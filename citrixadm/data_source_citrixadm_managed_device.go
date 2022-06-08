package citrixadm

import (
	"context"
	"log"

	"terraform-provider-citrixadm/service"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceManagedDevice() *schema.Resource {
	return &schema.Resource{
		Description: "Get a managed device ID by IP address",
		ReadContext: dataSourceManagedDeviceRead,
		Schema: map[string]*schema.Schema{
			"ip_address": {
				Description: "IP Address for this managed device",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func dataSourceManagedDeviceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In dataSourceManagedDeviceRead")
	var diags diag.Diagnostics
	c := m.(*service.NitroClient)

	resourceID, err := getManagedDeviceID(c, d.Get("ip_address").(string))

	if err != nil {
		return diag.Errorf("unable to get Managed Device ID: %s", err.Error())
	}
	d.SetId(resourceID)
	d.Set("ip_address", d.Get("ip_address").(string))

	return diags
}
