package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDomainResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + testAccDomainResourceConfig("one"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("graviteeioam_domain.test", "configurable_attribute", "one"),
					resource.TestCheckResourceAttr("graviteeioam_domain.test", "defaulted", "example value when not configured"),
					resource.TestCheckResourceAttr("graviteeioam_domain.test", "id", "example-id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "graviteeioam_domain.test",
				ImportState:       true,
				ImportStateVerify: true,
				// This is not normally necessary, but is here because this
				// example code does not have an actual upstream service.
				// Once the Read method is able to refresh information from
				// the upstream service, this can be removed.
				ImportStateVerifyIgnore: []string{"configurable_attribute", "defaulted"},
			},
			// Update and Read testing
			{
				Config: providerConfig + testAccDomainResourceConfig("two"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("graviteeioam_domain.test", "configurable_attribute", "two"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDomainResourceConfig(configurableAttribute string) string {
	return fmt.Sprintf(`
resource "graviteeioam_domain" "test" {
  configurable_attribute = %[1]q
}
`, configurableAttribute)
}
