package openapi

import "errors"

// codebeat:disable[TOO_MANY_IVARS]

// Info Object
type Info struct {
	Title          string
	Description    string
	TermsOfService string `yaml:"termsOfService"`
	Contact        *Contact
	License        *License
	Version        string
}

// Validate the values of Info object.
func (info Info) Validate() error {
	if info.Title == "" {
		return errors.New("info.title is required")
	}
	if err := mustURL("info.termsOfService", info.TermsOfService); err != nil {
		return err
	}
	validaters := []validater{}
	if info.Contact != nil {
		validaters = append(validaters, info.Contact)
	}
	if info.License != nil {
		validaters = append(validaters, info.License)
	}
	if err := validateAll(validaters); err != nil {
		return err
	}
	if info.Version == "" {
		return errors.New("info.version is required")
	}
	return nil
}
