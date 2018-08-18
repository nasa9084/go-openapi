package openapi

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
		return ErrRequired{Target: "securityScheme.type"}
	case "apiKey":
		return secScheme.validateFieldForAPIKey()
	case "http":
		return secScheme.validateFieldForHTTP()
	case "oauth2":
		return secScheme.validateFieldForOAuth2()
	case "openIdConnect":
		return secScheme.validateFieldForOpenIDConnect()
	}
	return ErrMustOneOf{Object: "securityScheme.type", ValidValues: []string{"apikey", "http", "oauth2", "openIdConnect"}}
}

func (secScheme SecurityScheme) validateFieldForAPIKey() error {
	if secScheme.Name == "" {
		return ErrRequired{"securityScheme.name"}
	}
	if secScheme.In == "" {
		return ErrRequired{"securityScheme.in"}
	}
	if secScheme.In != "query" && secScheme.In != "header" && secScheme.In != "cookie" {
		return ErrMustOneOf{Object: "securityScheme.in", ValidValues: []string{"query", "header", "cookie"}}
	}
	return nil
}

func (secScheme SecurityScheme) validateFieldForHTTP() error {
	if secScheme.Scheme == "" {
		return ErrRequired{Target: "securityScheme.scheme"}
	}
	return nil
}

func (secScheme SecurityScheme) validateFieldForOAuth2() error {
	if secScheme.Flows == nil {
		return ErrRequired{Target: "securityScheme.flows"}
	}
	return secScheme.Flows.Validate()
}

func (secScheme SecurityScheme) validateFieldForOpenIDConnect() error {
	return mustURL("securityScheme.openIdConnectUrl", secScheme.OpenIDConnectURL)
}
