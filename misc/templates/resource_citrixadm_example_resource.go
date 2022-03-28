package citrixadm

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"terraform-provider-citrixadm/service"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceExampleResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceExampleResourceCreate,
		ReadContext:   resourceExampleResourceRead,
		UpdateContext: resourceExampleResourceUpdate,
		DeleteContext: resourceExampleResourceDelete,
		Schema: map[string]*schema.Schema{

			"sample_attribute": {
				// This description is used by the documentation generator and the language server.
				Description: "Sample attribute.",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
		},
	}
}

func getExampleResourcePayload(d *schema.ResourceData) []interface{} {
	data := make(map[string]interface{})

	if v, ok := d.GetOk("name"); ok {
		data["name"] = v.(string)
	}

	var payload []interface{}
	payload = append(payload, data)

	return payload

}
func resourceExampleResourceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceExampleResourceCreate")

	c := m.(*service.NitroClient)

	endpoint := "example_resource"

	n := service.NitroRequestParams{
		Resource: endpoint,

		ResourcePath:       fmt.Sprintf("massvc/%s/nitro/v2/config/%s", c.CustomerID, endpoint),
		ResourceData:       getExampleResourcePayload(d),
		Method:             "POST",
		SuccessStatusCodes: []int{200, 201},
	}

	body, err := c.MakeNitroRequest(n)
	if err != nil {
		return diag.FromErr(err)
	}
	var returnData map[string]interface{}

	err = json.Unmarshal(body, &returnData)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("return data %v", returnData)

	resourceID := returnData[endpoint].([]interface{})[0].(map[string]interface{})["id"].(string)
	log.Printf("id %s", resourceID)

	d.SetId(resourceID)
	return resourceExampleResourceRead(ctx, d, m)
}

func resourceExampleResourceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceExampleResourceRead")
	var diags diag.Diagnostics
	c := m.(*service.NitroClient)

	resourceID := d.Id()
	endpoint := "example_resource"

	n := service.NitroRequestParams{
		ResourcePath:       fmt.Sprintf("massvc/%s/nitro/v2/config/%s/%s", c.CustomerID, endpoint, resourceID),
		Method:             "GET",
		Resource:           endpoint,
		ResourceData:       d,
		SuccessStatusCodes: []int{200},
	}

	body, err := c.MakeNitroRequest(n)
	if err != nil {
		return diag.FromErr(err)
	}
	var returnData map[string]interface{}

	err = json.Unmarshal(body, &returnData)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("return data %v", returnData)
	getResponseData := returnData[endpoint].([]interface{})[0].(map[string]interface{})

	log.Println("getResponseData", getResponseData)

	d.Set("name", getResponseData["name"].(string))

	return diags
}

func resourceExampleResourceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceExampleResourceUpdate")
	c := m.(*service.NitroClient)

	resourceID := d.Id()
	endpoint := "example_resource"

	n := service.NitroRequestParams{
		Resource:           endpoint,
		ResourcePath:       fmt.Sprintf("massvc/%s/nitro/v2/config/%s/%s", c.CustomerID, endpoint, resourceID),
		ResourceData:       getExampleResourcePayload(d),
		Method:             "PUT",
		SuccessStatusCodes: []int{200, 201},
	}

	body, err := c.MakeNitroRequest(n)
	if err != nil {
		return diag.FromErr(err)
	}
	var returnData map[string]interface{}

	err = json.Unmarshal(body, &returnData)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("return data %v", returnData)
	return resourceExampleResourceRead(ctx, d, m)
}

func resourceExampleResourceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceExampleResourceDelete")
	var diags diag.Diagnostics

	c := m.(*service.NitroClient)

	endpoint := "example_resource"
	resourceID := d.Id()

	n := service.NitroRequestParams{
		ResourcePath:       fmt.Sprintf("massvc/%s/nitro/v2/config/%s/%s", c.CustomerID, endpoint, resourceID),
		Method:             "DELETE",
		Resource:           endpoint,
		SuccessStatusCodes: []int{200, 204},
	}

	body, err := c.MakeNitroRequest(n)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("delete response %v", body)

	d.SetId("")

	return diags
}
