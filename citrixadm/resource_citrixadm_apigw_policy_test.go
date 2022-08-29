package citrixadm

import (
	"fmt"
	"log"
	"terraform-provider-citrixadm/service"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const testAccApiGwPolicyAdd = `

	data "citrixadm_apigw_deployment" "demo" {
		name = "tf_deploy"
	}
	data "citrixadm_apigw_upstream_service" "upstreamdemo" {
		name 		  = "sdf"
		deployment_id = data.citrixadm_apigw_deployment.demo.id
	}
	resource "citrixadm_apigw_policy" "policy" {

		policygroup_name = "tf_policy"
		requestpath {
			order_index = 1
			policy_name = "asaaa"
			policytype  = "NoAuth"
			config_spec {
				api_resource_paths {
					delete = false
					endpoints = [
					"/user"
					]
					get   = true
					patch = false
					post  = true
					put   = false
				}
				custom_rules {
					delete = false
					endpoints = [
						"/sss"
					]
					get   = false
					patch = false
					post  = false
					put   = false
				}
			}
		}
		upstreamservice_id = data.citrixadm_apigw_upstream_service.upstreamdemo.id
  		deployment_id      = data.citrixadm_apigw_deployment.demo.id
		deploy             = true
	}  
`

func TestAccApiGwPolicy_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck:   func() { testAccPreCheck(t) },
		// ProviderFactories: providerFactories,
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckApiGwPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccApiGwPolicyAdd,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApiGwPolicyExists("citrixadm_apigw_policy.policy", nil),
				),
			},
		},
	})
}

func testAccCheckApiGwPolicyExists(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// retrieve the resource by name from state
		log.Println("[DEBUG] testAccCheckApiGwPolicyExists")
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ApiGw Policy ID is set")
		}

		if id != nil {
			if *id != "" && *id != rs.Primary.ID {
				return fmt.Errorf("Resource ID has changed")
			}
			*id = rs.Primary.ID
		}

		// retrieve the client from the test provider
		c := testAccProvider.Meta().(*service.NitroClient)
		data, err := c.GetResource("policies", rs.Primary.ID)

		if err != nil {
			return err
		}

		if data == nil {
			return fmt.Errorf("Resource %s not found", n)
		}

		return nil
	}
}

func testAccCheckApiGwPolicyDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*service.NitroClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "citrixadm_apigw_policy" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No name is set")
		}

		_, err := c.GetResource("policies", rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("ApiGw Policy %s still exists", rs.Primary.ID)
		}

	}

	return nil
}
