package citrixadm

import (
	"fmt"
	"log"
	"terraform-provider-citrixadm/service"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const testAccApiGwDeploymentAdd = `
	data "citrixadm_apigw_definition" "demo" {
		name = "new_definition"
	}
	data "citrixadm_managed_device" "device1" {
		ip_address = "10.10.10.10"
	}
	data "citrixadm_apigw_proxy" "demo_proxy" {
		name = "testing"
	}
  	resource "citrixadm_apigw_deployment" "demo" {
		api_id = data.citrixadm_apigw_definition.demo.id
		apiproxy_ref {
		  	id = data.citrixadm_apigw_proxy.demo_proxy.id
		}
		name   = "tf_deploy111"
		tags   = []
		deploy = true
		target_apigw {
			id           = data.citrixadm_managed_device.device1.id
			display_name = data.citrixadm_managed_device.device1.ip_address
		}
		routes {
			name                 = "routing_name_0"
			route_param          = "/user/action"
			route_paramtype      = "resource_path"
			upstreamservice_name = "sdf"
		}
		routes {
			name                 = "route_name2_0"
			route_param          = "/user"
			route_paramtype      = "resource_path"
			upstreamservice_name = "sdf"
		}
		upstreamservices {        
			scheme            = "http"
			service_fqdn      = "www.example.com"
			service_fqdn_port = 80
	  		backend_servers {
				ip4_addr = "2.2.2.2"
				port     = 80
	  		}
			name              = "sdf"
		}
		upstreamservices {
	  		backend_servers {
				ip4_addr = "2.2.2.2"
				port     = 80
	  		}
			name              = "sdf1"
			scheme            = "http"
			service_fqdn      = "www.example.com"
			service_fqdn_port = 80
		}
  	}
  
`

func TestAccApiGwDeployment_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck:   func() { testAccPreCheck(t) },
		// ProviderFactories: providerFactories,
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckApiGwDeploymentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccApiGwDeploymentAdd,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApiGwDeploymentExists("citrixadm_apigw_deployment.demo", nil),
				),
			},
		},
	})
}

func testAccCheckApiGwDeploymentExists(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// retrieve the resource by name from state
		log.Println("[DEBUG] testAccCheckApiGwDeploymentExists")
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ApiGw Deployment ID is set")
		}

		if id != nil {
			if *id != "" && *id != rs.Primary.ID {
				return fmt.Errorf("Resource ID has changed")
			}
			*id = rs.Primary.ID
		}

		// retrieve the client from the test provider
		c := testAccProvider.Meta().(*service.NitroClient)
		data, err := c.GetResource("deployments", rs.Primary.ID)

		if err != nil {
			return err
		}

		if data == nil {
			return fmt.Errorf("Resource %s not found", n)
		}

		return nil
	}
}

func testAccCheckApiGwDeploymentDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*service.NitroClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "citrixadm_apigw_deployment" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No name is set")
		}

		_, err := c.GetResource("deployments", rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("ApiGw Deployment %s still exists", rs.Primary.ID)
		}

	}

	return nil
}
