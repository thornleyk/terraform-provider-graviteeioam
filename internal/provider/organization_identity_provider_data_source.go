package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/thornleyk/graviteeioam-service/client"
	organizationIdentityModel "github.com/thornleyk/terraform-provider-graviteeioam/internal/model/organization_identity"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &OrganizationIdentityProviderDataSource{}

func NewOrganizationIdentityProviderDataSource() datasource.DataSource {
	return &OrganizationIdentityProviderDataSource{}
}

type OrganizationIdentityProviderDataSource struct {
	client *client.Client
}

func (d *OrganizationIdentityProviderDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_identity_provider"
}

func (d *OrganizationIdentityProviderDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = *organizationIdentityModel.GetOrganizationIdentityDataSourceSchema()
}

func (d *OrganizationIdentityProviderDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func ParseOrganizationIdentityProviderID(id string) (string, string, error) {
	parts := strings.SplitN(id, ":", 3)
	if len(parts) != 4 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return "", "", fmt.Errorf("unexpected format of ID (%s), expected organizationId:identityProviderId", id)
	}
	return parts[0], parts[1], nil
}

func (d *OrganizationIdentityProviderDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data organizationIdentityModel.OrganizationIdentityDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	organizationId, identityProviderId, idErr := ParseOrganizationIdentityProviderID(data.OrganizationId.ValueString())
	if idErr != nil {
		resp.Diagnostics.AddError(
			"Error parsing id",
			data.OrganizationId.String(),
		)
		return
	}

	httpRes, err := d.client.OrganizationGetPlatformIdentityProvider(ctx, organizationId, identityProviderId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read item",
			err.Error(),
		)
		return
	}

	var apiRes client.IdentityProvider
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
	data, mapErr := organizationIdentityModel.MapOrganizationIdentityDataSource(&apiRes, data)
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
