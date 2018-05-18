package openapi

import (
	"errors"
	"net/url"
)

// codebeat:disable[TOO_MANY_IVARS]

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
