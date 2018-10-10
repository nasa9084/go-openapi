package openapi

import (
	"net/url"
)

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
		return ErrRequired{"info.title"}
	}
	if info.TermsOfService != "" {
		if _, err := url.ParseRequestURI(info.TermsOfService); err != nil {
			return ErrFormatInvalid{Target: "info.termsOfService", Format: "URL"}
		}
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
		return ErrRequired{Target: "info.version"}
	}
	return nil
}
