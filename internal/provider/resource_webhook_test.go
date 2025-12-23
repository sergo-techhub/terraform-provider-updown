package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccUpdownWebhook_basic(t *testing.T) {
	resourceName := "updown_webhook.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckUpdownWebhookDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccUpdownWebhookConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUpdownWebhookExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "url", "https://example.com/webhook/test"),
				),
			},
		},
	})
}

func testAccCheckUpdownWebhookDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "updown_webhook" {
			continue
		}
	}
	return nil
}

func testAccCheckUpdownWebhookExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Webhook ID is set")
		}

		return nil
	}
}

func testAccUpdownWebhookConfig_basic() string {
	return `
resource "updown_webhook" "test" {
  url = "https://example.com/webhook/test"
}
`
}
