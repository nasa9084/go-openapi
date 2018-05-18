package openapi

import "github.com/nasa9084/go-openapi/oauth"

// codebeat:disable[TOO_MANY_IVARS]

// OAuthFlows Object
type OAuthFlows struct {
	Implicit          *OAuthFlow
	Password          *OAuthFlow
	ClientCredentials *OAuthFlow `yaml:"clientCredentials"`
	AuthorizationCode *OAuthFlow `yaml:"authorizationCode"`
}

// Validate the values of OAuthFlows Object.
func (oauthFlows OAuthFlows) Validate() error {
	if oauthFlows.Implicit != nil {
		if err := oauthFlows.Implicit.Validate(oauth.ImplicitFlow); err != nil {
			return err
		}
	}
	if oauthFlows.Password != nil {
		if err := oauthFlows.Password.Validate(oauth.PasswordFlow); err != nil {
			return err
		}
	}
	if oauthFlows.ClientCredentials != nil {
		if err := oauthFlows.ClientCredentials.Validate(oauth.ClientCredentialsFlow); err != nil {
			return err
		}
	}
	if oauthFlows.AuthorizationCode != nil {
		if err := oauthFlows.AuthorizationCode.Validate(oauth.AuthorizationCodeFlow); err != nil {
			return err
		}
	}
	return nil
}
