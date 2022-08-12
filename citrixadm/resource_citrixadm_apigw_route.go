package citrixadm

import (
	"context"
	"log"

	"terraform-provider-citrixadm/service"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApiGwRoute() *schema.Resource {
	return &schema.Resource{
		Description:   "API Route for the provided API Deployment Id.",
		CreateContext: resourceApiGwRouteCreate,
		ReadContext:   resourceApiGwRouteRead,
		DeleteContext: resourceApiGwRouteDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Name for API Route.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"route_param": {
				Description: "API Resource Path in the API definition or user entered Path.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"route_paramtype": {
				Description: "API Route Param type",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"upstreamservice_name": {
				Description: "Upstream Service Name where API traffic will be sent",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"deployment_id": {
				Description: "API Deployment Id",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func getApiGwRoutePayload(d *schema.ResourceData) interface{} {
	data := make(map[string]interface{})

	data["name"] = d.Get("name").(string)

	data["route_param"] = d.Get("route_param").(string)

	data["route_paramtype"] = d.Get("route_paramtype").(string)

	data["upstreamservice_name"] = d.Get("upstreamservice_name").(string)

	return data

}

func resourceApiGwRouteCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceApiGwRouteCreate")

	c := m.(*service.NitroClient)

	endpoint := "routes"
	parent_Name := "deployments"
	parent_Id := d.Get("deployment_id").(string)

	returnData, err := c.AddChildResource(endpoint, getApiGwRoutePayload(d), parent_Name, parent_Id)
	if err != nil {
		return diag.Errorf("error creating api_route: %s", err.Error())
	}

	resourceID := returnData["route"].(map[string]interface{})["id"].(string)
	log.Printf("id %s", resourceID)

	d.SetId(resourceID)

	return resourceApiGwRouteRead(ctx, d, m)
}

func resourceApiGwRouteRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceApiGwRouteRead")
	var diags diag.Diagnostics
	c := m.(*service.NitroClient)

	resourceID := d.Id()
	endpoint := "routes"

	returnData, err := c.GetResource(endpoint, resourceID)
	if err != nil {
		return diag.Errorf("error reading Api_route %s: %s", resourceID, err.Error())
	}

	getResponseData := returnData["route"].(map[string]interface{})

	log.Println("getResponseData", getResponseData)

	d.Set("name", getResponseData["name"].(string))
	d.Set("route_param", getResponseData["route_param"].(string))
	d.Set("route_paramtype", getResponseData["route_paramtype"].(string))
	d.Set("upstreamservice_name", getResponseData["upstreamservice_name"].(string))

	return diags
}

func resourceApiGwRouteDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceApiGwRouteDelete")
	var diags diag.Diagnostics

	c := m.(*service.NitroClient)

	endpoint := "routes"
	resourceID := d.Id()
	_, err := c.DeleteResource(endpoint, resourceID)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
