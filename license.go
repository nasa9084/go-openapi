package openapi

import "errors"

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
	return mustURL("license.url", license.URL)
}
