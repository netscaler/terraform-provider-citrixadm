package citrixadm

import (
	"fmt"
	"log"
	"terraform-provider-citrixadm/service"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const testAccStylebookAdd = `
resource "citrixadm_stylebook" "stylebook1" {
	name      = "basic-lb-config"
	namespace = "com.example.stylebooks"
	version   = "0.1"
	source    = "---\nname: basic-lb-config\nnamespace: com.example.stylebooks\nversion: \"0.1\"\ndisplay-name: Load Balancing Configuration\ndescription: This StyleBook defines a simple load balancing configuration.\nschema-version: \"1.0\"\nimport-stylebooks:\n  - namespace: netscaler.nitro.config\n    version: \"10.5\"\n    prefix: ns\nparameters:\n  - name: name\n    type: string\n    label: Application Name\n    description: Give a name to the application configuration.\n    required: true\n  - name: ip\n    type: ipaddress\n    label: Application Virtual IP (VIP)\n    description: The Application VIP that clients access\n    required: true\n  - name: lb-alg\n    type: string\n    label: LoadBalancing Algorithm\n    description: Choose the loadbalancing algorithm (method) used for loadbalancing client requests between the application servers.\n    allowed-values:\n      - ROUNDROBIN\n      - LEASTCONNECTION\n    default: ROUNDROBIN\n  - name: svc-servers\n    type: ipaddress[]\n    label: Application Server IPs\n    description: The IP addresses of all the servers of this application\n    required: true\n  - name: svc-port\n    type: tcp-port\n    label: Server Port\n    description: The TCP port open on the Application Servers to receive requests.\n    default: 80\ncomponents:\n  - name: lbvserver-comp\n    type: ns::lbvserver\n    properties:\n      name: $parameters.name + \"-lb\"\n      servicetype: HTTP\n      ipv46: $parameters.ip\n      port: 80\n      lbmethod: $parameters.lb-alg\n    components:\n      - name: svcg-comp\n        type: ns::servicegroup\n        properties:\n          servicegroupname: $parameters.name + \"-svcgrp\"\n          servicetype: HTTP\n        components:\n          - name: lbvserver-svg-binding-comp\n            type: ns::lbvserver_servicegroup_binding\n            properties:\n              name: $parent.parent.properties.name\n              servicegroupname: $parent.properties.servicegroupname\n          - name: members-svcg-comp\n            type: ns::servicegroup_servicegroupmember_binding\n            repeat: $parameters.svc-servers\n            repeat-item: srv\n            properties:\n              ip: $srv\n              port: $parameters.svc-port\n              servicegroupname: $parent.properties.servicegroupname\noutputs:\n  - name: lbvserver-comp\n    value: $components.lbvserver-comp\n    description: The component that builds the Nitro lbvserver configuration object\n  - name: servicegroup-comp\n    value: $components.lbvserver-comp.components.svcg-comp\n    description: The component that builds the Nitro servicegroup configuration object\n"
}
`

// example.Widget represents a concrete Go type that represents an API resource
func TestAccStylebook_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck:   func() { testAccPreCheck(t) },
		// ProviderFactories: providerFactories,
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckStylebookDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStylebookAdd,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStylebookExists("citrixadm_stylebook.stylebook1", nil),
				),
			},
		},
	})
}

func testAccCheckStylebookExists(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// retrieve the resource by name from state
		log.Println("[DEBUG] testAccCheckStylebookExists")
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Stylebook ID is set")
		}

		if id != nil {
			if *id != "" && *id != rs.Primary.ID {
				return fmt.Errorf("Resource ID has changed")
			}
			*id = rs.Primary.ID
		}

		// retrieve the client from the test provider
		c := testAccProvider.Meta().(*service.NitroClient)
		data, err := c.GetResource("stylebooks", rs.Primary.ID)

		if err != nil {
			return err
		}

		if data == nil {
			return fmt.Errorf("Resource %s not found", n)
		}

		return nil
	}
}

func testAccCheckStylebookDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*service.NitroClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "citrixadm_stylebook" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No name is set")
		}

		_, err := c.GetResource("stylebooks", rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Stylebook %s still exists", rs.Primary.ID)
		}

	}

	return nil
}
