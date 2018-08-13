package openapi

import (
	"errors"
	"net/url"
)

// codebeat:disable[TOO_MANY_IVARS]

// License Object
type License struct {
	Name string
	URL  string
}

// Validate the values of License object.
func (license License) Validate() error {
	if license.Name == "" {
		return errors.New("license.name is required")
	}
	if license.URL != "" {
		_, err := url.ParseRequestURI(license.URL)
		return err
	}
	return nil
}
