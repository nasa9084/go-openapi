package openapi

import (
	"fmt"
	"strings"
)

// ErrFormatInvalid is returned some error caused by string format is occurred.
type ErrFormatInvalid struct {
	Target string
	Format string
}

func (fe ErrFormatInvalid) Error() string {
	if fe.Format == "" {
		return fmt.Sprintf("%s format is invalid", fe.Target)
	}
	return fmt.Sprintf("%s format is invalid: should be %s", fe.Target, fe.Format)
}

// central error variables relating format
var (
	EmailFormatError  = ErrFormatInvalid{Target: "email"}
	MapKeyFormatError = ErrFormatInvalid{Target: "map key"}
	PathFormatError   = ErrFormatInvalid{Target: "path"}
)

// ErrRequired is returned when missing some required parameter
type ErrRequired struct {
	Target string
}

func (re ErrRequired) Error() string {
	return fmt.Sprintf("%s is required", re.Target)
}

type errString string

func (ge errString) Error() string {
	return string(ge)
}

const (
	// UnsupportedVersionError is returned when the openapi version
	// is unsupported by this package.
	UnsupportedVersionError errString = "the OAS version is not supported"
	// InvalidFlowTypeError is returned when the OAuth flow type is invalid
	// or not set to the object.
	InvalidFlowTypeError errString = "invalid flow type"
	// RequiredMustTrueError is returned when the value of parameter.required is
	// false when parameter.in is path.
	RequiredMustTrueError errString = "required must be true if parameter.in is path"
	// AllowEmptyValueNotValidError is returned when allowEmptyValue is specified
	// but parameter.in is not query.
	AllowEmptyValueNotValidError errString = "allowEmptyValue is valid only for query parameters"
	// InvalidStatusCodeError is returned when specified status code is not
	// valid as HTTP status code.
	InvalidStatusCodeError errString = "status code is invalid"
	// MissingRootDocumentError is returned when validating securityRequirement
	// object but root document is not set.
	MissingRootDocumentError errString = "missing root document for security requirement"
)

type errTooManyContentEntry struct {
	target string
}

func (etme errTooManyContentEntry) Error() string {
	return fmt.Sprintf("%s.content must only contain one entry", etme.target)
}

var (
	// ErrTooManyHeaderContent is returned when the length of header.content
	// is more than 2.
	ErrTooManyHeaderContent = errTooManyContentEntry{target: "header"}
	// ErrTooManyParameterContent is returned when the length of parameter.content
	// is more than 2.
	ErrTooManyParameterContent = errTooManyContentEntry{target: "parameter"}
)

type errDuplicated struct {
	target string
}

func (de errDuplicated) Error() string {
	return fmt.Sprintf("some %s are duplicated", de.target)
}

var (
	// ErrOperationIDDuplicated is returned when some operation ids are
	// duplicated but operation ids cannot be duplicated.
	ErrOperationIDDuplicated = errDuplicated{target: "operation ids"}
	// ErrParameterDuplicated is returned when some parameters are duplicated
	// but cannot be duplicated.
	ErrParameterDuplicated = errDuplicated{target: "parameters"}
)

// ErrNotDeclared is returned when the securityScheme name is
// not defined in components object in the document.
type ErrNotDeclared struct {
	Name string
}

func (snde ErrNotDeclared) Error() string {
	return fmt.Sprintf("%s is not declared in components.securitySchemes", snde.Name)
}

// ErrMustEmpty returned when the securityRequirement is not
// empty but must be empty.
type ErrMustEmpty struct {
	Type string
}

func (srmee ErrMustEmpty) Error() string {
	return fmt.Sprintf("securityRequirements for %s type must be empty", srmee.Type)
}

// ErrMustOneOf is returned some value must be one of given list, but not one.
type ErrMustOneOf struct {
	Object      string
	ValidValues []string
}

func (ooe ErrMustOneOf) Error() string {
	return fmt.Sprintf("%s must be one of: %s", ooe.Object, strings.Join(ooe.ValidValues, ", "))
}
