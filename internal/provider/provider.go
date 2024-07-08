package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-privx/internal/client"
)

// Ensure privxProvider satisfies various provider interfaces.
var _ provider.Provider = &privxProvider{}

// privxProvider defines the provider implementation.
type privxProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// privxProviderModel describes the provider data model.
type privxProviderModel struct {
	APIBaseURL        types.String `tfsdk:"api_base_url"`
	APIBearerToken    types.String `tfsdk:"api_bearer_token"`
	APIClientID       types.String `tfsdk:"api_client_id"`
	APIClientSecret   types.String `tfsdk:"api_client_secret"`
	OAuthClientID     types.String `tfsdk:"api_oauth_client_id"`
	OAuthClientSecret types.String `tfsdk:"api_oauth_client_secret"`
	Debug             types.Bool   `tfsdk:"debug"`
}

func (p *privxProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "privx"
	resp.Version = p.version
}

func (p *privxProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_base_url": schema.StringAttribute{
				MarkdownDescription: "PrivX API Base URL",
				Optional:            true,
			},
			"api_bearer_token": schema.StringAttribute{
				MarkdownDescription: "PrivX bearer token",
				Optional:            true,
				Sensitive:           true,
			},
			"api_client_id": schema.StringAttribute{
				MarkdownDescription: "PrivX API OAuth client ID",
				Optional:            true,
			},
			"api_client_secret": schema.StringAttribute{
				MarkdownDescription: "PrivX API OAuth client ID",
				Optional:            true,
				Sensitive:           true,
			},
			"api_oauth_client_id": schema.StringAttribute{
				MarkdownDescription: "PrivX API OAuth client ID",
				Optional:            true,
			},
			"api_oauth_client_secret": schema.StringAttribute{
				MarkdownDescription: "Privx API OAuth Client Secret",
				Optional:            true,
				Sensitive:           true,
			},
			"debug": schema.BoolAttribute{
				MarkdownDescription: "PrivX debug mode",
				Optional:            true,
			},
		},
	}
}

