package openapi

import (
	"errors"
	"net/url"
	"regexp"
)

// codebeat:disable[TOO_MANY_IVARS]

var tmplVarRegexp = regexp.MustCompile("{[^}]+}")

// Server Object
type Server struct {
	URL         string
	Description string
	Variables   map[string]*ServerVariable
}

// Validate the values of Server object.
func (server Server) Validate() error {
	if server.URL == "" {
		return errors.New("server.url is required")
	}
	// replace template variable with placeholder to validate the replaced string
	// is valid URL or not
	serverURL := tmplVarRegexp.ReplaceAllLiteralString(server.URL, "ph")
	// use url.Parse because relative URL is allowed
	if _, err := url.Parse(serverURL); err != nil {
		return err
	}
	validaters := []validater{}
	for _, sv := range server.Variables {
		validaters = append(validaters, sv)
	}
	return validateAll(validaters)
}
