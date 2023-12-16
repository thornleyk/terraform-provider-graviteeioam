package organization_identity

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/thornleyk/graviteeioam-service/client"
)

type OrganizationIdentityDataSourceModel struct {
	Id             types.String `tfsdk:"id"`
	OrganizationId types.String `tfsdk:"organization_id"`
}

func MapOrganizationIdentityDataSource(source *client.IdentityProvider, target OrganizationIdentityDataSourceModel) (OrganizationIdentityDataSourceModel, error) {

	return target, nil
}

func GetOrganizationIdentityDataSourceSchema() *schema.Schema {
	return &schema.Schema{
		MarkdownDescription: "OrganizationIdentity data source",
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
				MarkdownDescription: "Organization Identity name",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Organization Identity type",
				Computed:            true,
			},
			"system": schema.BoolAttribute{
				MarkdownDescription: "Organization Identity system provided identity",
				Computed:            true,
			},
			"configuration": schema.StringAttribute{
				MarkdownDescription: "Organization Identity configuration",
				Computed:            true,
			},
			"user_mappers": schema.MapNestedAttribute{
				MarkdownDescription: "Organization Identity user mapping",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"mapping": schema.StringAttribute{
							MarkdownDescription: "Organization Identity user attribute mapping",
							Computed:            true,
						},
					},
				},
				Computed: true,
			},
			"role_mappers": schema.MapNestedAttribute{
				MarkdownDescription: "Organization Identity role mapping",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"mapping": schema.ListAttribute{
							MarkdownDescription: "Organization Identity role map",
							ElementType:         types.StringType,
							Computed:            true,
						},
					},
				},
				Computed: true,
			},
			"reference_type": schema.StringAttribute{
				MarkdownDescription: "Organization Identity reference type",
				Computed:            true,
			},
			"reference_id": schema.StringAttribute{
				MarkdownDescription: "Organization Identity reference id",
				Computed:            true,
			},
			"external": schema.BoolAttribute{
				MarkdownDescription: "Organization Identity exposed externally",
				Computed:            true,
			},
			"whitelist": schema.ListAttribute{
				MarkdownDescription: "Organization Identity whitelist",
				ElementType:         types.StringType,
				Computed:            true,
			},
		},
	}
}
