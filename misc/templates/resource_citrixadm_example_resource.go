package citrixadm

import (
	"context"
	"log"

	// "terraform-provider-citrixadm/service"
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

func resourceExampleResourceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	log.Printf("In resourceExampleResourceCreate")

	return resourceExampleResourceRead(ctx, d, m)
}

func resourceExampleResourceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	log.Printf("In resourceExampleResourceRead")
	var diags diag.Diagnostics

	return diags
}

func resourceExampleResourceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	log.Printf("In resourceExampleResourceUpdate")

	return resourceExampleResourceRead(ctx, d, m)
}

func resourceExampleResourceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("In resourceExampleResourceDelete")
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}
