package citrixadm

import (
	"fmt"
	"log"
	"terraform-provider-citrixadm/service"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const testAccApiGwProxyAdd = `
	data "citrixadm_managed_device" "device1" {
		ip_address = "10.0.1.42"
	}
	resource "citrixadm_apigw_proxy" "proxy1" {
		proxy_name           = "tf_proxy101"
		host                 = "2.2.2.2"
		port                 = 22
		protocol             = "https"
		deploy               = "false"
		service_fqdns        = ["api.example.com"]
		tls_security_profile = "High Security"
		tls_certkey_objref {
		adm_certkey {
			id = "abcdefgh-1234-5678-ijkl-mnopqrstuvwx"
			}
		}
		target_apigw {
			target_device {
				id           = data.citrixadm_managed_device.device1.id
				display_name = data.citrixadm_managed_device.device1.ip_address
			}
		}
	}  
`

func TestAccApiGwProxy_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck:   func() { testAccPreCheck(t) },
		// ProviderFactories: providerFactories,
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckApiGwProxyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccApiGwProxyAdd,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApiGwProxyExists("citrixadm_apigw_proxy.proxy1", nil),
				),
			},
		},
	})
}

func testAccCheckApiGwProxyExists(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// retrieve the resource by name from state
		log.Println("[DEBUG] testAccCheckApiGwProxyExists")
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ApiGw Proxy ID is set")
		}

		if id != nil {
			if *id != "" && *id != rs.Primary.ID {
				return fmt.Errorf("Resource ID has changed")
			}
			*id = rs.Primary.ID
		}

		// retrieve the client from the test provider
		c := testAccProvider.Meta().(*service.NitroClient)
		data, err := c.GetResource("apiproxies", rs.Primary.ID)

		if err != nil {
			return err
		}

		if data == nil {
			return fmt.Errorf("Resource %s not found", n)
		}

		return nil
	}
}

func testAccCheckApiGwProxyDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*service.NitroClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "citrixadm_apigw_proxy" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No name is set")
		}

		_, err := c.GetResource("apiproxies", rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("ApiGw Proxy %s still exists", rs.Primary.ID)
		}

	}

	return nil
}
