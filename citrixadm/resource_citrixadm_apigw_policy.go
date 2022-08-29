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

func resourceApiGwPolicy() *schema.Resource {
	return &schema.Resource{
		Description:   "Configure API Policies for the provided API Deployment Id",
		CreateContext: resourceApiGwPolicyCreate,
		ReadContext:   resourceApiGwPolicyRead,
		UpdateContext: resourceApiGwPolicyUpdate,
		DeleteContext: resourceApiGwPolicyDelete,
		Schema: map[string]*schema.Schema{
			"policygroup_name": {
				Description: "API Policy Group Name",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"deployment_id": {
				Description: "Deployment Id",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"deploy": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"upstreamservice_id": {
				Description: "Upstream Service Id in the API Deployment",
				Type:        schema.TypeString,
				Required:    true,
			},
			"requestpath": {
				Description: "List of API Policy with corresponding JSON configuration to be applied on Request Path",
				Type:        schema.TypeList,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policytype": {
							Description: "Type for API Policy",
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
						},
						"policy_name": {
							Description: "API Policy name",
							Type:        schema.TypeString,
							Required:    true,
						},
						"order_index": {
							Description: "Order in which to apply API Policy",
							Type:        schema.TypeInt,
							Required:    true,
						},
						"config_spec": {
							Description: "JSON dictionary of configuration parameters and values as required for the selected API Policy type. Please refer to the document for details on config_spec JSON schema for different API policy type",
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"api_resource_paths": {
										Description: "API resource and its methods for which you want to apply a policy",
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"delete": {
													Description: "DELETE Method",
													Type:        schema.TypeBool,
													Optional:    true,
												},
												"endpoints": {
													Description: "API Path",
													Type:        schema.TypeList,
													Elem:        &schema.Schema{Type: schema.TypeString},
													Optional:    true,
												},
												"get": {
													Description: "GET Method",
													Type:        schema.TypeBool,
													Optional:    true,
												},
												"patch": {
													Description: "PATCH Method",
													Type:        schema.TypeBool,
													Optional:    true,
												},
												"post": {
													Description: "POST Method",
													Type:        schema.TypeBool,
													Optional:    true,
												},
												"put": {
													Description: "PUT Method",
													Type:        schema.TypeBool,
													Optional:    true,
												},
											},
										},
									},
									"custom_rules": {
										Description: "Custom path prefixes and multiple methods",
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"delete": {
													Description: "DELETE Method",
													Type:        schema.TypeBool,
													Optional:    true,
												},
												"endpoints": {
													Description: "API Path",
													Type:        schema.TypeList,
													Elem:        &schema.Schema{Type: schema.TypeString},
													Optional:    true,
												},
												"get": {
													Description: "GET Method",
													Type:        schema.TypeBool,
													Optional:    true,
												},
												"patch": {
													Description: "PATCH Method",
													Type:        schema.TypeBool,
													Optional:    true,
												},
												"post": {
													Description: "POST Method",
													Type:        schema.TypeBool,
													Optional:    true,
												},
												"put": {
													Description: "PUT Method",
													Type:        schema.TypeBool,
													Optional:    true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func getApiGwPolicyPayload(d *schema.ResourceData) interface{} {
	data := make(map[string]interface{})
	// requestpath_list := make([]interface{},2)
	requestpath_list := make([]interface{}, 0)
	requestpath := make(map[string]interface{})

	config_spec := make(map[string]interface{})

	data["policygroup_name"] = d.Get("policygroup_name").(string)
	log.Printf("%v and %v and %v and %v and %v and %v", d.Get("policytype"), d.Get("policy_name"), d.Get("order_index"), d.Get("config_spec"), d.Get("endpoints"), d.Get("put"))
	if v, ok := d.GetOk("upstreamservice_id"); ok {
		data["upstreamservice_id"] = v.(string)
	}
	//iterate through the requestpath list and store the same in the list, and convert the internal list to map
	//this is done to store the value of api_resource_paths as api-resource-paths
	for i := range d.Get("requestpath").([]interface{}) {

		if v, ok := d.GetOk("requestpath"); ok {
			requestpath["policytype"] = v.([]interface{})[i].(map[string]interface{})["policytype"].(string)
			requestpath["policy_name"] = v.([]interface{})[i].(map[string]interface{})["policy_name"].(string)
			requestpath["order_index"] = v.([]interface{})[i].(map[string]interface{})["order_index"].(int)

			//To convert the list to map
			config_spec["api-resource-paths"] = v.([]interface{})[i].(map[string]interface{})["config_spec"].([]interface{})[0].(map[string]interface{})["api_resource_paths"].([]interface{})[0]

			config_spec["custom-rules"] = v.([]interface{})[i].(map[string]interface{})["config_spec"].([]interface{})[0].(map[string]interface{})["custom_rules"].([]interface{})[0]
			requestpath["config_spec"] = config_spec

		}
		requestpath_list = append(requestpath_list, requestpath)

	}
	data["requestpath"] = requestpath_list
	// var payload []interface{}
	// payload = append(payload, data)

	return data
}
func toJSONIndent(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "    ")
	return string(b)
}

func resourceApiGwPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceApiGwPolicyCreate")

	c := m.(*service.NitroClient)

	endpoint := "policies"
	parentId := d.Get("deployment_id").(string)
	parentName := "deployments"
	returnData, err := c.AddChildResource(endpoint, getApiGwPolicyPayload(d), parentName, parentId)

	if err != nil {
		return diag.Errorf("error creating ApiGw Policies: %s", err.Error())
	}

	resourceID := returnData["policy"].(map[string]interface{})["id"].(string)
	log.Printf("id %s", resourceID)

	d.SetId(resourceID)

	if d.Get("deploy").(bool) == true {
		_, err := c.AddResourceWithActionParams(endpoint, resourceID, "applyconfig")
		if err != nil {
			return diag.FromErr(err)
		}
		err1 := c.WaitForDeplymentCompletion("policiesdeploy", resourceID, time.Duration(c.ActivityTimeout)*time.Second, "applyconfig")
		if err1 != nil {
			return diag.FromErr(err1)
		}

	}

	return resourceApiGwPolicyRead(ctx, d, m)
}

func resourceApiGwPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceApiGwPolicyRead")
	var diags diag.Diagnostics
	c := m.(*service.NitroClient)

	resourceID := d.Id()
	endpoint := "policies"

	returnData, err := c.GetResource(endpoint, resourceID)
	if err != nil {
		return diag.Errorf("error reading ApiGw Policies %s: %s", resourceID, err.Error())
	}

	getResponseData := returnData["policy"].(map[string]interface{})

	returnData1, err1 := c.GetResource("policiesdeploy", resourceID)
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

	d.Set("policygroup_name", getResponseData["policygroup_name"].(string))
	d.Set("upstreamservice_id", getResponseData["upstreamservice_id"].(string))
	d.Set("requestpath", flattenRequestPath(getResponseData["requestpath"].([]interface{})))

	return diags
}

func flattenRequestPath(requestpath []interface{}) []interface{} {
	flattenRequestList := make([]interface{}, 0)
	//Iterate through requestPath and Flatten the attributes accordingly...and also manage the names of attributes eg:api_resource_paths and custom_rules
	for i := range requestpath {
		config_spec := make([]interface{}, 0)
		requestpath_in := make(map[string]interface{})
		a := make(map[string]interface{})
		api := make([]interface{}, 0)
		api = append(api, requestpath[i].(map[string]interface{})["config_spec"].(map[string]interface{})["api-resource-paths"])
		a["api_resource_paths"] = api
		api1 := make([]interface{}, 0)
		api1 = append(api1, requestpath[i].(map[string]interface{})["config_spec"].(map[string]interface{})["custom-rules"])
		a["custom_rules"] = api1
		requestpath_in["policytype"] = requestpath[i].(map[string]interface{})["policytype"]
		requestpath_in["policy_name"] = requestpath[i].(map[string]interface{})["policy_name"]
		requestpath_in["order_index"] = requestpath[i].(map[string]interface{})["order_index"]
		config_spec = append(config_spec, a)

		requestpath_in["config_spec"] = config_spec

		flattenRequestList = append(flattenRequestList, requestpath_in)
	}
	log.Printf("datatatd %v", flattenRequestList)

	return flattenRequestList
}

func resourceApiGwPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceApiGwPolicyUpdate")
	c := m.(*service.NitroClient)

	resourceID := d.Id()
	endpoint := "policies"

	a := getApiGwPolicyPayload(d).(map[string]interface{})
	a["id"] = resourceID
	_, err := c.UpdateResource(endpoint, a, resourceID)

	returnData, err := c.GetResource("policiesdeploy", resourceID)
	if err != nil {
		return diag.Errorf("error Reading apigw deploy state %s", err)
	}
	status := returnData["configstatus"].(map[string]interface{})["status"].(string)
	// To check if the Policy is undeployed and user wants to deploy i.e., applyconfig action
	if status == "Indraft" && d.Get("deploy").(bool) == true {
		_, err := c.AddResourceWithActionParams(endpoint, resourceID, "applyconfig")
		if err != nil {
			return diag.FromErr(err)
		}
		err1 := c.WaitForDeplymentCompletion("policiesdeploy", resourceID, time.Duration(c.ActivityTimeout)*time.Second, "applyconfig")
		if err1 != nil {
			return diag.FromErr(err1)
		}
		// If the Policy is deployed and user wants to Undeploy it i.e, undeployconfig action
	} else if status == "Applied" && d.Get("deploy").(bool) == false {
		_, err := c.AddResourceWithActionParams(endpoint, resourceID, "undeployconfig")
		if err != nil {
			return diag.FromErr(err)
		}
		err1 := c.WaitForDeplymentCompletion("policiesdeploy", resourceID, time.Duration(c.ActivityTimeout)*time.Second, "undeployconfig")
		if err1 != nil {
			return diag.FromErr(err1)
		}
	} else {
		// nothing
	}

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceApiGwPolicyRead(ctx, d, m)

}

func resourceApiGwPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceApiGwPolicyDelete")
	var diags diag.Diagnostics

	c := m.(*service.NitroClient)

	endpoint := "policies"
	resourceID := d.Id()
	_, err := c.DeleteResource(endpoint, resourceID)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
