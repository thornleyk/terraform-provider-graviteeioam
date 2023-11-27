package organization_identity

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/thornleyk/graviteeioam-service/client"
)

type OrganizationIdentityDataSourceModel struct {
	Id       types.String `tfsdk:"id"`
	DomainId types.String `tfsdk:"domain_id"`
}

func MapOrganizationIdentityDataSource(source *client.IdentityProvider, target OrganizationIdentityDataSourceModel) (OrganizationIdentityDataSourceModel, error) {

	return target, nil
}

func GetOrganizationIdentityDataSourceSchema() *schema.Schema {
	return &schema.Schema{
		MarkdownDescription: "OrganizationIdentity data source",
	}
}
