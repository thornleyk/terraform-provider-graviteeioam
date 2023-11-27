package domain_identity

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/thornleyk/graviteeioam-service/client"
)

type DomainIdentityDataSourceModel struct {
	Id       types.String `tfsdk:"id"`
	DomainId types.String `tfsdk:"domain_id"`
}

func MapDomainIdentityDataSource(source *client.IdentityProvider, target DomainIdentityDataSourceModel) (DomainIdentityDataSourceModel, error) {

	return target, nil
}

func GetDomainIdentityDataSourceSchema() *schema.Schema {
	return &schema.Schema{
		MarkdownDescription: "DomainIdentity data source",
	}
}
