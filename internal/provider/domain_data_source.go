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

func (d *DomainDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain"
}

func (d *DomainDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Domain data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "TF identifier",
				Computed:            true,
			},
			"domain_id": schema.StringAttribute{
				MarkdownDescription: "Domain id",
				Required:            true,
			},
			"hrid": schema.StringAttribute{
				MarkdownDescription: "Domain hrid",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Domain name",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Domain description",
				Computed:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Domain enabled",
				Computed:            true,
			},
			"master": schema.BoolAttribute{
				MarkdownDescription: "Domain master",
				Computed:            true,
			},
			"vhost_mode": schema.BoolAttribute{
				MarkdownDescription: "Domain vhost_mode",
				Computed:            true,
			},
			"oidc": schema.SetNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"client_registration_settings": schema.SetNestedAttribute{
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"allow_localhost_redirect_uri": schema.BoolAttribute{
										Computed: true,
									},
									"allow_http_scheme_redirect_uri": schema.BoolAttribute{
										Computed: true,
									},
									"allow_wild_card_redirect_uri": schema.BoolAttribute{
										Computed: true,
									},
									"is_dynamic_client_registration_enabled": schema.BoolAttribute{
										Computed: true,
									},
									"is_allowed_scopes_enabled": schema.BoolAttribute{
										Computed: true,
									},
									"is_client_template_enabled": schema.BoolAttribute{
										Computed: true,
									},
								},
							},
							Computed: true,
						},
						"security_profile_settings": schema.SetNestedAttribute{
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"enable_plain_fapi": schema.BoolAttribute{
										Computed: true,
									},
									"enable_fapi_brazil": schema.BoolAttribute{
										Computed: true,
									},
								},
							},
							Computed: true,
						},
						"redirect_uri_strict_matching": schema.BoolAttribute{
							Computed: true,
						},
						"ciba_settings": schema.SetNestedAttribute{
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"enabled": schema.BoolAttribute{
										Computed: true,
									},
									"auth_req_expiry": schema.NumberAttribute{
										Computed: true,
									},
									"token_req_interval": schema.NumberAttribute{
										Computed: true,
									},
									"binding_message_length": schema.NumberAttribute{
										Computed: true,
									},
								},
							},
							Computed: true,
						},
					},
				},
				Computed: true,
			},
			"login_settings": schema.SetNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"inherited": schema.BoolAttribute{
							Computed: true,
						},
						"forgot_password_enabled": schema.BoolAttribute{
							Computed: true,
						},
						"register_enabled": schema.BoolAttribute{
							Computed: true,
						},
						"remember_me_enabled": schema.BoolAttribute{
							Computed: true,
						},
						"passwordless_enabled": schema.BoolAttribute{
							Computed: true,
						},
						"passwordless_remember_device_enabled": schema.BoolAttribute{
							Computed: true,
						},
						"passwordless_enforce_password_enabled": schema.BoolAttribute{
							Computed: true,
						},
						"passwordless_device_naming_enabled": schema.BoolAttribute{
							Computed: true,
						},
						"hide_form": schema.BoolAttribute{
							Computed: true,
						},
						"identifier_first_enabled": schema.BoolAttribute{
							Computed: true,
						},
					},
				},
				Computed: true,
			},
		},
	}
}

func ParseDomainID(id string) (string, string, string, error) {
	parts := strings.SplitN(id, ":", 3)
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return "", "", "", fmt.Errorf("unexpected format of ID (%s), expected organizationId:environmentId:domainId", id)
	}
	return parts[0], parts[1], parts[2], nil
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
	var data model.DomainDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	organizationId, environmentId, domainId, idErr := ParseDomainID(data.DomainId.ValueString())
	if idErr != nil {
		resp.Diagnostics.AddError(
			"Error parsing id for Domain Id",
			data.DomainId.String(),
		)
		return
	}

	domainResponse, err := d.client.DomainGetByHrid(ctx, organizationId, environmentId, domainId)
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
			"Invalid format received for Domain",
			err.Error(),
		)
		return
	}
	data, mapErr := model.MapDomainDataSource(&domain, data)
	if mapErr != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Item",
			mapErr.Error(),
		)
		return
	}
	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
