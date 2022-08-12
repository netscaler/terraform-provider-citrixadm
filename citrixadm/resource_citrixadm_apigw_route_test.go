package citrixadm

import (
	"fmt"
	"log"
	"terraform-provider-citrixadm/service"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const testAccApiGwRouteAdd = `
	data "citrixadm_apigw_deployment" "demo" {
		name = "deploy"
	}
	resource "citrixadm_apigw_route" "route1" {
		name                 = "tf_route1"
		route_param          = "/user/action"
		route_paramtype      = "resource_path"
		upstreamservice_name = "upstream_service10"
		deployment_id        = data.citrixadm_apigw_deployment.demo.id
	}
`

func TestAccApiGwRoute_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck:   func() { testAccPreCheck(t) },
		// ProviderFactories: providerFactories,
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckApiGwRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccApiGwRouteAdd,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApiGwRouteExists("citrixadm_apigw_route.route1", nil),
				),
			},
		},
	})
}

func testAccCheckApiGwRouteExists(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// retrieve the resource by name from state
		log.Println("[DEBUG] testAccCheckApiGwRouteExists")
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No API route ID is set")
		}

		if id != nil {
			if *id != "" && *id != rs.Primary.ID {
				return fmt.Errorf("Resource ID has changed")
			}
			*id = rs.Primary.ID
		}

		// retrieve the client from the test provider
		c := testAccProvider.Meta().(*service.NitroClient)
		data, err := c.GetResource("routes", rs.Primary.ID)

		if err != nil {
			return err
		}

		if data == nil {
			return fmt.Errorf("Resource %s not found", n)
		}

		return nil
	}
}

func testAccCheckApiGwRouteDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*service.NitroClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "citrixadm_apigw_route" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No name is set")
		}

		_, err := c.GetResource("routes", rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("API Route %s still exists", rs.Primary.ID)
		}

	}

	return nil
}
