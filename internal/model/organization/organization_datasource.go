package organization

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/thornleyk/graviteeioam-service/client"
)

type OrganizationDataSourceModel struct {
	Id             types.String   `tfsdk:"id"`
	OrganizationId types.String   `tfsdk:"organization_id"`
	Name           types.String   `tfsdk:"name"`
	Identities     []types.String `tfsdk:"identities"`
	HrIds          []types.String `tfsdk:"hrids"`
}

func MapOrganizationDataSource(source *client.Organization, target OrganizationDataSourceModel) (OrganizationDataSourceModel, error) {
	target.Id = target.OrganizationId
	target.Name = types.StringValue(*source.Name)

	for _, hrid := range *source.Hrids {
		target.HrIds = append(target.Identities, types.StringValue(hrid))
	}

	for _, identity := range *source.Identities {
		target.Identities = append(target.Identities, types.StringValue(identity))
	}
	return target, nil
}

func GetOrganizationDataSourceSchema() *schema.Schema {
	return &schema.Schema{
		MarkdownDescription: "Organization data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "TF identifier",
				Computed:            true,
			},
			"organization_id": schema.StringAttribute{
				MarkdownDescription: "Organization id",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Organization name",
				Computed:            true,
			},
			"identities": schema.ListAttribute{
				MarkdownDescription: "Organization identities",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"hrids": schema.ListAttribute{
				MarkdownDescription: "Organization HrIds",
				ElementType:         types.StringType,
				Computed:            true,
			},
		},
	}
}
