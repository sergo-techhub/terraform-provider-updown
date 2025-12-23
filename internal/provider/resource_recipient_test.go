package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccUpdownRecipient_email(t *testing.T) {
	email := fmt.Sprintf("test-%s@example.com", acctest.RandString(10))
	resourceName := "updown_recipient.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckUpdownRecipientDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccUpdownRecipientConfig_email(email),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUpdownRecipientExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", "email"),
					resource.TestCheckResourceAttr(resourceName, "value", email),
				),
			},
		},
	})
}

func TestAccUpdownRecipient_webhook(t *testing.T) {
	resourceName := "updown_recipient.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckUpdownRecipientDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccUpdownRecipientConfig_webhook(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUpdownRecipientExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", "webhook"),
					resource.TestCheckResourceAttr(resourceName, "value", "https://example.com/webhook"),
				),
			},
		},
	})
}

func TestAccUpdownRecipient_slackCompatible(t *testing.T) {
	resourceName := "updown_recipient.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckUpdownRecipientDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccUpdownRecipientConfig_slack(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUpdownRecipientExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", "slack_compatible"),
				),
			},
		},
	})
}

func testAccCheckUpdownRecipientDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "updown_recipient" {
			continue
		}
	}
	return nil
}

func testAccCheckUpdownRecipientExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Recipient ID is set")
		}

		return nil
	}
}

func testAccUpdownRecipientConfig_email(email string) string {
	return fmt.Sprintf(`
resource "updown_recipient" "test" {
  type  = "email"
  value = %[1]q
}
`, email)
}

func testAccUpdownRecipientConfig_webhook() string {
	return `
resource "updown_recipient" "test" {
  type  = "webhook"
  value = "https://example.com/webhook"
}
`
}

func testAccUpdownRecipientConfig_slack() string {
	return `
resource "updown_recipient" "test" {
  type  = "slack_compatible"
  value = "https://hooks.slack.com/services/TEST/WEBHOOK/URL"
}
`
}
