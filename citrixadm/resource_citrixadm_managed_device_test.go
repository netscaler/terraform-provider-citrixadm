package citrixadm

import (
	"github.com/hashicorp/terraform/helper/resource"
	"terraform-provider-citrixadm/service"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const testAccManagedDeviceAdd = `




data "citrixadm_mps_agent" "agent1" {
	name = "10.0.1.91"
  }

resource "citrixadm_managed_device" "device5" {
	ip_address    = "10.0.1.92"
	profile_name  = "nsroot_notnsroot_profile"
	datacenter_id = data.citrixadm_mps_agent.agent1.datacenter_id
	agent_id      = data.citrixadm_mps_agent.agent1.id
}
`

// example.Widget represents a concrete Go type that represents an API resource
func TestAccManagedDevice_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck:   func() { testAccPreCheck(t) },
		// ProviderFactories: providerFactories,
		Providers:    testAccProviders,
		CheckDestroy: nil, //testAccCheckManagedDeviceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccManagedDeviceAdd,
				Check:  resource.ComposeTestCheckFunc(
				// testAccCheckManagedDeviceExists("example_widget.foo", &widgetBefore),
				),
			},
		},
	})
}

// // testAccCheckManagedDeviceDestroy ...
// func testAccCheckManagedDeviceDestroy(s *terraform.State) error {
// 	// retrieve the connection established in Provider configuration
// 	conn := testAccProvider.Meta().(*service.NitroClient)

// 	// loop through the resources in state, verifying each widget
// 	// is destroyed
// 	for _, rs := range s.RootModule().Resources {
// 	  if rs.Type != "example_widget" {
// 		continue
// 	  }

// 	  // Retrieve our widget by referencing it's state ID for API lookup
// 	  request := &example.DescribeWidgets{
// 		IDs: []string{rs.Primary.ID},
// 	  }

// 	  response, err := conn.DescribeWidgets(request)
// 	  if err == nil {
// 		if len(response.Widgets) > 0 && *response.Widgets[0].ID == rs.Primary.ID {
// 		  return fmt.Errorf("Widget (%s) still exists.", rs.Primary.ID)
// 		}

// 		return nil
// 	  }

// 	  // If the error is equivalent to 404 not found, the widget is destroyed.
// 	  // Otherwise return the error
// 	  if !strings.Contains(err.Error(), "Widget not found") {
// 		return err
// 	  }
// 	}

// 	return nil
//   }
