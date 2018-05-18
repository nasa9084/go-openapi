package openapi

import "errors"

// codebeat:disable[TOO_MANY_IVARS]

// Tag Object
type Tag struct {
	Name         string
	Description  string
	ExternalDocs *ExternalDocumentation `yaml:"externalDocs"`
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
