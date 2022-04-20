package citrixadm

import (
	"fmt"
	"log"
	"os"
	"terraform-provider-citrixadm/service"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Assumption before running this test:
// 1. There is a ADM agent with IP `AGENT_IP` already registered with CitrixADM
const (
	testAccManagedDevicePlaceholder = `
		data "citrixadm_mps_agent" "agent1" {
			name = "%s"
		}

		resource "citrixadm_ns_device_profile" "profile1" {
			name     = "tf_acc_test_profile"
			username = "%s"
			password = "%s"
		}

		resource "citrixadm_managed_device" "device1" {
			ip_address    = "%s"
			profile_name  = citrixadm_ns_device_profile.profile1.name
			datacenter_id = data.citrixadm_mps_agent.agent1.datacenter_id
			agent_id      = data.citrixadm_mps_agent.agent1.id
		}
	`
)

var testAccManagedDeviceAdd = fmt.Sprintf(testAccManagedDevicePlaceholder,
	os.Getenv("AGENT_IP"),
	os.Getenv("VPX_USER"),
	os.Getenv("VPX_PASSWORD"),
	os.Getenv("VPX_IP"),
)

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
