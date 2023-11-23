package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/thornleyk/graviteeioam-service/client"
	"github.com/thornleyk/terraform-provider-graviteeioam/internal/model"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &OrganizationDataSource{}

func NewOrganizationDataSource() datasource.DataSource {
	return &OrganizationDataSource{}
}

// ExampleDataSource defines the data source implementation.
type OrganizationDataSource struct {
	client *client.Client
}

func (d *OrganizationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization"
}

func (d *OrganizationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
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
			"hrid": schema.StringAttribute{
				MarkdownDescription: "Organization Hrid",
				Computed:            true,
			},
		},
	}
}

func (d *OrganizationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *OrganizationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data model.OrganizationDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	domainResponse, err := d.client.OrganizationGetPlatformSettings(ctx, data.OrganizationId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Item",
			err.Error(),
		)
		return
	}

	var domain client.Domain
	if domainResponse.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP error code received for Organization",
			domainResponse.Status,
		)
		return
	}

	if err := json.NewDecoder(domainResponse.Body).Decode(&domain); err != nil {
		resp.Diagnostics.AddError(
			"Invalid format received for Item",
			err.Error(),
		)
		return
	}

	data.OrganizationId = types.StringValue(*domain.Id)
	data.Id = types.StringValue(*domain.Id)
	data.Name = types.StringValue(*domain.Name)
	//data.HrId = types.StringValue(*domain.Hrid) //FIXME: this is returning from Gravitee as a [] hrid
	for _, identity := range *domain.Identities {
		data.Identities = append(data.Identities, types.StringValue(identity))
	}
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
