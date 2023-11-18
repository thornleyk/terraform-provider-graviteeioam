package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOrganizationDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + testAccOrganizationDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.graviteeioam_organization.test", "id", "DEFAULT"),
				),
			},
		},
	})
}

const testAccOrganizationDataSourceConfig = `
data "graviteeioam_organization" "test" {
  organization_id = "DEFAULT"
}
`
