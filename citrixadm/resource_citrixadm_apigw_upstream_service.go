package citrixadm

import (
	"context"
	"log"

	"terraform-provider-citrixadm/service"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApiGwUpstreamService() *schema.Resource {
	return &schema.Resource{
		Description:   "Configure Upstream service for the provided API Deployment Id",
		CreateContext: resourceApiGwUpstreamServiceCreate,
		ReadContext:   resourceApiGwUpstreamServiceRead,
		UpdateContext: resourceApiGwUpstreamServiceUpdate,
		DeleteContext: resourceApiGwUpstreamServiceDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Upstream Service Name",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"service_fqdn": {
				Description: "Hostname of Upstream Service as FQDN where API traffic will be sent",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"service_fqdn_port": {
				Description: "Optional Listening port for Service FQDN",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
			"deployment_id": {
				Description: "API Deployment Id",
				Type:        schema.TypeString,
				Required:    true,
				Computed:    false,
				ForceNew:    true,
			},
			"scheme": {
				Description: "Protocol used to exchange data with the service",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"backend_servers": {
				Description: "List of all IPv4 or IPv6 host address and Port for Upstream Service.",
				Type:        schema.TypeList,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip4_addr": {
							Description: "IPv4 host address for the Upstream Server/Service where API traffic will be sent",
							Type:        schema.TypeString,
							Required:    true,
						},
						"port": {
							Description: "Port number for Upstream Server/Service",
							Type:        schema.TypeInt,
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func getApiGwUpstreamServicePayload(d *schema.ResourceData) interface{} {
	data := make(map[string]interface{})

	data["name"] = d.Get("name").(string)
	if v, ok := d.GetOk("service_fqdn"); ok {
		data["service_fqdn"] = v.(string)
	}
	if v, ok := d.GetOk("service_fqdn_port"); ok {
		data["service_fqdn_port"] = v.(int)
	}
	if v, ok := d.GetOk("scheme"); ok {
		data["scheme"] = v.(string)
	}
	if v, ok := d.GetOk("backend_servers"); ok {
		data["backend_servers"] = v.(interface{})
	}

	return data

}

func resourceApiGwUpstreamServiceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceApiGwUpstreamServiceCreate")

	c := m.(*service.NitroClient)

	endpoint := "upstreamservices"
	parent_name := "deployments"
	parent_id := d.Get("deployment_id").(string)

	returnData, err := c.AddChildResource(endpoint, getApiGwUpstreamServicePayload(d), parent_name, parent_id)
	if err != nil {
		return diag.Errorf("error creating upstream_service: %s", err.Error())
	}

	resourceID := returnData["upstreamservice"].(map[string]interface{})["id"].(string)
	log.Printf("id %s", resourceID)

	d.SetId(resourceID)

	return resourceApiGwUpstreamServiceRead(ctx, d, m)
}

func resourceApiGwUpstreamServiceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceApiGwUpstreamServiceRead")
	var diags diag.Diagnostics
	c := m.(*service.NitroClient)

	resourceID := d.Id()
	endpoint := "upstreamservices"

	returnData, err := c.GetResource(endpoint, resourceID)
	if err != nil {
		return diag.Errorf("error reading upstream_service %s: %s", resourceID, err.Error())
	}

	getResponseData := returnData["upstreamservice"].(map[string]interface{})

	log.Println("getResponseData", getResponseData)

	
	if _, ok := d.GetOk("service_fqdn"); ok {
		d.Set("service_fqdn", getResponseData["service_fqdn"].(string))
	}
	if _, ok := d.GetOk("service_fqdn_port"); ok {
		d.Set("service_fqdn_port", getResponseData["service_fqdn_port"].(float64))
	}
	d.Set("scheme", getResponseData["scheme"].(string))
	d.Set("name", getResponseData["name"].(string))
	d.Set("backend_servers", getResponseData["backend_servers"].([]interface{}))

	return diags
}

func resourceApiGwUpstreamServiceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceApiGwUpstreamServiceUpdate")
	c := m.(*service.NitroClient)

	resourceID := d.Id()
	endpoint := "upstreamservices"
	get_payload := getApiGwUpstreamServicePayload(d).(map[string](interface{}))
	get_payload["id"] = resourceID
	_, err := c.UpdateResource(endpoint, get_payload, resourceID)

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceApiGwUpstreamServiceRead(ctx, d, m)

}

func resourceApiGwUpstreamServiceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceApiGwUpstreamServiceDelete")
	var diags diag.Diagnostics

	c := m.(*service.NitroClient)

	endpoint := "upstreamservices"
	resourceID := d.Id()
	_, err := c.DeleteResource(endpoint, resourceID)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
