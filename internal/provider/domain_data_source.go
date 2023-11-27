package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/thornleyk/graviteeioam-service/client"
	domainModel "github.com/thornleyk/terraform-provider-graviteeioam/internal/model/domain"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &DomainDataSource{}

func NewDomainDataSource() datasource.DataSource {
	return &DomainDataSource{}
}

type DomainDataSource struct {
	client *client.Client
}

func (d *DomainDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain"
}

func (d *DomainDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = *domainModel.GetDomainDataSourceSchema()
}

func ParseDomainID(id string) (string, string, string, error) {
	parts := strings.SplitN(id, ":", 3)
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return "", "", "", fmt.Errorf("unexpected format of ID (%s), expected organizationId:environmentId:domainId", id)
	}
	return parts[0], parts[1], parts[2], nil
}

func (d *DomainDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DomainDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data domainModel.DomainDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	organizationId, environmentId, domainId, idErr := ParseDomainID(data.DomainId.ValueString())
	if idErr != nil {
		resp.Diagnostics.AddError(
			"Error parsing id",
			data.DomainId.String(),
		)
		return
	}

	httpRes, err := d.client.DomainGetByHrid(ctx, organizationId, environmentId, domainId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read item",
			err.Error(),
		)
		return
	}

	var apiRes client.Domain
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
	data, mapErr := domainModel.MapDomainDataSource(&apiRes, data)
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
