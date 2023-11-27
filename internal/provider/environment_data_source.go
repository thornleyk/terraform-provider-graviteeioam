package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/thornleyk/graviteeioam-service/client"
	environmentModel "github.com/thornleyk/terraform-provider-graviteeioam/internal/model/environment"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &EnvironmentDataSource{}

func NewEnvironmentDataSource() datasource.DataSource {
	return &EnvironmentDataSource{}
}

type EnvironmentDataSource struct {
	client *client.Client
}

func (d *EnvironmentDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_environment"
}

func (d *EnvironmentDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = *environmentModel.GetEnvironmentDataSourceSchema()
}

func (d *EnvironmentDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	var data environmentModel.EnvironmentDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	organizationId, environmentId, idErr := ParseEnvironmentID(data.EnvironmentId.ValueString())
	if idErr != nil {
		resp.Diagnostics.AddError(
			"Error parsing id",
			data.EnvironmentId.ValueString(),
		)
		return
	}

	var listParams = client.EnvironmentListDomainsPaginatedParams{}

	httpRes, err := d.client.EnvironmentListDomainsPaginated(ctx, organizationId, environmentId, &listParams)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read item",
			err.Error(),
		)
		return
	}

	if httpRes.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Invalid format received",
			httpRes.Status,
		)
		return
	}

	var listWrapper client.Page
	if listWrapperErr := json.NewDecoder(httpRes.Body).Decode(&listWrapper); err != nil {
		resp.Diagnostics.AddError(
			"Invalid list format received",
			listWrapperErr.Error(),
		)
		return
	}

	dataJSON, err := json.Marshal(listWrapper.Data)
	if err != nil {
		fmt.Println("Error converting Data to JSON string:", err)
		return
	}

	var apiRes []client.Domain
	if marshalErr := json.Unmarshal(dataJSON, &apiRes); marshalErr != nil {
		fmt.Println("Error re-marshaling JSON string to concrete type:", marshalErr)
		return
	}

	data, mapErr := environmentModel.MapEnvironmentDataSource(apiRes, data)
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
