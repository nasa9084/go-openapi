package openapi

import (
	"errors"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var (
	compatibleVersions = []string{"3.0.0", "3.0.1"}
	emailRegexp        = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	mapKeyRegexp       = regexp.MustCompile("^[a-zA-Z0-9\\.\\-_]+$")
)

type validater interface {
	Validate() error
}

func validateAll(vs []validater) error {
	for _, v := range vs {
		if v == nil {
			continue
		}
		if err := v.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// Validate the values of spec.
func (doc Document) Validate() error {
	if err := validateOASVersion(doc.Version); err != nil {
		return err
	}
	if err := doc.validateRequiredObjects(); err != nil {
		return err
	}
	var validaters []validater
	validaters = append(validaters, doc.Info)
	for _, s := range doc.Servers {
		validaters = append(validaters, s)
	}
	validaters = append(validaters, doc.Paths)
	validaters = append(validaters, doc.Components)
	validaters = append(validaters, doc.Security)
	for _, t := range doc.Tags {
		validaters = append(validaters, t)
	}
	validaters = append(validaters, doc.ExternalDocs)
	return validateAll(validaters)
}

func validateOASVersion(version string) error {
	if version == "" {
		return errors.New("openapi is required")
	}
	for _, v := range compatibleVersions {
		if version == v {
			return nil
		}
	}
	return errors.New("the OAS version is not supported")
}

func (doc Document) validateRequiredObjects() error {
	if doc.Info == nil {
		return errors.New("info is required")
	}
	if doc.Paths == nil {
		return errors.New("paths is required")
	}
	return nil
}

// Validate the values of Info object.
func (info Info) Validate() error {
	if info.Title == "" {
		return errors.New("info.title is required")
	}
	if err := mustURL("info.termsOfService", info.TermsOfService); err != nil {
		return err
	}
	validaters := []validater{info.Contact, info.License}
	if err := validateAll(validaters); err != nil {
		return err
	}
	if info.Version == "" {
		return errors.New("info.version is required")
	}
	return nil
}

// Validate the values of Contact object.
func (contact Contact) Validate() error {
	if err := mustURL("contact.url", contact.URL); err != nil {
		return err
	}
	if contact.Email != "" {
		if !emailRegexp.MatchString(contact.Email) {
			return errors.New("email format invalid")
		}
	}
	return nil
}

// Validate the values of License object.
func (license License) Validate() error {
	if license.Name == "" {
		return errors.New("license.name is required")
	}
	return mustURL("license.url", license.URL)
}

// Validate the values of Server object.
func (server Server) Validate() error {
	if server.URL == "" {
		return errors.New("server.url is required")
	}
	// use url.Parse because relative URL is allowed
	if _, err := url.Parse(server.URL); err != nil {
		return err
	}
	validaters := []validater{}
	for _, sv := range server.Variables {
		validaters = append(validaters, sv)
	}
	return validateAll(validaters)
}

// Validate the values of Server Variable object.
func (sv ServerVariable) Validate() error {
	if sv.Default == "" {
		return errors.New("server.default is required")
	}
	return nil
}

// Validate the values of Components object.
func (components Components) Validate() error {
	if err := components.validateKeys(); err != nil {
		return err
	}
	validaters := reduceComponentObjects(components)
	return validateAll(validaters)
}

func (components Components) validateKeys() error {
	keys := reduceComponentKeys(components)
	for _, k := range keys {
		if !mapKeyRegexp.MatchString(k) {
			return errors.New("map key format is invalid")
		}
	}
	return nil
}

func reduceComponentKeys(components Components) []string {
	keys := []string{}
	for k := range components.Schemas {
		keys = append(keys, k)
	}
	for k := range components.Responses {
		keys = append(keys, k)
	}
	for k := range components.Parameters {
		keys = append(keys, k)
	}
	for k := range components.Examples {
		keys = append(keys, k)
	}
	for k := range components.RequestBodies {
		keys = append(keys, k)
	}
	for k := range components.Headers {
		keys = append(keys, k)
	}
	for k := range components.SecuritySchemes {
		keys = append(keys, k)
	}
	for k := range components.Links {
		keys = append(keys, k)
	}
	for k := range components.Callbacks {
		keys = append(keys, k)
	}
	return keys
}

func reduceComponentObjects(components Components) []validater {
	validaters := []validater{}
	for _, schema := range components.Schemas {
		validaters = append(validaters, schema)
	}
	for _, response := range components.Responses {
		validaters = append(validaters, response)
	}
	for _, parameter := range components.Parameters {
		validaters = append(validaters, parameter)
	}

	// example has no validation

	for _, reqBody := range components.RequestBodies {
		validaters = append(validaters, reqBody)
	}
	for _, header := range components.Headers {
		validaters = append(validaters, header)
	}
	for _, secScheme := range components.SecuritySchemes {
		validaters = append(validaters, secScheme)
	}
	for _, link := range components.Links {
		validaters = append(validaters, link)
	}
	for _, callback := range components.Callbacks {
		validaters = append(validaters, callback)
	}
	return validaters
}

// Validate the values of Paths object.
func (paths Paths) Validate() error {
	for path, pathItem := range paths {
		if !strings.HasPrefix(path, "/") {
			return errors.New("path name must begin with a slash")
		}
		if err := pathItem.Validate(); err != nil {
			return err
		}
	}
	if hasDuplicatedOperationID(paths) {
		return errors.New("operation id is duplicated")
	}
	return nil
}

func hasDuplicatedOperationID(paths Paths) bool {
	// TODO
	return false
}

// Validate the values of PathItem object.
func (pathItem PathItem) Validate() error {
	validaters := []validater{
		pathItem.Get,
		pathItem.Put,
		pathItem.Post,
		pathItem.Delete,
		pathItem.Options,
		pathItem.Head,
		pathItem.Patch,
		pathItem.Trace,
	}
	for _, s := range pathItem.Servers {
		validaters = append(validaters, s)
	}
	if hasDuplicatedParameter(pathItem.Parameters) {
		return errors.New("some parameter is duplicated")
	}
	for _, p := range pathItem.Parameters {
		validaters = append(validaters, p)
	}
	return validateAll(validaters)
}

func hasDuplicatedParameter(parameters []*Parameter) bool {
	for i, p := range parameters {
		for _, q := range parameters[i+1:] {
			if p.Name == q.Name && p.In == q.In {
				return true
			}
		}
	}
	return false
}

// Validate the values of Operation object.
func (operation Operation) Validate() error {
	validaters := []validater{}
	validaters = append(validaters, operation.ExternalDocs)
	if hasDuplicatedParameter(operation.Parameters) {
		return errors.New("some parameter is duplicated")
	}
	validaters = append(validaters, operation.RequestBody)
	if operation.Responses == nil {
		return errors.New("operation.responses is required")
	}
	validaters = append(validaters, operation.Responses)
	for _, callback := range operation.Callbacks {
		validaters = append(validaters, callback)
	}
	for _, security := range operation.Security {
		validaters = append(validaters, security)
	}
	for _, server := range operation.Servers {
		validaters = append(validaters, server)
	}
	return validateAll(validaters)
}

// Validate the values of ExternalDocumentaion object.
func (externalDocumentation ExternalDocumentation) Validate() error {
	return mustURL("externalDocumentation.url", externalDocumentation.URL)
}

// Validate the values of Parameter object.
// This function DOES NOT check whether the name field correspond to the associated path or not.
func (parameter Parameter) Validate() error {
	if parameter.Name == "" {
		return errors.New("parameter.name is required")
	}
	if parameter.In == "" {
		return errors.New("parameter.in is required")
	}
	if parameter.In == "path" && !parameter.Required {
		return errors.New("if parameter.in is path, required must be true")
	}
	validaters := []validater{parameter.Schema}
	if v, ok := parameter.Example.(validater); ok {
		validaters = append(validaters, v)
	}

	// example has no validation

	if len(parameter.Content) > 1 {
		return errors.New("parameter.content must only contain one entry")
	}
	for _, mediaType := range parameter.Content {
		validaters = append(validaters, mediaType)
	}
	return validateAll(validaters)
}

// Validate the values of RequestBody object.
func (requestBody RequestBody) Validate() error {
	if requestBody.Content == nil || len(requestBody.Content) == 0 {
		return errors.New("requestBody.content is required")
	}
	for _, mediaType := range requestBody.Content {
		if err := mediaType.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// Validate the values of MediaType object.
// This function DOES NOT check whether the encoding object is in schema or not.
func (mediaType MediaType) Validate() error {
	validaters := []validater{mediaType.Schema}
	if v, ok := mediaType.Example.(validater); ok {
		validaters = append(validaters, v)
	}

	// example has no validation

	for _, e := range mediaType.Encoding {
		validaters = append(validaters, e)
	}
	return validateAll(validaters)
}

// Validate the values of Encoding object.
func (encoding Encoding) Validate() error {
	for _, header := range encoding.Headers {
		if err := header.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// Validate the values of Responses object.
func (responses Responses) Validate() error {
	for status, response := range responses {
		if err := validateStatusCode(status); err != nil {
			return err
		}
		if err := response.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func validateStatusCode(statusStr string) error {
	switch statusStr {
	case "default", "1XX", "2XX", "3XX", "4XX", "5XX":
		return nil
	}
	statusInt, err := strconv.Atoi(statusStr)
	if err != nil {
		return err
	}
	if statusInt < 100 || 599 < statusInt {
		return errors.New("status code is invalid")
	}
	return nil
}

// Validate the value of Response object.
func (response Response) Validate() error {
	if response.Description == "" {
		return errors.New("response.description is required")
	}
	validaters := []validater{}
	for _, header := range response.Headers {
		validaters = append(validaters, header)
	}
	for _, mediaType := range response.Content {
		validaters = append(validaters, mediaType)
	}
	for _, link := range response.Links {
		validaters = append(validaters, link)
	}
	return validateAll(validaters)
}

// Validate the values of Callback object.
func (callback Callback) Validate() error {
	for _, pathItem := range callback {
		if err := pathItem.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// Validate the values of Link object.
func (link Link) Validate() error {
	validaters := []validater{}
	for _, i := range link.Parameters {
		if v, ok := i.(validater); ok {
			validaters = append(validaters, v)
		}
	}
	if v, ok := link.RequestBody.(validater); ok {
		validaters = append(validaters, v)
	}
	validaters = append(validaters, link.Server)
	return validateAll(validaters)
}

// Validate the values of Header object.
func (header Header) Validate() error {
	validaters := []validater{header.Schema}
	if v, ok := header.Example.(validater); ok {
		validaters = append(validaters, v)
	}

	// example has no validation

	if len(header.Content) > 1 {
		return errors.New("header.content must only contain one entry")
	}
	for _, mediaType := range header.Content {
		validaters = append(validaters, mediaType)
	}
	return validateAll(validaters)
}

// Validate the values of Tag object.
func (tag Tag) Validate() error {
	if tag.Name == "" {
		return errors.New("tag.name is required")
	}
	if tag.ExternalDocs != nil {
		return tag.ExternalDocs.Validate()
	}
	return nil
}

// Validate the values of Schema object.
func (schema Schema) Validate() error {
	validaters := []validater{
		schema.AllOf,
		schema.OneOf,
		schema.AnyOf,
		schema.Not,
		schema.Items,
		schema.Discriminator,
		schema.XML,
		schema.ExternalDocs,
	}
	for _, property := range schema.Properties {
		validaters = append(validaters, property)
	}
	if e, ok := schema.Example.(validater); ok {
		validaters = append(validaters, e)
	}
	return validateAll(validaters)
}

// Validate the values of Descriminator object.
func (discriminator Discriminator) Validate() error {
	if discriminator.PropertyName == "" {
		return errors.New("discriminator.propertyName is required")
	}
	return nil
}

// Validate the values of XML object.
func (xml XML) Validate() error {
	return mustURL("xml.namespace", xml.Namespace)
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

// Validate the values of OAuthFlows Object.
func (oauthFlows OAuthFlows) Validate() error {
	if oauthFlows.Implicit != nil {
		if err := oauthFlows.Implicit.Validate("implicit"); err != nil {
			return err
		}
	}
	if oauthFlows.Password != nil {
		if err := oauthFlows.Password.Validate("password"); err != nil {
			return err
		}
	}
	if oauthFlows.ClientCredentials != nil {
		if err := oauthFlows.ClientCredentials.Validate("clientCredentials"); err != nil {
			return err
		}
	}
	if oauthFlows.AuthorizationCode != nil {
		if err := oauthFlows.AuthorizationCode.Validate("authorizationCode"); err != nil {
			return err
		}
	}
	return nil
}

// Validate the values of OAuthFlow object.
func (oauthFlow OAuthFlow) Validate(typ string) error {
	if typ == "implicit" || typ == "authorizationCode" {
		if err := mustURL("oauthFlow.authorizationUrl", oauthFlow.AuthorizationURL); err != nil {
			return err
		}
	}
	if typ == "password" || typ == "clientCredentials" || typ == "authorizationCode" {
		if err := mustURL("oauthFlow.tokenUrl", oauthFlow.TokenURL); err != nil {
			return err
		}
	}
	if err := mustURL("oauthFlow.refreshUrl", oauthFlow.RefreshURL); err != nil {
		return err
	}
	if oauthFlow.Scopes == nil || len(oauthFlow.Scopes) == 0 {
		return errors.New("oauthFlow.scopes is required")
	}

	return nil
}

func mustURL(name, urlStr string) error {
	if urlStr == "" {
		return errors.New(name + " is required")
	}
	_, err := url.ParseRequestURI(urlStr)
	return err
}

// Validate the values of SecurityRequirement object.
// TODO: implement varidation
func (secReq SecurityRequirement) Validate() error {
	return nil
}
