package privx

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/SSHcom/privx-sdk-go/api/authorizer"
	"github.com/SSHcom/privx-sdk-go/api/hoststore"
	"github.com/SSHcom/privx-sdk-go/api/rolestore"
	"github.com/SSHcom/privx-sdk-go/api/userstore"
	"github.com/SSHcom/privx-sdk-go/api/vault"
	"github.com/SSHcom/privx-sdk-go/oauth"
	"github.com/SSHcom/privx-sdk-go/restapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

// Config struct ...
type Config struct {
	base_url            string
	oauth_client_id     string
	oauth_client_secret string
	api_client_id       string
	api_client_secret   string
	token               string
	debug               bool
}

type privx_API_client_connector struct {
	Connector restapi.Connector
	Config    *Config
	Token     string
}

func authorize(cfg *Config) restapi.Authorizer {
	if cfg.token != "" {
		return oauth.WithToken("Bearer " + cfg.token)
	}
	return oauth.WithClientID(
		restapi.New(
			restapi.BaseURL(cfg.base_url),
		),
		oauth.Access(cfg.api_client_id),
		oauth.Secret(cfg.api_client_secret),
		oauth.Digest(cfg.oauth_client_id, cfg.oauth_client_secret),
	)
}

// 2. Create HTTP transport for PrivX API
func createConnector(cfg *Config) (restapi.Connector, []string) {
	var debug []string
	auth := authorize(cfg)
	AccessToken, err := auth.AccessToken()
	if err != nil {
		debug = append(debug, "API Acces token generated :"+AccessToken)
		debug = append(debug, "Debug mode : \n"+strconv.FormatBool(cfg.debug))
		debug = append(debug, "Error message from oauth : "+err.Error())
		debug = append(debug, "Privx baseURL: "+cfg.base_url)
		debug = append(debug, "Privx api_client_id: "+cfg.api_client_id)
		debug = append(debug, "Privx oauth_client_id"+cfg.oauth_client_id)
	}
	connector := restapi.New(restapi.Auth(auth), restapi.Verbose(), restapi.BaseURL(cfg.base_url))
	return connector, debug
}

// NewClient func...
func (cfg *Config) NewClientConnector(ctx context.Context) (interface{}, diag.Diagnostics) {
	connector, debug := createConnector(cfg)
	api_client_connector := privx_API_client_connector{
		Connector: connector,
		Config:    cfg,
	}
	if cfg.debug || len(debug) != 0 {
		return api_client_connector, diag.FromErr(errors.New(strings.Join(debug, "\n")))
	}
	return api_client_connector, nil
}

// Transform slice interface into map of strings
func flattenSimpleSlice(slice []interface{}) []string {
	results := make([]string, 0, len(slice))
	for _, obj := range slice {
		str := fmt.Sprintf("%v", obj)
		results = append(results, str)
	}
	return results
}

/*
--------------------------------------------
Functions to instanciate privx object/endpoint API Client.
*/
func createAuthorizerClient(ctx context.Context, connector restapi.Connector) *authorizer.Client {
	return authorizer.New(connector)
}

func createHostClient(ctx context.Context, connector restapi.Connector) *hoststore.HostStore {
	return hoststore.New(connector)
}

func createRoleClient(ctx context.Context, connector restapi.Connector) *rolestore.RoleStore {
	return rolestore.New(connector)
}

func createUserStoreClient(ctx context.Context, connector restapi.Connector) *userstore.UserStore {
	return userstore.New(connector)
}

func createVaultClient(ctx context.Context, connector restapi.Connector) *vault.Vault {
	return vault.New(connector)
}

/*
--------------------------------------------
*/
