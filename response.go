package openapi

import "errors"

// codebeat:disable[TOO_MANY_IVARS]

// Response Object
type Response struct {
	Description string
	Headers     map[string]*Header
	Content     map[string]*MediaType
	Links       map[string]*Link

	Ref string `yaml:"$ref"`
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
