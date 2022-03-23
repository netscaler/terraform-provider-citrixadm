package citrixadm

import (
	"fmt"
	"log"
	"terraform-provider-citrixadm/service"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const testAccNsDeviceProfileAdd = `
	resource "citrixadm_ns_device_profile" "profile1" {
		name     = "tf_acc_test_profile"
		username = "nsroot"
		password = "verysecretpassword"
	}
`

func TestAccNsDeviceProfile_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck:   func() { testAccPreCheck(t) },
		// ProviderFactories: providerFactories,
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNsDeviceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNsDeviceProfileAdd,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNsDeviceProfileExists("citrixadm_ns_device_profile.profile1", nil),
				),
			},
		},
	})
}

func testAccCheckNsDeviceProfileExists(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// retrieve the resource by name from state
		log.Println("[DEBUG] testAccCheckNsDeviceProfileExists")
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No NS Device Profile ID is set")
		}

		if id != nil {
			if *id != "" && *id != rs.Primary.ID {
				return fmt.Errorf("Resource ID has changed")
			}
			*id = rs.Primary.ID
		}

		// retrieve the client from the test provider
		c := testAccProvider.Meta().(*service.NitroClient)
		data, err := c.GetResource("ns_device_profile", rs.Primary.ID)

		if err != nil {
			return err
		}

		if data == nil {
			return fmt.Errorf("Resource %s not found", n)
		}

		return nil
	}
}

func testAccCheckNsDeviceProfileDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*service.NitroClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "citrixadm_ns_device_profile" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No name is set")
		}

		_, err := c.GetResource("ns_device_profile", rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("NS Device Profile %s still exists", rs.Primary.ID)
		}

	}

	return nil
}
