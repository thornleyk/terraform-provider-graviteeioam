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
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "TF identifier",
				Computed:            true,
			},
			"domain_id": schema.StringAttribute{
				MarkdownDescription: "Domain id",
				Required:            true,
			},
			"hrid": schema.StringAttribute{
				MarkdownDescription: "Domain hrid",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Domain Identity name",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Domain Identity type",
				Computed:            true,
			},
			"configuration": schema.StringAttribute{
				MarkdownDescription: "Domain Identity configuration",
				Computed:            true,
			},
			"user_mappers": schema.MapNestedAttribute{
				MarkdownDescription: "Domain Identity user mapping",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"mapping": schema.StringAttribute{
							MarkdownDescription: "Domain Identity user attribute mapping",
							Computed:            true,
						},
					},
				},
				Computed: true,
			},
			"role_mappers": schema.MapNestedAttribute{
				MarkdownDescription: "Domain Identity role mapping",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"mapping": schema.ListAttribute{
							MarkdownDescription: "Domain Identity role map",
							ElementType:         types.StringType,
							Computed:            true,
						},
					},
				},
				Computed: true,
			},
			"reference_type": schema.StringAttribute{
				MarkdownDescription: "Domain Identity reference type",
				Computed:            true,
			},
			"reference_id": schema.StringAttribute{
				MarkdownDescription: "Domain Identity reference id",
				Computed:            true,
			},
			"external": schema.BoolAttribute{
				MarkdownDescription: "Domain Identity exposed externally",
				Computed:            true,
			},
			"whitelist": schema.ListAttribute{
				MarkdownDescription: "Domain Identity whitelist",
				ElementType:         types.StringType,
				Computed:            true,
			},
		},
	}
}
