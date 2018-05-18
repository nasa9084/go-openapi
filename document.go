package openapi

import (
	"errors"
	"strconv"
	"strings"
)

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
	splited := strings.Split(version, ".")
	if len(splited) != 3 {
		return errors.New("openapi version is not valid: version format should be X.Y.Z")
	}
	major, err := strconv.Atoi(splited[0])
	if err != nil {
		return errors.New("major part of openapi version is invalid format")
	}
	minor, err := strconv.Atoi(splited[1])
	if err != nil {
		return errors.New("minor part of openapi version is invalid format")
	}
	_, err := strconv.Atoi(splited[2])
	if err != nil {
		return errors.New("patch part of openapi version is invalid format")
	}
	if major == 3 && 0 <= minor {
		return nil
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
