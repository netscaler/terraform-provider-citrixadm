package citrixadm

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"terraform-provider-citrixadm/service"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApiGwDeployment() *schema.Resource {
	return &schema.Resource{
		Description:   "Configure API Proxy endpoint on API Gateway",
		CreateContext: resourceApiGwDeploymentCreate,
		ReadContext:   resourceApiGwDeploymentRead,
		UpdateContext: resourceApiGwDeploymentUpdate,
		DeleteContext: resourceApiGwDeploymentDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "API Deployment Name",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"api_id": {
				Description: "API Definition Id",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"deploy": {
				Description: "Flag to deply/undeploy the API Deployment Resource",
				Type:        schema.TypeBool,
				Required:    true,
			},
			"target_apigw": {
				Description: "Details of Target ADC instance acting as API Gateway",
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "ADC Instance Id",
							Type:        schema.TypeString,
							Required:    true,
						},
						"display_name": {
							Description: "ADC Instance Name(IP-Address)",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
			},
			"tags": {
				Description: "Tags",
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"apiproxy_ref": {
				Description: "API Proxy Reference",
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "API proxy ID",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
			},
			"upstreamservices": {
				Description: "List of all Upstream Service",
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "Upstream Service Name",
							Type:        schema.TypeString,
							Required:    true,
						},
						"service_fqdn": {
							Description: "Hostname of Upstream Service as FQDN where API traffic will be sent",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"service_fqdn_port": {
							Description: "Optional Listening port for Service FQDN",
							Type:        schema.TypeInt,
							Optional:    true,
						},
						"scheme": {
							Description: "Protocol used to exchange data with the service",
							Type:        schema.TypeString,
							Optional:    true,
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
				},
			},
			"routes": {
				Description: "List of all API Routes",
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "Name for API Route",
							Type:        schema.TypeString,
							Required:    true,
						},
						"route_param": {
							Description: "API Resource Path in the API definition or user entered Path.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"route_paramtype": {
							Description: "API Route Param type",
							Type:        schema.TypeString,
							Required:    true,
						},
						"upstreamservice_name": {
							Description: "Upstream Service Name where API traffic will be sent",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func getApiGwDeploymentPayload(d *schema.ResourceData) interface{} {
	data := make(map[string]interface{})

	deployment := make(map[string]interface{})
	deployment["name"] = d.Get("name").(string)
	if v, ok := d.GetOk("api_id"); ok {
		deployment["api_id"] = v.(string)
	}
	if v, ok := d.GetOk("tags"); ok {
		deployment["tags"] = v.([]interface{})
	}
	if v, ok := d.GetOk("apiproxy_ref"); ok {
		deployment["apiproxy_ref"] = v.([]interface{})[0].(map[string]interface{})

	}

	a := d.Get("target_apigw").([]interface{})

	deployment["target_apigw"] = mapToString(a[0].(map[string]interface{}))

	data["deployment"] = deployment

	if _, ok := d.GetOk("routes"); ok {
		data["routes"] = d.Get("routes").([]interface{})
	}
	if _, ok := d.GetOk("upstreamservices"); ok {
		data["upstreamservices"] = d.Get("upstreamservices").([]interface{})
	}

	return data

}

func resourceApiGwDeploymentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceApiGwDeploymentCreate")

	c := m.(*service.NitroClient)

	endpoint := "deployments"
	parentName := "apidefs"
	parentId := d.Get("api_id").(string)

	returnData, err := c.AddChildResource(endpoint, getApiGwDeploymentPayload(d), parentName, parentId)

	if err != nil {
		return diag.Errorf("error creating apigw_deployment: %s", err.Error())
	}

	resourceID := returnData["deployment"].(map[string]interface{})["id"].(string)
	log.Printf("id %s", resourceID)

	d.SetId(resourceID)

	if d.Get("deploy").(bool) == true {
		_, err := c.AddResourceWithActionParams(endpoint, resourceID, "applyconfig")
		if err != nil {
			return diag.FromErr(err)
		}
		err1 := c.WaitForDeplymentCompletion("deploymentsdeploy", resourceID, time.Duration(c.ActivityTimeout)*time.Second, "applyconfig")
		if err1 != nil {
			return diag.FromErr(err1)
		}

	}

	return resourceApiGwDeploymentRead(ctx, d, m)
}

func resourceApiGwDeploymentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceApiGwDeploymentRead")
	var diags diag.Diagnostics
	c := m.(*service.NitroClient)

	resourceID := d.Id()
	endpoint := "deployments"

	returnData, err := c.GetResource(endpoint, resourceID)
	if err != nil {
		return diag.Errorf("error reading apigw_deployment %s: %s", resourceID, err.Error())
	}

	getResponseData := returnData["deployment"].(map[string]interface{})

	returnData1, err1 := c.GetResource("deploymentsdeploy", resourceID)
	if err1 != nil {
		return diag.Errorf("error Reading apigw deploy state %s", err1)
	}
	status := returnData1["configstatus"].(map[string]interface{})["status"].(string)
	if status == "Applied" {
		d.Set("deploy", true)
	} else if status == "Indraft" {
		d.Set("deploy", false)
	}
	log.Println("getResponseData", getResponseData)
	log.Printf("The value is %v and %T", getResponseData["tags"], getResponseData["tags"])

	d.Set("api_id", getResponseData["api_id"].(string))
	d.Set("name", getResponseData["name"].(string))

	// if the tags is given as empty list we receieve empty string
	if getResponseData["tags"] == "" {
		tag := make([]interface{}, 0)
		d.Set("tags", tag)
	} else {
		d.Set("tags", getResponseData["tags"].([]interface{}))
	}
	d.Set("target_apigw", flattenApiGwDeploymentAttribute(getResponseData["target_apigw"].(string)))
	d.Set("apiproxy_ref", []interface{}{(getResponseData["apiproxy_ref"].(map[string]interface{}))})

	return diags
}

func flattenApiGwDeploymentAttribute(attribute string) []interface{} {
	a := make(map[string]interface{})

	b := attribute
	//To covert jsonstring(b) to map(a)
	json.Unmarshal([]byte(b), &a)

	return []interface{}{a}
}

func resourceApiGwDeploymentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceApiGwDeploymentUpdate")
	c := m.(*service.NitroClient)

	resourceID := d.Id()
	endpoint := "deployments"

	returnData, err := c.GetResource("deploymentsdeploy", resourceID)
	if err != nil {
		return diag.Errorf("error Reading apigw deployment deploy status %s", err)
	}
	status := returnData["configstatus"].(map[string]interface{})["status"].(string)
	// To check if the deployment is undeployed and user wants to deploy i.e., applyconfig action
	if status == "Indraft" && d.Get("deploy").(bool) == true {
		_, err := c.AddResourceWithActionParams(endpoint, resourceID, "applyconfig")
		if err != nil {
			return diag.FromErr(err)
		}
		err1 := c.WaitForDeplymentCompletion("deploymentsdeploy", resourceID, time.Duration(c.ActivityTimeout)*time.Second, "applyconfig")
		if err1 != nil {
			return diag.FromErr(err1)
		}
		// If the deployment is deployed and user wants to Undeploy it i.e, undeployconfig action
	} else if status == "Applied" && d.Get("deploy").(bool) == false {
		_, err := c.AddResourceWithActionParams(endpoint, resourceID, "undeployconfig")
		if err != nil {
			return diag.FromErr(err)
		}
		err1 := c.WaitForDeplymentCompletion("deploymentsdeploy", resourceID, time.Duration(c.ActivityTimeout)*time.Second, "undeployconfig")
		if err1 != nil {
			return diag.FromErr(err1)
		}
	} else {
		// nothing
	}

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceApiGwDeploymentRead(ctx, d, m)

}

func resourceApiGwDeploymentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceApiGwDeploymentDelete")
	var diags diag.Diagnostics

	c := m.(*service.NitroClient)

	endpoint := "deployments"
	resourceID := d.Id()
	_, err := c.DeleteResource(endpoint, resourceID)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
