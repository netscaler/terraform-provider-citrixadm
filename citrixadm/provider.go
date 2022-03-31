package citrixadm

import (
	"context"
	"terraform-provider-citrixadm/service"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CITRIXADM_HOST", nil),
			},
			"client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CITRIXADM_CLIENT_ID", nil),
			},
			"client_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("CITRIXADM_CLIENT_SECRET", nil),
			},
			"host_location": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CITRIXADM_HOST_LOCATION", ""),
			},
			"customer_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CITRIXADM_CUSTOMER_ID", ""),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"citrixadm_managed_device":                  resourceManagedDevice(),
			"citrixadm_ns_device_profile":               resourceNsDeviceProfile(),
			"citrixadm_managed_device_allocate_license": resourceManagedDeviceAllocateLicense(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"citrixadm_mps_agent": dataSourceMpsAgent(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	tflog.Trace(ctx, "In providerConfigure")
	var diags diag.Diagnostics
	params := service.NitroParams{
		Host:         d.Get("host").(string),
		HostLocation: d.Get("host_location").(string),
		ID:           d.Get("client_id").(string),
		Secret:       d.Get("client_secret").(string),
		CustomerID:   d.Get("customer_id").(string),
	}
	c, err := service.NewNitroClientFromParams(params)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	if err != nil {
		return nil, diag.FromErr(err)
	}
	return c, diags
}
