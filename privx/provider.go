package privx

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider returns the provider to be use by the code.
func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"privx_api_client_id": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"PRIVX_API_CLIENT_ID",
				}, ""),
				Description: "PrivX API Oauth client ID",
			},
			"privx_api_client_secret": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"PRIVX_API_CLIENT_SECRET",
				}, ""),
				Description: "PrivX API Oauth client ID",
				Sensitive:   true,
			},
			"privx_oauth_client_id": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"PRIVX_OAUTH_CLIENT_ID",
				}, ""),
				Description: "PrivX API Oauth client ID",
			},
			"privx_oauth_client_secret": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"PRIVX_OAUTH_CLIENT_SECRET",
				}, ""),
				Description: "Privx API Oauth Client Secret",
				Sensitive:   true,
			},
			"privx_api_bearer_token": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"PRIVX_API_BEARER_TOKEN",
				}, ""),
				Description: "PrivX bearer token",
				Sensitive:   true,
			},
			"privx_api_base_url": {
				Type:     schema.TypeString,
				Required: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"PRIVX_API_BASE_URL",
				}, ""),
				Description: "PrivX API Base URL",
			},
			"privx_debug": {
				Type:     schema.TypeBool,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"PRIVX_DEBUG",
				}, ""),
				Description: "PrivX debug mode",
			},
		},
		DataSourcesMap:       getDataSourcesMap(),
		ResourcesMap:         getResourcesMap(),
		ConfigureContextFunc: providerConfigure,
	}
	return provider
}

func getDataSourcesMap() map[string]*schema.Resource {
	dataSourcesMap := map[string]*schema.Resource{
		"privx_extenders":       dataSourcePrivxExtender(),
		"privx_roles":           dataSourcePrivxRoles(),
		"privx_access_groups":   dataSourcePrivxAccessGroups(),
		"privx_extender_config": dataSourcePrivxExtenderConfig(),
	}
	return dataSourcesMap
}

func getResourcesMap() map[string]*schema.Resource {
	resourcesMap := map[string]*schema.Resource{
		"privx_access_group":      resourcePrivXAccessGroup(),
		"privx_api_client":        resourcePrivXApiClient(),
		"privx_extender":          resourcePrivXExtender(),
		"privx_host":              resourcePrivXHost(),
		"privx_secret_credential": resourcePrivXSecretCredential(),
		"privx_secret_json":       resourcePrivxJsonSecret(),
		"privx_role":              resourcePrivXRole(),
		"privx_source":            resourcePrivXSource(),
	}
	return resourcesMap
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	config := Config{
		oauth_client_id:     d.Get("privx_oauth_client_id").(string),
		oauth_client_secret: d.Get("privx_oauth_client_secret").(string),
		api_client_id:       d.Get("privx_api_client_id").(string),
		api_client_secret:   d.Get("privx_api_client_secret").(string),
		base_url:            d.Get("privx_api_base_url").(string),
		token:               d.Get("privx_api_bearer_token").(string),
		debug:               d.Get("privx_debug").(bool),
	}
	api_client_connector, err := config.NewClientConnector(ctx)
	if err != nil {
		return nil, err
	}
	return api_client_connector, nil
}
