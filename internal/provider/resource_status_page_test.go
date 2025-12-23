package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccUpdownStatusPage_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	resourceName := "updown_status_page.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckUpdownStatusPageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccUpdownStatusPageConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUpdownStatusPageExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "visibility", "private"),
					resource.TestCheckResourceAttrSet(resourceName, "url"),
				),
			},
		},
	})
}

func TestAccUpdownStatusPage_public(t *testing.T) {
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	resourceName := "updown_status_page.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckUpdownStatusPageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccUpdownStatusPageConfig_public(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUpdownStatusPageExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "visibility", "public"),
					resource.TestCheckResourceAttr(resourceName, "description", "Test public status page"),
				),
			},
		},
	})
}

func TestAccUpdownStatusPage_protected(t *testing.T) {
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	resourceName := "updown_status_page.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckUpdownStatusPageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccUpdownStatusPageConfig_protected(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUpdownStatusPageExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "visibility", "protected"),
					resource.TestCheckResourceAttr(resourceName, "access_key", "test-access-key-123"),
				),
			},
		},
	})
}

func TestAccUpdownStatusPage_update(t *testing.T) {
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	rNameUpdated := rName + "-updated"
	resourceName := "updown_status_page.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckUpdownStatusPageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccUpdownStatusPageConfig_update(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUpdownStatusPageExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
			{
				Config: testAccUpdownStatusPageConfig_update(rName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUpdownStatusPageExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
				),
			},
		},
	})
}

func testAccCheckUpdownStatusPageDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "updown_status_page" {
			continue
		}
	}
	return nil
}

func testAccCheckUpdownStatusPageExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Status Page ID is set")
		}

		return nil
	}
}

func testAccUpdownStatusPageConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "updown_check" "test" {
  url   = "https://example.com"
  alias = "%[1]s-check"
}

resource "updown_status_page" "test" {
  name       = %[1]q
  visibility = "private"
  checks     = [updown_check.test.id]
}
`, rName)
}

func testAccUpdownStatusPageConfig_update(checkName, statusPageName string) string {
	return fmt.Sprintf(`
resource "updown_check" "test" {
  url   = "https://example.com"
  alias = "%[1]s-check"
}

resource "updown_status_page" "test" {
  name       = %[2]q
  visibility = "private"
  checks     = [updown_check.test.id]
}
`, checkName, statusPageName)
}

func testAccUpdownStatusPageConfig_public(rName string) string {
	return fmt.Sprintf(`
resource "updown_check" "test" {
  url   = "https://example.com"
  alias = "%[1]s-check"
}

resource "updown_status_page" "test" {
  name        = %[1]q
  description = "Test public status page"
  visibility  = "public"
  checks      = [updown_check.test.id]
}
`, rName)
}

func testAccUpdownStatusPageConfig_protected(rName string) string {
	return fmt.Sprintf(`
resource "updown_check" "test" {
  url   = "https://example.com"
  alias = "%[1]s-check"
}

resource "updown_status_page" "test" {
  name       = %[1]q
  visibility = "protected"
  access_key = "test-access-key-123"
  checks     = [updown_check.test.id]
}
`, rName)
}
