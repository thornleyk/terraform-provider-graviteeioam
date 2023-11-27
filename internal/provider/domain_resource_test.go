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
			{
				Config: providerConfig + testAccDomainResourceConfig("one"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("graviteeioam_domain.test", "configurable_attribute", "one"),
					resource.TestCheckResourceAttr("graviteeioam_domain.test", "defaulted", "example value when not configured"),
					resource.TestCheckResourceAttr("graviteeioam_domain.test", "id", "example-id"),
				),
			},
			{
				ResourceName:            "graviteeioam_domain.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"configurable_attribute", "defaulted"},
			},
			{
				Config: providerConfig + testAccDomainResourceConfig("two"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("graviteeioam_domain.test", "configurable_attribute", "two"),
				),
			},
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
