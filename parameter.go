package openapi

import "errors"

// codebeat:disable[TOO_MANY_IVARS]

// Parameter Object
type Parameter struct {
	Name            string
	In              string
	Description     string
	Required        bool
	Deprecated      string
	AllowEmptyValue bool `yaml:"allowEmptyValue"`

	Style         string
	Explode       bool
	AllowReserved bool `yaml:"allowReserved"`
	Schema        *Schema
	Example       interface{}
	Examples      map[string]*Example

	Content map[string]*MediaType

	Ref string `yaml:"$ref"`
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
	switch parameter.In {
	case "query", "header", "path", "cookie":
	default:
		return errors.New("parameter.in should be one of: query, header, path, cookie")
	}
	if parameter.In == "path" && !parameter.Required {
		return errors.New("if parameter.in is path, required must be true")
	}
	if parameter.In != "query" && parameter.AllowEmptyValue {
		return errors.New("allowEmptyValue is valid only for query parameters")
	}
	validaters := []validater{}
	if parameter.Schema != nil {
		validaters = append(validaters, parameter.Schema)
	}
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
