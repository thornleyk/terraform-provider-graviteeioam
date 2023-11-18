package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/thornleyk/graviteeioam-service/client"
)

type DomainClientRegistrationSettings struct {
	AllowLocalhostRedirectURI          types.Bool `tfsdk:"allow_localhost_redirect_uri"`
	AllowHTTPSchemeRedirectURI         types.Bool `tfsdk:"allow_http_scheme_redirect_uri"`
	AllowWildCardRedirectURI           types.Bool `tfsdk:"allow_wild_card_redirect_uri"`
	IsDynamicClientRegistrationEnabled types.Bool `tfsdk:"is_dynamic_client_registration_enabled"`
	IsAllowedScopesEnabled             types.Bool `tfsdk:"is_allowed_scopes_enabled"`
	IsClientTemplateEnabled            types.Bool `tfsdk:"is_client_template_enabled"`
}

type DomainSecurityProfileSettings struct {
	EnablePlanFAPI   types.Bool `tfsdk:"enable_plain_fapi"`
	EnableFAPIBrazil types.Bool `tfsdk:"enable_fapi_brazil"`
}

type DomainCIBASettings struct {
	Enabled              types.Bool   `tfsdk:"enabled"`
	AuthReqExpiry        types.Number `tfsdk:"auth_req_expiry"`
	TokenReqInterval     types.Number `tfsdk:"token_req_interval"`
	BindingMessageLength types.Number `tfsdk:"binding_message_length"`
}

type DomainOIDC struct {
	ClientRegistrationSettings types.String `tfsdk:"client_registration_settings"`
	SecurityProfileSettings    types.String `tfsdk:"security_profile_settings"`
	RedirectURIStrictMatching  types.Bool   `tfsdk:"redirect_uri_strict_matching"`
	CIBASettings               types.String `tfsdk:"ciba_settings"`
}

type DomainLoginSettings struct {
	Inherited                          types.String `tfsdk:"inherited"`
	ForgotPasswordEnabled              types.String `tfsdk:"forgot_password_enabled"`
	RegisterEnabled                    types.Bool   `tfsdk:"register_enabled"`
	RememberMeEnabled                  types.Bool   `tfsdk:"remember_me_enabled"`
	PasswordlessEnabled                types.String `tfsdk:"passwordless_enabled"`
	PasswordlessRememberDeviceEnabled  types.String `tfsdk:"passwordless_remember_device_enabled"`
	PasswordlessEnforcePasswordEnabled types.String `tfsdk:"passwordless_enforce_password_enabled"`
	PasswordlessDeviceNamingEnabled    types.String `tfsdk:"passwordless_device_naming_enabled"`
	HideForm                           types.String `tfsdk:"hide_form"`
	IdentifierFirstEnabled             types.String `tfsdk:"identifier_first_enabled"`
}

type DomainDataSourceModel struct {
	Id                  types.String         `tfsdk:"id"`
	DomainId            types.String         `tfsdk:"domain_id"`
	Hrid                types.String         `tfsdk:"hrid"`
	Name                types.String         `tfsdk:"name"`
	Description         types.String         `tfsdk:"description"`
	Enabled             types.Bool           `tfsdk:"enabled"`
	Master              types.Bool           `tfsdk:"master"`
	VHostMode           types.Bool           `tfsdk:"vhost_mode"`
	DomainOIDC          *DomainOIDC          `tfsdk:"oidc"`
	DomainLoginSettings *DomainLoginSettings `tfsdk:"login_settings"`
}

func MapDomainDataSource(source *client.Domain, target DomainDataSourceModel) (DomainDataSourceModel, error) {
	target.Id = target.DomainId
	target.Hrid = types.StringValue(*source.Hrid)
	target.Name = types.StringValue(*source.Name)
	target.Description = types.StringValue(*source.Description)
	target.Enabled = types.BoolValue(*source.Enabled)
	target.Master = types.BoolValue(*source.Master)
	target.VHostMode = types.BoolValue(*source.VhostMode)
	return target, nil
}