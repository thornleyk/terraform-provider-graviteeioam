package environment

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/thornleyk/graviteeioam-service/client"
)

type DomainLightDataSourceModel struct {
	Hrid types.String `tfsdk:"hrid"`
}

type EnvironmentDataSourceModel struct {
	Id            types.String                 `tfsdk:"id"`
	EnvironmentId types.String                 `tfsdk:"environment_id"`
	Domains       []DomainLightDataSourceModel `tfsdk:"domains"`
}

func MapEnvironmentDataSource(source []client.Domain, target EnvironmentDataSourceModel) (EnvironmentDataSourceModel, error) {
	target.Id = target.EnvironmentId
	for _, domain := range source {
		var domainData = DomainLightDataSourceModel{
			Hrid: types.StringValue(*domain.Hrid),
		}
		target.Domains = append(target.Domains, domainData)
	}
	return target, nil
}

func GetEnvironmentDataSourceSchema() *schema.Schema {
	return &schema.Schema{
		MarkdownDescription: "Environment data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "TF identifier",
				Computed:            true,
			},
			"environment_id": schema.StringAttribute{
				MarkdownDescription: "Environment id",
				Required:            true,
			},
			"domains": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"hrid": schema.StringAttribute{
							Computed: true,
						},
					},
				},
				Computed: true,
			},
		},
	}
}
