package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/thornleyk/graviteeioam-service/client"
	organizationModel "github.com/thornleyk/terraform-provider-graviteeioam/internal/model/organization"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &OrganizationDataSource{}

func NewOrganizationDataSource() datasource.DataSource {
	return &OrganizationDataSource{}
}

type OrganizationDataSource struct {
	client *client.Client
}

func (d *OrganizationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization"
}

func (d *OrganizationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = *organizationModel.GetOrganizationDataSourceSchema()
}

func (d *OrganizationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	var data organizationModel.OrganizationDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	httpRes, err := d.client.OrganizationGetPlatformSettings(ctx, data.OrganizationId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read item",
			err.Error(),
		)
		return
	}

	var apiRes client.Organization
	if httpRes.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP error code received",
			httpRes.Status,
		)
		return
	}

	if err := json.NewDecoder(httpRes.Body).Decode(&apiRes); err != nil {
		resp.Diagnostics.AddError(
			"Invalid format received",
			err.Error(),
		)
		return
	}

	data, mapErr := organizationModel.MapOrganizationDataSource(&apiRes, data)
	if mapErr != nil {
		resp.Diagnostics.AddError(
			"Unable to read data source",
			mapErr.Error(),
		)
		return
	}
	tflog.Trace(ctx, "Read the data source")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
