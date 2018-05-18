package openapi

import "errors"

// codebeat:disable[TOO_MANY_IVARS]

// SecurityScheme Object
type SecurityScheme struct {
	Type             string
	Description      string
	Name             string
	In               string
	Scheme           string
	BearerFormat     string `yaml:"bearerFormat"`
	Flows            *OAuthFlows
	OpenIDConnectURL string `yaml:"openIdConnectUrl"`

	Ref string `yaml:"$ref"`
}

// Validate the values of SecurityScheme object.
func (secScheme SecurityScheme) Validate() error {
	switch secScheme.Type {
	case "":
		return errors.New("securityScheme.type is required")
	case "apiKey":
		return secScheme.validateFieldForAPIKey()
	case "http":
		return secScheme.validateFieldForHTTP()
	case "oauth2":
		return secScheme.validateFieldForOAuth2()
	case "openIdConnect":
		return secScheme.validateFieldForOpenIDConnect()
	}
	return errors.New("securityScheme.type must be one of [apikey, http, oauth2, openIdConnect]")
}

func (secScheme SecurityScheme) validateFieldForAPIKey() error {
	if secScheme.Name == "" {
		return errors.New("securityScheme.name is required")
	}
	if secScheme.In == "" {
		return errors.New("securityScheme.in is required")
	}
	if secScheme.In != "query" && secScheme.In != "header" && secScheme.In != "cookie" {
		return errors.New("securityScheme.in must be one of [query, header, cookie]")
	}
	return nil
}

func (secScheme SecurityScheme) validateFieldForHTTP() error {
	if secScheme.Scheme == "" {
		return errors.New("securityScheme.scheme is required")
	}
	return nil
}

func (secScheme SecurityScheme) validateFieldForOAuth2() error {
	if secScheme.Flows == nil {
		return errors.New("securityScheme.flows is required")
	}
	return secScheme.Flows.Validate()
}

func (secScheme SecurityScheme) validateFieldForOpenIDConnect() error {
	return mustURL("securityScheme.openIdConnectUrl is required", secScheme.OpenIDConnectURL)
}
