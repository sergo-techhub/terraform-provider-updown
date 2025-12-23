package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// providerFactories are used to instantiate the provider during acceptance testing.
var providerFactories = map[string]func() (*schema.Provider, error){
	"updown": func() (*schema.Provider, error) {
		return New()(), nil
	},
}

func TestProvider(t *testing.T) {
	if err := New()().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = New()()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("UPDOWN_API_KEY"); v == "" {
		t.Fatal("UPDOWN_API_KEY must be set for acceptance tests")
	}
}
