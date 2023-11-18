package model

import "github.com/hashicorp/terraform-plugin-framework/types"

// EnvironmentDataSourceModel describes the data source data model.
type EnvironmentDataSourceModel struct {
	Id             types.String            `tfsdk:"id"`
	OrganizationId types.String            `tfsdk:"organization_id"`
	EnvironmentId  types.String            `tfsdk:"environment_id"`
	Domains        []DomainDataSourceModel `tfsdk:"domains"`
}
