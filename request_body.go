package openapi

import "errors"

// codebeat:disable[TOO_MANY_IVARS]

// RequestBody Object
type RequestBody struct {
	Description string
	Content     map[string]*MediaType
	Required    bool

	Ref string `yaml:"$ref"`
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
