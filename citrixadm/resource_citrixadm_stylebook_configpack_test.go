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
	testAccStylebookConfigpackPlaceholder = `
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

		resource "citrixadm_stylebook_configpack" "cfgpack1" {
			stylebook {
				name      = "lb"
				namespace = "com.citrix.adc.stylebooks"
				version   = "1.1"
			}
			parameters = {
				lb-appname       = "tf-acc-lb"
				lb-service-type  = "HTTP"
				lb-virtual-ip    = "4.4.4.4"
				lb-virtual-port  = "80"
				svc-service-type = "HTTP"
			}
			targets {
				instance_id = citrixadm_managed_device.device1.id
			}
		}
	`
)

var testAccStylebookConfigpackAdd = fmt.Sprintf(testAccStylebookConfigpackPlaceholder,
	os.Getenv("AGENT_IP"),
	os.Getenv("VPX_USER"),
	os.Getenv("VPX_PASSWORD"),
	os.Getenv("VPX_IP"),
)

// example.Widget represents a concrete Go type that represents an API resource
func TestAccStylebookConfigpack_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck:   func() { testAccPreCheck(t) },
		// ProviderFactories: providerFactories,
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckStylebookConfigpackDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStylebookConfigpackAdd,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStylebookConfigpackExists("citrixadm_stylebook_configpack.cfgpack1", nil),
				),
			},
		},
	})
}

func testAccCheckStylebookConfigpackExists(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// retrieve the resource by name from state
		log.Println("[DEBUG] testAccCheckStylebookConfigpackExists")
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Stylebook Configpack ID is set")
		}

		if id != nil {
			if *id != "" && *id != rs.Primary.ID {
				return fmt.Errorf("Resource ID has changed")
			}
			*id = rs.Primary.ID
		}

		// retrieve the client from the test provider
		c := testAccProvider.Meta().(*service.NitroClient)
		data, err := c.GetResource("configpacks", rs.Primary.ID)

		if err != nil {
			return err
		}

		if data == nil {
			return fmt.Errorf("Resource %s not found", n)
		}

		return nil
	}
}

func testAccCheckStylebookConfigpackDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*service.NitroClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "citrixadm_stylebook_configpack" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No name is set")
		}

		_, err := c.GetResource("configpacks", rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Stylebook Configpack %s still exists", rs.Primary.ID)
		}

	}

	return nil
}
