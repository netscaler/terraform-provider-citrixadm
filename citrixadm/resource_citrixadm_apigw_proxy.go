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

func resourceApiGwProxy() *schema.Resource {
	return &schema.Resource{
		Description:   "Configure API Proxy endpoint on API Gateway",
		CreateContext: resourceApiGwProxyCreate,
		ReadContext:   resourceApiGwProxyRead,
		UpdateContext: resourceApiGwProxyUpdate,
		DeleteContext: resourceApiGwProxyDelete,
		Schema: map[string]*schema.Schema{
			"proxy_name": {
				Description: "API Proxy name",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"host": {
				Description: "Hostname or IP address (Citrix ADC VIP) where API traffic will be received by Citrix ADC from API clients.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"port": {
				Description: "API Proxy Port number",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"protocol": {
				Description: "Protocol: http/https",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"service_fqdns": {
				Description: "Service FQDNs where API traffic will be received from api clients",
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tls_security_profile": {
				Description: "TLS Security profile",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"deploy": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"tls_certkey_objref": {
				Description: "Reference to existing TLS certkey object in ADM",
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"adm_certkey": {
							Description: "Reference to existing TLS certkey object in ADM",
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Description: "Cert File Id",
										Type:        schema.TypeString,
										Optional:    true,
									},
								},
							},
						},
					},
				},
			},
			"target_apigw": {
				Description: "Target ADC instance acting as API GW",
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_device": {
							Description: "ADC device details",
							Type:        schema.TypeList,
							Required:    true,
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
					},
				},
			},
		},
	}
}

func getApiGwProxyPayload(d *schema.ResourceData) interface{} {
	data := make(map[string]interface{})

	apiproxy := make(map[string]interface{})
	apiproxy["proxy_name"] = d.Get("proxy_name").(string)
	if v, ok := d.GetOk("host"); ok {
		apiproxy["host"] = v.(string)
	}
	if v, ok := d.GetOk("port"); ok {
		apiproxy["port"] = v.(int)
	}
	if v, ok := d.GetOk("protocol"); ok {
		apiproxy["protocol"] = v.(string)
	}
	if v, ok := d.GetOk("service_fqdns"); ok {
		apiproxy["service_fqdns"] = v.([]interface{})
	}
	if v, ok := d.GetOk("tls_security_profile"); ok {
		apiproxy["tls_security_profile"] = v.(string)
	}

	//To Take a map and then convert it into string(required Payload)-and it is then added to the payload
	tls_certkey_objref := make(map[string]interface{})
	z := d.Get("tls_certkey_objref").([]interface{})[0].(map[string]interface{})
	x := z["adm_certkey"].([]interface{})
	tls_certkey_objref["adm_certkey"] = mapToString(x[0].(map[string]interface{}))
	apiproxy["tls_certkey_objref"] = tls_certkey_objref

	data["apiproxy"] = apiproxy

	target_apigw := make(map[string]interface{})
	a := d.Get("target_apigw").([]interface{})[0].(map[string]interface{})
	e := a["target_device"].([]interface{})
	target_apigw["target_device"] = mapToString(e[0].(map[string]interface{}))

	data["target_apigw"] = target_apigw

	return data

}

//To convert the map[string]interface{} to string
func mapToString(a map[string]interface{}) string {
	b, _ := json.Marshal(a)
	c := string(b)
	return c
}

func resourceApiGwProxyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceApiGwProxyCreate")

	c := m.(*service.NitroClient)

	endpoint := "apiproxies"

	returnData, err := c.AddResource(endpoint, getApiGwProxyPayload(d))

	if err != nil {
		return diag.Errorf("error creating apigw_proxy: %s", err.Error())
	}

	resourceID := returnData["apiproxy"].(map[string]interface{})["id"].(string)
	log.Printf("id %s", resourceID)

	d.SetId(resourceID)

	if d.Get("deploy").(bool) == true {
		_, err := c.AddResourceWithActionParams(endpoint, resourceID, "applyconfig")
		if err != nil {
			return diag.FromErr(err)
		}
		err1 := c.WaitForDeplymentCompletion("apiproxiesdeploy", resourceID, time.Duration(c.ActivityTimeout)*time.Second, "applyconfig")
		if err1 != nil {
			return diag.FromErr(err1)
		}

	}

	return resourceApiGwProxyRead(ctx, d, m)
}

func resourceApiGwProxyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceApiGwProxyRead")
	var diags diag.Diagnostics
	c := m.(*service.NitroClient)

	resourceID := d.Id()
	endpoint := "apiproxies"

	returnData, err := c.GetResource(endpoint, resourceID)
	if err != nil {
		return diag.Errorf("error reading apigw_proxy %s: %s", resourceID, err.Error())
	}

	getResponseData := returnData["apiproxy"].(map[string]interface{})

	returnData1, err1 := c.GetResource("apiproxiesdeploy", resourceID)
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

	d.Set("proxy_name", getResponseData["proxy_name"].(string))
	d.Set("host", getResponseData["host"].(string))
	d.Set("port", getResponseData["port"].(float64))
	d.Set("protocol", getResponseData["protocol"].(string))
	d.Set("service_fqdns", getResponseData["service_fqdns"].([]interface{}))
	d.Set("tls_security_profile", getResponseData["tls_security_profile"].(string))
	if _, ok := getResponseData["tls_certkey_objref"]; ok {
		d.Set("tls_certkey_objref", flattenApiGwProxyAttribute(getResponseData["tls_certkey_objref"].(map[string]interface{})))
	}
	d.Set("target_apigw", flattenApiGwProxyAttribute(returnData["target_apigw"].(map[string]interface{})))

	return diags
}

func flattenApiGwProxyAttribute(attribute map[string]interface{}) []interface{} {
	s := make(map[string]interface{})
	a := make(map[string]interface{})
	var b string
	var c string
	if val, ok := attribute["target_device"]; ok {
		b = val.(string)
		c = "target_device"
	} else if val, ok := attribute["adm_certkey"]; ok {
		b = val.(string)
		c = "adm_certkey"
	}
	//To covert jsonstring(b) to map(a)
	json.Unmarshal([]byte(b), &a)

	//To list/Slice
	s[c] = []interface{}{a}
	//To list/Slice
	return []interface{}{s}
}

func resourceApiGwProxyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceApiGwProxyUpdate")
	c := m.(*service.NitroClient)

	resourceID := d.Id()
	endpoint := "apiproxies"

	_, err := c.UpdateResource(endpoint, getApiGwProxyPayload(d), resourceID)

	//To Get Deploy Status of resource
	returnData, err := c.GetResource("apiproxiesdeploy", resourceID)
	if err != nil {
		return diag.Errorf("error Reading apigw deploy state %s", err)
	}
	status := returnData["configstatus"].(map[string]interface{})["status"].(string)

	// To check if the Proxy is undeployed and user wants to deploy i.e., applyconfig action
	if status == "Indraft" && d.Get("deploy").(bool) == true {
		_, err := c.AddResourceWithActionParams(endpoint, resourceID, "applyconfig")
		if err != nil {
			return diag.FromErr(err)
		}
		err1 := c.WaitForDeplymentCompletion("apiproxiesdeploy", resourceID, time.Duration(c.ActivityTimeout)*time.Second, "applyconfig")
		if err1 != nil {
			return diag.FromErr(err1)
		}
		// If the Proxy is deployed and user wants to Undeploy it i.e, undeployconfig action
	} else if status == "Applied" && d.Get("deploy").(bool) == false {
		_, err := c.AddResourceWithActionParams(endpoint, resourceID, "undeployconfig")
		if err != nil {
			return diag.FromErr(err)
		}
		err1 := c.WaitForDeplymentCompletion("apiproxiesdeploy", resourceID, time.Duration(c.ActivityTimeout)*time.Second, "undeployconfig")
		if err1 != nil {
			return diag.FromErr(err1)
		}
	} else {
		// nothing
	}

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceApiGwProxyRead(ctx, d, m)

}

func resourceApiGwProxyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceApiGwProxyDelete")
	var diags diag.Diagnostics

	c := m.(*service.NitroClient)

	endpoint := "apiproxies"
	resourceID := d.Id()
	_, err := c.DeleteResource(endpoint, resourceID)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
