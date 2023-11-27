package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDomainDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccDomainDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.graviteeioam_domain.test", "id", "DEFAULT:DEFAULT:test-domain"),
				),
			},
		},
	})
}

const testAccDomainDataSourceConfig = `
data "graviteeioam_domain" "test" {
  domain_id = "DEFAULT:DEFAULT:test-domain"
}
`
