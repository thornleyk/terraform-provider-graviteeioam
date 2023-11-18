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
			"organization_id": schema.StringAttribute{
				MarkdownDescription: "Organiation id",
				Required:            true,
			},
			"environment_id": schema.StringAttribute{
				MarkdownDescription: "Environment id",
				Required:            true,
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
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
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

	var listParams = client.ListDomainsParams{}

	domainsResponse, err := d.client.ListDomains(ctx, data.OrganizationId.ValueString(), data.EnvironmentId.ValueString(), &listParams)
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

	var domains []client.Domain
	if err := json.NewDecoder(domainsResponse.Body).Decode(&domains); err != nil {
		resp.Diagnostics.AddError(
			"Invalid format received for Item",
			err.Error(),
		)
		return
	}

	data.Id = types.StringValue(fmt.Sprintf("%s:%s", data.OrganizationId, data.EnvironmentId))

	for _, domain := range domains {
		var domainData = model.DomainDataSourceModel{
			Id:   types.StringValue(*domain.Id),
			Name: types.StringValue(*domain.Name),
		}
		data.Domains = append(data.Domains, domainData)
	}

	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
