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
// 1. There is a managed device with IP 10.0.1.42 already registered with CitrixADM, OR change the IP address in the test
// 2. There is no lb vserver with name "tf-acc-lb" with VIP IP "4.4.4.4"

const testAccStylebookConfigpackAdd = `

data "citrixadm_managed_device" "device1" {
  ip_address = "10.0.1.42"
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
    instance_id = data.citrixadm_managed_device.device1.id
  }
}
`

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
