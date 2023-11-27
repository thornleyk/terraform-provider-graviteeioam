package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/thornleyk/graviteeioam-service/client"
	domainIdentityModel "github.com/thornleyk/terraform-provider-graviteeioam/internal/model/domain_identity"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &DomainIdentityProviderDataSource{}

func NewDomainIdentityProviderDataSource() datasource.DataSource {
	return &DomainIdentityProviderDataSource{}
}

type DomainIdentityProviderDataSource struct {
	client *client.Client
}

func (d *DomainIdentityProviderDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain_identity_provider"
}

func (d *DomainIdentityProviderDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = *domainIdentityModel.GetDomainIdentityDataSourceSchema()
}

func (d *DomainIdentityProviderDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func ParseDomainIdentityProviderID(id string) (string, string, string, string, error) {
	parts := strings.SplitN(id, ":", 4)
	if len(parts) != 4 || parts[0] == "" || parts[1] == "" || parts[2] == "" || parts[3] == "" {
		return "", "", "", "", fmt.Errorf("unexpected format of ID (%s), expected organizationId:environmentId:domainId:identityProviderId", id)
	}
	return parts[0], parts[1], parts[2], parts[3], nil
}

func (d *DomainIdentityProviderDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data domainIdentityModel.DomainIdentityDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	organizationId, environmentId, domainId, identityProviderId, idErr := ParseDomainIdentityProviderID(data.DomainId.ValueString())
	if idErr != nil {
		resp.Diagnostics.AddError(
			"Error parsing id",
			data.DomainId.String(),
		)
		return
	}

	httpRes, err := d.client.DomainGetIdentityProvider(ctx, organizationId, environmentId, domainId, identityProviderId)
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
	data, mapErr := domainIdentityModel.MapDomainIdentityDataSource(&apiRes, data)
	if mapErr != nil {
		resp.Diagnostics.AddError(
			"Unable to read data source",
			mapErr.Error(),
		)
		return
	}

	tflog.Trace(ctx, "read a data source")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
