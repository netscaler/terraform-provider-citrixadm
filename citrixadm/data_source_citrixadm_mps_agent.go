package citrixadm

import (
	"context"
	"fmt"
	"log"

	"terraform-provider-citrixadm/service"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMpsAgent() *schema.Resource {
	return &schema.Resource{
		Description: "Get a mps agent ID and Datacenter ID by MPS Agent IP address",
		ReadContext: dataSourceMpsAgentRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Agent IP Address",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"datacenter_id": {
				Description: "Datacenter Id is system generated key for data center",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func dataSourceMpsAgentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In dataSourceMpsAgentRead")
	var diags diag.Diagnostics
	c := m.(*service.NitroClient)

	endpoint := "mps_agent"

	agentIP := d.Get("name").(string)

	returnData, err := c.GetAllResource(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}

	// Find the correct resource with the given name and store agent_id and datacenter_id from the object
	for _, v := range returnData[endpoint].([]interface{}) {
		agent := v.(map[string]interface{})
		if agent["name"].(string) == agentIP {
			d.SetId(agent["id"].(string))
			d.Set("name", agent["name"].(string))
			d.Set("datacenter_id", agent["datacenter_id"].(string))
			break
		}
	}
	if d.Id() == "" {
		return diag.FromErr(fmt.Errorf("Mps Agent with name %s not found", agentIP))
	}
	return diags
}
