package openapi

import (
	"errors"
	"net/url"

	"github.com/nasa9084/go-openapi/oauth"
)

// codebeat:disable[TOO_MANY_IVARS]

// OAuthFlow Object
type OAuthFlow struct {
	AuthorizationURL string `yaml:"authorizationUrl"`
	TokenURL         string `yaml:"tokenUrl"`
	RefreshURL       string `yaml:"refreshUrl"`
	Scopes           map[string]string
}

var defined = struct{}{}

var validFlowTypes = map[string]struct{}{
	oauth.ImplicitFlow:          defined,
	oauth.PasswordFlow:          defined,
	oauth.ClientCredentialsFlow: defined,
	oauth.AuthorizationCodeFlow: defined,
}

var requireAuthorizationURL = map[string]struct{}{
	oauth.ImplicitFlow:          defined,
	oauth.AuthorizationCodeFlow: defined,
}

var requireTokenURL = map[string]struct{}{
	oauth.PasswordFlow:          defined,
	oauth.ClientCredentialsFlow: defined,
	oauth.AuthorizationCodeFlow: defined,
}

// Validate the values of OAuthFlow object.
func (oauthFlow OAuthFlow) Validate(typ string) error {
	if _, ok := validFlowTypes[typ]; !ok {
		return errors.New("invalid type name")
	}
	if _, ok := requireAuthorizationURL[typ]; ok {
		if err := mustURL("oauthFlow.authorizationUrl", oauthFlow.AuthorizationURL); err != nil {
			return err
		}
	}
	if _, ok := requireTokenURL[typ]; ok {
		if err := mustURL("oauthFlow.tokenUrl", oauthFlow.TokenURL); err != nil {
			return err
		}
	}
	if oauthFlow.RefreshURL != "" {
		if _, err := url.ParseRequestURI(oauthFlow.RefreshURL); err != nil {
			return err
		}
	}
	if oauthFlow.Scopes == nil || len(oauthFlow.Scopes) == 0 {
		return errors.New("oauthFlow.scopes is required")
	}

	return nil
}
