package privx

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/SSHcom/privx-sdk-go/api/authorizer"
	"github.com/SSHcom/privx-sdk-go/api/rolestore"
	"github.com/SSHcom/privx-sdk-go/api/userstore"
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
	debug               bool
}

type privx_API_client_connector struct {
	Connector restapi.Connector
	Config    *Config
	Token     string
}

func authorize(cfg *Config) restapi.Authorizer {
	return oauth.WithClientID(
		restapi.New(
			restapi.BaseURL(cfg.base_url),
			restapi.Verbose(), /** DEBUG **/
		),
		oauth.Access(cfg.api_client_id),
		oauth.Secret(cfg.api_client_secret),
		oauth.Digest(cfg.oauth_client_id, cfg.oauth_client_secret),
	)
}

// 2. Create HTTP transport for PrivX API
func createConnector(cfg *Config) (restapi.Connector, string, []string) {
	var debug []string
	auth := authorize(cfg)
	AccessToken, err := auth.AccessToken()
	if err != nil {
		debug = append(debug, "API Acces token generated :"+AccessToken)
		debug = append(debug, "Debug mode : \n"+strconv.FormatBool(cfg.debug))
		if err != nil {
			debug = append(debug, "Error message from token generation : "+err.Error())
		}
		debug = append(debug, "Privx baseURL: "+cfg.base_url)
		debug = append(debug, "Privx api_client_id: "+cfg.api_client_id)
		debug = append(debug, "Privx oauth_client_id"+cfg.oauth_client_id)
	}
	return restapi.New(restapi.Auth(*&auth), restapi.BaseURL(cfg.base_url)), AccessToken, debug
}

// NewClient func...
func (cfg *Config) NewClientConnector(ctx context.Context) (interface{}, diag.Diagnostics) {
	connector, token, debug := createConnector(cfg)
	api_client_connector := privx_API_client_connector{
		Connector: connector,
		Config:    cfg,
		Token:     token, /** DEBUG **/
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

func createAuthorizerClient(ctx context.Context, connector restapi.Connector) *authorizer.Client {
	return authorizer.New(connector)
}

func createUserStoreClient(ctx context.Context, connector restapi.Connector) *userstore.UserStore {
	return userstore.New(connector)
}

func createRoleClient(ctx context.Context, connector restapi.Connector) *rolestore.RoleStore {
	return rolestore.New(connector)
}
