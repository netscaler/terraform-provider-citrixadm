package citrixadm

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// providerFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.
// var providerFactories = map[string]func() (*schema.Provider, error){
// 	"citrixadm": func() (*schema.Provider, error) {
// 		return New("dev")(), nil
// 	},
// }

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.

	requiredEnvVariables := []string{
		"CITRIXADM_HOST",
		"CITRIXADM_CLIENT_ID",
		"CITRIXADM_CLIENT_SECRET",
		"CITRIXADM_HOST_LOCATION",
		"CITRIXADM_CUSTOMER_ID",
	}
	for _, envVar := range requiredEnvVariables {
		if v := os.Getenv(envVar); v == "" {
			t.Fatalf("%s must be set for acceptance tests", envVar)
		}
	}
}

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"citrixadm": testAccProvider,
	}
}
