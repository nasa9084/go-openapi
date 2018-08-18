package openapi

import (
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
	Security     []*SecurityRequirement
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
	for _, securityRequirement := range doc.Security {
		validaters = append(validaters, securityRequirement)
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
		return ErrRequired{Target: "openapi"}
	}
	splited := strings.Split(version, ".")
	if len(splited) != 3 {
		return ErrFormatInvalid{Target: "openapi version", Format: "X.Y.Z"}
	}
	major, err := strconv.Atoi(splited[0])
	if err != nil {
		return ErrFormatInvalid{Target: "major part of openapi version"}
	}
	minor, err := strconv.Atoi(splited[1])
	if err != nil {
		return ErrFormatInvalid{Target: "minor part of openapi version"}
	}
	_, err = strconv.Atoi(splited[2])
	if err != nil {
		return ErrFormatInvalid{Target: "patch part of openapi version"}
	}
	if major == 3 && 0 <= minor {
		return nil
	}
	return UnsupportedVersionError
}

func (doc Document) validateRequiredObjects() error {
	if doc.Info == nil {
		return ErrRequired{Target: "info"}
	}
	if doc.Paths == nil {
		return ErrRequired{Target: "paths"}
	}
	return nil
}
