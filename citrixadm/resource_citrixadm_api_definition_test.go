package citrixadm

import (
	"fmt"
	"log"
	"terraform-provider-citrixadm/service"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const testAccApiDefinitionAdd = `
	resource "citrixadm_api_definition" "tf_def1" {
		name     = "tf-def"
		version  = "V2"
		title    = "my_tf_api"
		host     = "example.com"
		basepath = "/user"
		schemes  = []
		apiresources {
			paths   = "/user"
			methods = ["GET", "PUT"]
		}
		apiresources {
			paths   = "/user/action"
			methods = ["POST"]
		}
	}  
`

func TestAccApiDefinition_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck:   func() { testAccPreCheck(t) },
		// ProviderFactories: providerFactories,
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckApiDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccApiDefinitionAdd,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApiDefinitionExists("citrixadm_api_definition.tf_def1", nil),
				),
			},
		},
	})
}

func testAccCheckApiDefinitionExists(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// retrieve the resource by name from state
		log.Println("[DEBUG] testAccCheckApiDefinitionExists")
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Api Definition ID is set")
		}

		if id != nil {
			if *id != "" && *id != rs.Primary.ID {
				return fmt.Errorf("Resource ID has changed")
			}
			*id = rs.Primary.ID
		}

		// retrieve the client from the test provider
		c := testAccProvider.Meta().(*service.NitroClient)
		data, err := c.GetResource("apidefs", rs.Primary.ID)

		if err != nil {
			return err
		}

		if data == nil {
			return fmt.Errorf("Resource %s not found", n)
		}

		return nil
	}
}

func testAccCheckApiDefinitionDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*service.NitroClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "citrixadm_api_definition" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No name is set")
		}

		_, err := c.GetResource("apidefs", rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Api Definition %s still exists", rs.Primary.ID)
		}

	}

	return nil
}
