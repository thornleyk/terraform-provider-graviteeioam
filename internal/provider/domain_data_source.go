// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/thornleyk/graviteeioam-service/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &DomainDataSource{}

func NewDomainDataSource() datasource.DataSource {
	return &DomainDataSource{}
}

// ExampleDataSource defines the data source implementation.
type DomainDataSource struct {
	client *client.Client
}

// DomainDataSourceModel describes the data source data model.
type DomainDataSourceModel struct {
	ConfigurableAttribute types.String `tfsdk:"configurable_attribute"`
	Id                    types.String `tfsdk:"id"`
}

func (d *DomainDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain"
}

func (d *DomainDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Example data source",

		Attributes: map[string]schema.Attribute{
			"configurable_attribute": schema.StringAttribute{
				MarkdownDescription: "Example configurable attribute",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Example identifier",
				Computed:            true,
			},
		},
	}
}

func (d *DomainDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DomainDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DomainDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	domainResponse, err := d.client.GetOrganizationSettings(ctx, "DEFAULT")
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
			"Unexpected HTTP error code received for Domain",
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

	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.
	data.Id = types.StringValue(*domain.Id)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
