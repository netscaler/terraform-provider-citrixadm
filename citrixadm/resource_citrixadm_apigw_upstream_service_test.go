package citrixadm

import (
	"fmt"
	"log"
	"terraform-provider-citrixadm/service"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const testAccApigwUpstreamServiceAdd = `
	data "citrixadm_apigw_deployment" "demo" {
		name = "deploy"
	}
	resource "citrixadm_apigw_upstream_service" "service1" {
		name              = "tf_upstream_service"
		scheme            = "http"
		service_fqdn      = "www.example.com"
		service_fqdn_port = 80
		deployment_id     = data.citrixadm_apigw_deployment.demo.id
		backend_servers {
		ip4_addr = "2.1.2.4"
		port     = 80
		}
	}  
`

func TestAccApigwUpstreamService_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck:   func() { testAccPreCheck(t) },
		// ProviderFactories: providerFactories,
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckApigwUpstreamServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccApigwUpstreamServiceAdd,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigwUpstreamServiceExists("citrixadm_apigw_upstream_service.service1", nil),
					resource.TestCheckResourceAttr("citrixadm_apigw_upstream_service.service1", "name", "tf_upstream_service"),
					resource.TestCheckResourceAttr("citrixadm_apigw_upstream_service.service1", "scheme", "http"),
					resource.TestCheckResourceAttr("citrixadm_apigw_upstream_service.service1", "service_fqdn", "www.example.com"),
				),
			},
		},
	})
}

func testAccCheckApigwUpstreamServiceExists(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// retrieve the resource by name from state
		log.Println("[DEBUG] testAccCheckApigwUpstreamServiceExists")
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ApigwUpstreamService ID is set")
		}

		if id != nil {
			if *id != "" && *id != rs.Primary.ID {
				return fmt.Errorf("Resource ID has changed")
			}
			*id = rs.Primary.ID
		}

		// retrieve the client from the test provider
		c := testAccProvider.Meta().(*service.NitroClient)
		data, err := c.GetResource("upstreamservices", rs.Primary.ID)

		if err != nil {
			return err
		}

		if data == nil {
			return fmt.Errorf("Resource %s not found", n)
		}

		return nil
	}
}

func testAccCheckApigwUpstreamServiceDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*service.NitroClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "citrixadm_apigw_upstream_service" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No name is set")
		}

		_, err := c.GetResource("upstreamservices", rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("ApiGw Upstream Service %s still exists", rs.Primary.ID)
		}

	}

	return nil
}
