package citrixadm

import (
	"context"
	"log"

	"terraform-provider-citrixadm/service"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApiGwDefinition() *schema.Resource {
	return &schema.Resource{
		Description:   "Create and Manage API Definition",
		CreateContext: resourceApiGwDefinitionCreate,
		ReadContext:   resourceApiGwDefinitionRead,
		UpdateContext: resourceApiGwDefinitionUpdate,
		DeleteContext: resourceApiGwDefinitionDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "API Definition name",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"title": {
				Description: "Title for API Definition",
				Type:        schema.TypeString,
				Required:    true,
			},
			"version": {
				Description: "API Definition version",
				Type:        schema.TypeString,
				Required:    true,
			},
			"host": {
				Description: "Host FQDN where API service is hosted",
				Type:        schema.TypeString,
				Required:    true,
			},
			"basepath": {
				Description: "API Definition base path - this is appended as prefix to all API resources",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"schemes": {
				Description: "Schemes of API Definition , HTTP/HTTPS",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
			},
			"apiresources": {
				Description: "List of HTTP Methods and API Resource paths.",
				Type:        schema.TypeList,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"paths": {
							Description: "API Path",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"methods": {
							Description: "API Method",
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func getApiGwDefinitionPayload(d *schema.ResourceData) interface{} {
	data := make(map[string]interface{})

	data["name"] = d.Get("name").(string)
	data["title"] = d.Get("title").(string)
	data["version"] = d.Get("version").(string)
	data["host"] = d.Get("host").(string)
	data["schemes"] = d.Get("schemes").([]interface{})

	if v, ok := d.GetOk("basepath"); ok {
		data["basepath"] = v.(string)
	}
	data["apiresources"] = d.Get("apiresources").([]interface{})

	// var payload []interface{}
	// payload = append(payload, data)

	return data

}

func resourceApiGwDefinitionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceApiGwDefinitionCreate")

	c := m.(*service.NitroClient)

	endpoint := "apidefs"

	returnData, err := c.AddResource(endpoint, getApiGwDefinitionPayload(d))

	if err != nil {
		return diag.Errorf("error creating apidef: %s", err.Error())
	}

	resourceID := returnData["apidef"].(map[string]interface{})["id"].(string)
	log.Printf("id %s", resourceID)

	d.SetId(resourceID)

	return resourceApiGwDefinitionRead(ctx, d, m)
}

func resourceApiGwDefinitionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceApiGwDefinitionRead")
	var diags diag.Diagnostics
	c := m.(*service.NitroClient)

	resourceID := d.Id()
	endpoint := "apidefs"

	returnData, err := c.GetResource(endpoint, resourceID)
	if err != nil {
		return diag.Errorf("error reading apidef %s: %s", resourceID, err.Error())
	}

	getResponseData := returnData["apidef"].(map[string]interface{})

	log.Println("getResponseData", getResponseData)

	d.Set("name", getResponseData["name"].(string))
	d.Set("title", getResponseData["title"].(string))
	d.Set("version", getResponseData["version"].(string))
	d.Set("host", getResponseData["host"].(string))
	d.Set("basepath", getResponseData["basepath"].(string))
	d.Set("schemes", getResponseData["schemes"].([]interface{}))
	d.Set("apiresources", getResponseData["apiresources"].([]interface{}))

	return diags
}

func resourceApiGwDefinitionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceApiGwDefinitionUpdate")
	c := m.(*service.NitroClient)

	resourceID := d.Id()
	endpoint := "apidefs"

	payload := getApiGwDefinitionPayload(d)
	payload_id := payload.(map[string]interface{})
	payload_id["id"] = resourceID
	_, err := c.UpdateResource(endpoint, payload_id, resourceID)

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceApiGwDefinitionRead(ctx, d, m)

}

func resourceApiGwDefinitionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceApiGwDefinitionDelete")
	var diags diag.Diagnostics

	c := m.(*service.NitroClient)

	endpoint := "apidefs"
	resourceID := d.Id()
	_, err := c.DeleteResource(endpoint, resourceID)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
