package client

import (
	"fmt"

	"github.com/SSHcom/privx-sdk-go/oauth"
	"github.com/SSHcom/privx-sdk-go/restapi"
)

func authorize(apiBaseURL, bearerToken, apiClientID, apiClientSecret, oauthClientID, oauthClientSecret string) restapi.Authorizer {
	if bearerToken != "" {
		return oauth.WithToken("Bearer " + bearerToken)
	}
	return oauth.WithClientID(
		restapi.New(
			restapi.BaseURL(apiBaseURL),
		),
		oauth.Access(apiClientID),
		oauth.Secret(apiClientSecret),
		oauth.Digest(oauthClientID, oauthClientSecret),
	)
}

func NewConnector(apiBaseURL, bearerToken, apiClientID, apiClientSecret, oauthClientID, oauthClientSecret string) (*restapi.Connector, error) {
	auth := authorize(apiBaseURL, bearerToken, apiClientID, apiClientSecret, oauthClientID, oauthClientSecret)
	_, err := auth.AccessToken()
	if err != nil {
		return nil, fmt.Errorf("PrivX client authentication failed: %v", err)
	}
	connector := restapi.New(restapi.Auth(auth), restapi.Verbose(), restapi.BaseURL(apiBaseURL))
	return &connector, nil
}
