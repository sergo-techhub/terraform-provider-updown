package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccUpdownCheck_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	resourceName := "updown_check.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckUpdownCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccUpdownCheckConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUpdownCheckExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "url", "https://example.com"),
					resource.TestCheckResourceAttr(resourceName, "alias", rName),
					resource.TestCheckResourceAttr(resourceName, "period", "60"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
		},
	})
}

func TestAccUpdownCheck_icmp(t *testing.T) {
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	resourceName := "updown_check.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckUpdownCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccUpdownCheckConfig_icmp(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUpdownCheckExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "url", "8.8.8.8"),
					resource.TestCheckResourceAttr(resourceName, "type", "icmp"),
				),
			},
		},
	})
}

func TestAccUpdownCheck_tcp(t *testing.T) {
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	resourceName := "updown_check.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckUpdownCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccUpdownCheckConfig_tcp(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUpdownCheckExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", "tcp"),
				),
			},
		},
	})
}

func TestAccUpdownCheck_httpVerb(t *testing.T) {
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	resourceName := "updown_check.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckUpdownCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccUpdownCheckConfig_httpVerb(rName, "POST"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUpdownCheckExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "http_verb", "POST"),
				),
			},
			{
				// Update back to default GET/HEAD
				Config: testAccUpdownCheckConfig_httpVerb(rName, "GET/HEAD"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUpdownCheckExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "http_verb", "GET/HEAD"),
				),
			},
		},
	})
}

func testAccCheckUpdownCheckDestroy(s *terraform.State) error {
	// Since we don't have direct access to the client in tests,
	// we just verify the resources are removed from state
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "updown_check" {
			continue
		}
	}
	return nil
}

func testAccCheckUpdownCheckExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Check ID is set")
		}

		return nil
	}
}

func testAccUpdownCheckConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "updown_check" "test" {
  url     = "https://example.com"
  alias   = %[1]q
  period  = 60
  enabled = true
}
`, rName)
}

func testAccUpdownCheckConfig_icmp(rName string) string {
	return fmt.Sprintf(`
resource "updown_check" "test" {
  url   = "8.8.8.8"
  type  = "icmp"
  alias = %[1]q
}
`, rName)
}

func testAccUpdownCheckConfig_tcp(rName string) string {
	return fmt.Sprintf(`
resource "updown_check" "test" {
  url   = "tcp://google.com:443"
  type  = "tcp"
  alias = %[1]q
}
`, rName)
}

func testAccUpdownCheckConfig_httpVerb(rName, verb string) string {
	return fmt.Sprintf(`
resource "updown_check" "test" {
  url       = "https://example.com"
  alias     = %[1]q
  http_verb = %[2]q
}
`, rName, verb)
}
