package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/thornleyk/graviteeioam-service/client"
	"github.com/thornleyk/terraform-provider-graviteeioam/internal/model"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &EnvironmentDataSource{}

func NewEnvironmentDataSource() datasource.DataSource {
	return &EnvironmentDataSource{}
}

// ExampleDataSource defines the data source implementation.
type EnvironmentDataSource struct {
	client *client.Client
}

func (d *EnvironmentDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_environment"
}

func (d *EnvironmentDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
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

func (d *EnvironmentDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func ParseEnvironmentID(id string) (string, string, error) {
	parts := strings.SplitN(id, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("unexpected format of ID (%s), expected organizationId:environmentId", id)
	}
	return parts[0], parts[1], nil
}

func (d *EnvironmentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data model.EnvironmentDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	organizationId, environmentId, idErr := ParseEnvironmentID(data.EnvironmentId.ValueString())
	if idErr != nil {
		resp.Diagnostics.AddError(
			"Error parsing id for Environment Id",
			data.EnvironmentId.ValueString(),
		)
		return
	}

	var listParams = client.ListDomainsParams{}

	domainsResponse, err := d.client.ListDomains(ctx, organizationId, environmentId, &listParams)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Item",
			err.Error(),
		)
		return
	}

	if domainsResponse.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP error code received for Environment",
			domainsResponse.Status,
		)
		return
	}

	var listWrapper model.GraviteeAMListWrapper
	if listWrapperErr := json.NewDecoder(domainsResponse.Body).Decode(&listWrapper); err != nil {
		resp.Diagnostics.AddError(
			"Invalid format received for GraviteeAMListWrapper",
			listWrapperErr.Error(),
		)
		return
	}

	dataJSON, err := json.Marshal(listWrapper.Data)
	if err != nil {
		fmt.Println("Error converting Data to JSON string:", err)
		return
	}

	fmt.Println("Data as JSON string:", string(dataJSON))

	var domains []client.Domain
	if domainsErr := json.Unmarshal(dataJSON, &domains); err != nil {
		fmt.Println("Error re-marshaling JSON string to concrete type:", domainsErr)
		return
	}

	data.Id = data.EnvironmentId

	for _, domain := range domains {
		var domainData = model.DomainLightDataSourceModel{
			Hrid: types.StringValue(*domain.Hrid),
		}
		data.Domains = append(data.Domains, domainData)
	}

	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
