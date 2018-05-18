package openapi

import "errors"

// codebeat:disable[TOO_MANY_IVARS]

// Document represents a OpenAPI Specification document.
type Document struct {
	Version      string `yaml:"openapi"`
	Info         *Info
	Servers      []*Server
	Paths        Paths
	Components   *Components
	Security     *SecurityRequirement
	Tags         []*Tag
	ExternalDocs *ExternalDocumentation `yaml:"externalDocs"`
}

// Validate the values of spec.
func (doc Document) Validate() error {
	if err := validateOASVersion(doc.Version); err != nil {
		return err
	}
	if err := doc.validateRequiredObjects(); err != nil {
		return err
	}
	var validaters []validater
	validaters = append(validaters, doc.Info) // doc.Info nil check has done
	for _, s := range doc.Servers {
		validaters = append(validaters, s)
	}
	validaters = append(validaters, doc.Paths) // doc.Paths nil check has done
	if doc.Components != nil {
		validaters = append(validaters, doc.Components)
	}
	if doc.Security != nil {
		validaters = append(validaters, doc.Security)
	}
	for _, t := range doc.Tags {
		validaters = append(validaters, t)
	}
	if doc.ExternalDocs != nil {
		validaters = append(validaters, doc.ExternalDocs)
	}
	return validateAll(validaters)
}

func validateOASVersion(version string) error {
	if version == "" {
		return errors.New("openapi is required")
	}
	for _, v := range compatibleVersions {
		if version == v {
			return nil
		}
	}
	return errors.New("the OAS version is not supported")
}

func (doc Document) validateRequiredObjects() error {
	if doc.Info == nil {
		return errors.New("info is required")
	}
	if doc.Paths == nil {
		return errors.New("paths is required")
	}
	return nil
}
