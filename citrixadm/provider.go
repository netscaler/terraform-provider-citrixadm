package citrixadm

import (
	"context"
	"fmt"
	"strings"
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
				Description: "Citrix Adm host. Can be specified with `CITRIXADM_HOST` environment variable. This has to start with https://",
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					// if the value does not start with http, throw an error
					if !strings.HasPrefix(v.(string), "https://") {
						errors = append(errors, fmt.Errorf("host must start with https://"))
					}
					return
				},
			},
			"client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CITRIXADM_CLIENT_ID", nil),
				Description: "Citrix Adm client id. Can be specified with `CITRIXADM_CLIENT_ID` environment variable.",
			},
			"client_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("CITRIXADM_CLIENT_SECRET", nil),
				Description: "Citrix Adm client secret. Can be specified with `CITRIXADM_CLIENT_SECRET` environment variable.",
			},
			"host_location": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CITRIXADM_HOST_LOCATION", ""),
				Description: "Citrix Adm host location, e.g. `us`, `eu`. Can be specified with `CITRIXADM_HOST_LOCATION` environment variable.",
			},
			"customer_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CITRIXADM_CUSTOMER_ID", ""),
				Description: "Citrix Adm customer/tenant id. Can be specified with `CITRIXADM_CUSTOMER_ID` environment variable.",
			},
			"fail_on_stall": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CITRIXADM_FAIL_ON_STALL", false),
				Description: "Boolean flag. Set to true for the module to fail when a status of job stalled is reported. Can be specified with `CITRIXADM_FAIL_ON_STALL` environment variable.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"citrixadm_managed_device":                  resourceManagedDevice(),
			"citrixadm_ns_device_profile":               resourceNsDeviceProfile(),
			"citrixadm_managed_device_allocate_license": resourceManagedDeviceAllocateLicense(),
			"citrixadm_stylebook_configpack":            resourceStylebookConfigpack(),
			"citrixadm_stylebook":                       resourceStylebook(),
			"citrixadm_provisioning_profile":            resourceProvisioningProfile(),
			"citrixadm_provision_vpx":                   resourceProvisionVpx(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"citrixadm_mps_agent":           dataSourceMpsAgent(),
			"citrixadm_managed_device":      dataSourceManagedDevice(),
			"citrixadm_config_job_template": dataSourceConfigJobTemplate(),
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

	return c, diags
}
