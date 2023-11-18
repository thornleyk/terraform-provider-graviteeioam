package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEnvironmentDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + testAccEnvironmentDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.graviteeioam_environment.test", "id", "DEFAULT:DEFAULT"),
				),
			},
		},
	})
}

const testAccEnvironmentDataSourceConfig = `
data "graviteeioam_environment" "test" {
  environment_id = "DEFAULT:DEFAULT"
}
`