func (p *privxProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data privxProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if data.APIBaseURL.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_base_url"),
			"Unknown PrivX API URL",
			"The provider cannot create the PrivX API client as there is an unknown configuration value for the PrivX API base URL. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the PRIVX_API_BASE_URL environment variable.",
		)
	}

	if data.APIBearerToken.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_bearer_token"),
			"Unknown PrivX API bearer token",
			"The provider cannot create the PrivX API client as there is an unknown configuration value for the PrivX API bearer token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the PRIVX_API_BEARER_TOKEN environment variable.",
		)
	}

	if data.APIClientID.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_client_id"),
			"Unknown PrivX API client ID",
			"The provider cannot create the PrivX API client as there is an unknown configuration value for the PrivX API client ID. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the PRIVX_API_CLIENT_ID environment variable.",
		)
	}

	if data.APIClientSecret.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_client_secret"),
			"Unknown PrivX API client secret",
			"The provider cannot create the PrivX API client as there is an unknown configuration value for the PrivX API client secret. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the PRIVX_API_CLIENT_SECRET environment variable.",
		)
	}

	if data.OAuthClientID.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_oauth_client_id"),
			"Unknown PrivX API OAuth client ID",
			"The provider cannot create the PrivX API client as there is an unknown configuration value for the PrivX API OAuth client ID. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the PRIVX_API_OAUTH_ID environment variable.",
		)
	}

	if data.OAuthClientSecret.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_oauth_client_secret"),
			"Unknown PrivX OAuth client secret",
			"The provider cannot create the PrivX API client as there is an unknown configuration value for the PrivX API OAuth client secret. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the PRIVX_API_OAUTH_SECRET environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	apiBaseURL := os.Getenv("PRIVX_API_BASE_URL")
	apiBearerToken := os.Getenv("PRIVX_API_BEARER_TOKEN")
	apiClientID := os.Getenv("PRIVX_API_CLIENT_ID")
	apiClientSecret := os.Getenv("PRIVX_API_CLIENT_SECRET")
	oauthClientID := os.Getenv("PRIVX_API_OAUTH_CLIENT_ID")
	oauthClientSecret := os.Getenv("PRIVX_API_OAUTH_CLIENT_SECRET")

	if !data.APIBaseURL.IsNull() {
		apiBaseURL = data.APIBaseURL.ValueString()
	}

	if !data.APIBearerToken.IsNull() {
		apiBearerToken = data.APIBearerToken.ValueString()
	}

	if !data.APIClientID.IsNull() {
		apiClientID = data.APIClientID.ValueString()
	}

	if !data.APIClientSecret.IsNull() {
		apiClientSecret = data.APIClientSecret.ValueString()
	}

	if !data.OAuthClientID.IsNull() {
		oauthClientID = data.OAuthClientID.ValueString()
	}

	if !data.OAuthClientSecret.IsNull() {
		oauthClientSecret = data.OAuthClientSecret.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if apiBaseURL == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_base_url"),
			"Missing PrivX API base URL",
			"The provider cannot create the PrivX API client as there is a missing or empty value for the PrivX API base URL. "+
				"Set the host value in the configuration or use the PRIVX_API_BASE_URL environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if apiBearerToken == "" {

		if apiClientID == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("api_client_id"),
				"Missing PrivX API client ID",
				"The provider cannot create the PrivX API client as there is a missing or empty value for the PrivX API client ID. "+
					"Set the username value in the configuration or use the PRIVX_API_CLIENT_ID environment variable. "+
					"If either is already set, ensure the value is not empty.",
			)
		}

		if apiClientSecret == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("api_client_secret"),
				"Missing PrivX API client secret",
				"The provider cannot create the PrivX API client as there is a missing or empty value for the PrivX API client secret. "+
					"Set the password value in the configuration or use the PRIVX_API_CLIENT_SECRET environment variable. "+
					"If either is already set, ensure the value is not empty.",
			)
		}

		if oauthClientID == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("api_oauth_client_id"),
				"Missing PrivX OAuth client ID",
				"The provider cannot create the PrivX API client as there is a missing or empty value for the PrivX OAuth client ID. "+
					"Set the username value in the configuration or use the PRIVX_API_OAUTH_ID environment variable. "+
					"If either is already set, ensure the value is not empty.",
			)
		}

		if oauthClientSecret == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("api_oauth_client_secret"),
				"Missing PrivX API client secret",
				"The provider cannot create the PrivX API client as there is a missing or empty value for the PrivX OAuth client secret. "+
					"Set the password value in the configuration or use the PRIVX_API_OAUTH_SECRET environment variable. "+
					"If either is already set, ensure the value is not empty.",
			)
		}

	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "api_base_url", apiBaseURL)
	ctx = tflog.SetField(ctx, "api_bearer_token", apiBearerToken)
	ctx = tflog.SetField(ctx, "api_client_id", apiClientID)
	ctx = tflog.SetField(ctx, "api_client_secret", apiClientSecret)
	ctx = tflog.SetField(ctx, "api_oauth_client_id", oauthClientID)
	ctx = tflog.SetField(ctx, "api_oauth_client_secret", oauthClientSecret)
	//	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "api_client_secret")
	//	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "api_oauth_client_secret")
	//	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "api_bearer_token")

	tflog.Debug(ctx, "Creating PrivX client")

	connector, err := client.NewConnector(apiBaseURL, apiBearerToken, apiClientID, apiClientSecret, oauthClientID, oauthClientSecret)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create PrivX client",
			"An unexpected error occurred while attempting to create the provider client:\n"+
				err.Error(),
		)
		return
	}
	resp.DataSourceData = connector
	resp.ResourceData = connector

	tflog.Info(ctx, "Configured PrivX API client", map[string]any{"success": true})
}

func (p *privxProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewAccessGroupResource,
		NewExtenderResource,
		NewHostResource,
		NewRoleResource,
		NewSecretResource,
		NewSourceResource,
		NewAPIClientResource,
		NewCarrierResource,
	}
}

func (p *privxProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewAccessGroupDataSource,
		NewAPIClientDataSource,
		NewCarrierConfigDataSource,
		NewExtenderDataSource,
		NewExtenderConfigDataSource,
		NewWebproxyConfigDataSource,
		NewWebproxyDataSource,
		NewRoleDataSource,
		NewSecretDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &privxProvider{
			version: version,
		}
	}
}
