package provider

import (
	"context"
	"encoding/json"
	"os"

	"github.com/thornleyk/graviteeioam-service/client"

	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure ScaffoldingProvider satisfies various provider interfaces.
var _ provider.Provider = &GraviteeIOAMProvider{}

// ScaffoldingProvider defines the provider implementation.
type GraviteeIOAMProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// ScaffoldingProviderModel describes the provider data model.
type GraviteeIOAMProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

func (p *GraviteeIOAMProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "graviteeioam"
	resp.Version = p.version
}

func (p *GraviteeIOAMProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "GraviteeIO AM endpoint",
				Optional:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "GraviteeIO AM username",
				Optional:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "GraviteeIO AM password",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *GraviteeIOAMProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring GraviteeIOAM client")

	var config GraviteeIOAMProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown GraviteeIOAM username",
			"The provider cannot create the GraviteeIOAM API client as there is an unknown configuration value for the GraviteeIOAM host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the GRAVITEEIOAM_USERNAME environment variable.",
		)
	}

	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown GraviteeIOAM password",
			"The provider cannot create the GraviteeIOAM API client as there is an unknown configuration value for the GraviteeIOAM password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the GRAVITEEIOAM_PASSWORD environment variable.",
		)
	}

	if config.Endpoint.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Unknown GraviteeIOAM endpoint",
			"The provider cannot create the GraviteeIOAM API client as there is an unknown configuration value for the GraviteeIOAM endpoint. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the GRAVITEEIOAM_ENDPOINT environment variable.",
		)
	}

	username := os.Getenv("GRAVITEEIOAM_USERNAME")
	password := os.Getenv("GRAVITEEIOAM_PASSWORD")
	endpoint := os.Getenv("GRAVITEEIOAM_ENDPOINT")

	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	if !config.Endpoint.IsNull() {
		endpoint = config.Endpoint.ValueString()
	}

	if username == "" {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("username"),
			"Missing GraviteeIOAM API username (using default value: admin)",
			"The provider is using a default value as there is a missing or empty value for the Inventory API host. "+
				"Set the host value in the configuration or use the GRAVITEEIOAM_USERNAME environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
		username = "admin"
	}

	if password == "" {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("password"),
			"Missing GraviteeIOAM API password (using default value: admin)",
			"The provider is using a default value as there is a missing or empty value for the GraviteeIOAM API host. "+
				"Set the host value in the configuration or use the GRAVITEEIOAM_USERNAME environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
		username = "admin"
	}

	if endpoint == "" {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("endpoint"),
			"Missing GraviteeIOAM API endpoint (using default value: https://localhost:8093/management/)",
			"The provider is using a default value as there is a missing or empty value for the GraviteeIOAM API host. "+
				"Set the host value in the configuration or use the GRAVITEEIOAM_USERNAME environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
		endpoint = "https://localhost:8093/management/"
	}

	if resp.Diagnostics.HasError() {
		return
	}

	basicAuthProvider, basicAuthProviderErr := securityprovider.NewSecurityProviderBasicAuth(username, password)
	if basicAuthProviderErr != nil {
		panic(basicAuthProviderErr)
	}
	authApi, authApiErr := client.NewClient(endpoint, client.WithRequestEditorFn(basicAuthProvider.Intercept))

	if authApiErr != nil {
		resp.Diagnostics.AddError(
			"Unable to Create GraviteeIOAM API Client",
			"An unexpected error occurred when creating the GraviteeIOAM API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"GraviteeIOAM Client Error: "+authApiErr.Error(),
		)
		return
	}

	authToken, authTokenErr := authApi.AuthToken(ctx)

	//authResp, _ := io.ReadAll(authToken.Body)

	if authTokenErr != nil {
		resp.Diagnostics.AddError(
			"Unable to Create GraviteeIOAM API Client",
			"An unexpected error occurred when authenticating the GraviteeIOAM API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"GraviteeIOAM Client Error: "+authTokenErr.Error(),
		)
		return
	}

	if authToken.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP error code received for Authentication",
			authToken.Status,
		)
		return
	}

	//tflog.Info(ctx, "Authenticated GraviteeIOAM client "+string(authResp))

	var token client.AuthToken
	if err := json.NewDecoder(authToken.Body).Decode(&token); err != nil {
		resp.Diagnostics.AddError(
			"Invalid format received for AuthToken",
			err.Error(),
		)
		return
	}

	bearerTokenProvider, bearerTokenProviderErr := securityprovider.NewSecurityProviderBearerToken(*token.AccessToken)
	if bearerTokenProviderErr != nil {
		panic(bearerTokenProviderErr)
	}
	api, apiErr := client.NewClient(endpoint, client.WithRequestEditorFn(bearerTokenProvider.Intercept))

	if apiErr != nil {
		resp.Diagnostics.AddError(
			"Unable to Create GraviteeIOAM API Client",
			"An unexpected error occurred when creating the GraviteeIOAM API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"GraviteeIOAM Client Error: "+apiErr.Error(),
		)
		return
	}

	resp.DataSourceData = api
	resp.ResourceData = api
	tflog.Info(ctx, "Configured GraviteeIOAM client", map[string]any{"success": true})

}

func (p *GraviteeIOAMProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewDomainResource,
	}
}

func (p *GraviteeIOAMProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewDomainDataSource,
		NewOrganizationDataSource,
		NewEnvironmentDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &GraviteeIOAMProvider{
			version: version,
		}
	}
}
