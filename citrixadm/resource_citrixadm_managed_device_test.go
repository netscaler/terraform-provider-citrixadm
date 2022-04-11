package citrixadm

import (
	"fmt"
	"log"
	"terraform-provider-citrixadm/service"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Assumption before running this test:
// 1. There is a ADM agent with IP `10.0.1.91` already registered with CitrixADM, OR change the IP address in the test
// 2. the ns device profile is already present with name `nsroot_verysecretpassword_profile`
const testAccManagedDeviceAdd = `

data "citrixadm_mps_agent" "agent1" {
	name = "10.0.1.91"
  }

resource "citrixadm_managed_device" "device1" {
	ip_address    = "10.0.1.165"
	profile_name  = "nsroot_verysecretpassword_profile"
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
		CheckDestroy: testAccCheckManagedDeviceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccManagedDeviceAdd,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckManagedDeviceExists("citrixadm_managed_device.device1", nil),
				),
			},
		},
	})
}

func testAccCheckManagedDeviceExists(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// retrieve the resource by name from state
		log.Println("[DEBUG] testAccCheckManagedDeviceExists")
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Managed Device ID is set")
		}

		if id != nil {
			if *id != "" && *id != rs.Primary.ID {
				return fmt.Errorf("Resource ID has changed")
			}

			*id = rs.Primary.ID
		}

		// retrieve the client from the test provider
		c := testAccProvider.Meta().(*service.NitroClient)
		data, err := c.GetResource("managed_device", rs.Primary.ID)

		if err != nil {
			return err
		}

		if data == nil {
			return fmt.Errorf("Resource %s not found", n)
		}

		return nil
	}
}

func testAccCheckManagedDeviceDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*service.NitroClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "citrixadm_managed_device" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No name is set")
		}

		_, err := c.GetResource("managed_device", rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Managed Device %s still exists", rs.Primary.ID)
		}

	}

	return nil
}
