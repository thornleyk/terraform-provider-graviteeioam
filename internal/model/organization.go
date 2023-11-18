package model

import "github.com/hashicorp/terraform-plugin-framework/types"

type OrganizationDataSourceModel struct {
	Id             types.String   `tfsdk:"id"`
	OrganizationId types.String   `tfsdk:"organization_id"`
	Name           types.String   `tfsdk:"name"`
	Identities     []types.String `tfsdk:"identities"`
	HrId           types.String   `tfsdk:"hrid"`
}
