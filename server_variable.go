package openapi

import "errors"

// codebeat:disable[TOO_MANY_IVARS]

// ServerVariable Object
type ServerVariable struct {
	Enum        []string
	Default     string
	Description string
}

// Validate the values of Server Variable object.
func (sv ServerVariable) Validate() error {
	if sv.Default == "" {
		return errors.New("server.default is required")
	}
	return nil
}
