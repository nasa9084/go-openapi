package openapi

import (
	"strings"
	"sync"
)

// codebeat:disable[TOO_MANY_IVARS]

// Schema Object
type Schema struct {
	Title            string
	MultipleOf       int `yaml:"multipleOf"`
	Maximum          int
	ExclusiveMaximum bool `yaml:"exclusiveMaximum"`
	Minimum          int
	ExclusiveMinimum bool `yaml:"exclusiveMinimum"`
	MaxLength        int  `yaml:"maxLength"`
	MinLength        int  `yaml:"minLength"`
	Pattern          string
	MaxItems         int `yaml:"maxItems"`
	MinItems         int `yaml:"minItems"`
	MaxProperties    int `yaml:"maxProperties"`
	MinProperties    int `yaml:"minProperties"`
	Required         []string
	Enum             []string

	Type                 string
	AllOf                []*Schema `yaml:"allOf"`
	OneOf                []*Schema `yaml:"oneOf"`
	AnyOf                []*Schema `yaml:"anyOf"`
	Not                  *Schema
	Items                *Schema
	Properties           map[string]*Schema
	AdditionalProperties *Schema `yaml:"additionalProperties"`
	Description          string
	Format               string
	Default              string

	Nullable      bool
	Discriminator *Discriminator
	ReadOnly      bool `yaml:"readOnly"`
	WriteOnly     bool `yaml:"writeOnly"`
	XML           *XML
	ExternalDocs  *ExternalDocumentation `yaml:"externalDocs"`
	Example       interface{}
	Deprecated    bool

	Ref string `yaml:"$ref"`
}

// Validate the values of Schema object.
func (schema Schema) Validate() error {
	validaters := []validater{}
	for _, s := range schema.AllOf {
		validaters = append(validaters, s)
	}
	for _, s := range schema.OneOf {
		validaters = append(validaters, s)
	}
	for _, s := range schema.AnyOf {
		validaters = append(validaters, s)
	}
	if schema.Not != nil {
		validaters = append(validaters, schema.Not)
	}
	if schema.Items != nil {
		validaters = append(validaters, schema.Items)
	}
	if schema.Discriminator != nil {
		validaters = append(validaters, schema.Discriminator)
	}
	if schema.XML != nil {
		validaters = append(validaters, schema.XML)
	}
	if schema.ExternalDocs != nil {
		validaters = append(validaters, schema.ExternalDocs)
	}
	for _, property := range schema.Properties {
		validaters = append(validaters, property)
	}
	if e, ok := schema.Example.(validater); ok {
		validaters = append(validaters, e)
	}
	return validateAll(validaters)
}

var typeMap sync.Map

// RegisterType registers type and format to go-type mapping for schema.GoType function.
func RegisterType(typ, format, gotype string) {
	formatMap, ok := typeMap.Load(typ)
	if !ok {
		formatMap = &sync.Map{}
		typeMap.Store(typ, formatMap)
	}
	formatMap.(*sync.Map).Store(format, gotype)
}

// DeregisterType removes type and format to go-type mapping.
func DeregisterType(typ, format string) {
	formatMap, ok := typeMap.Load(typ)
	if !ok {
		return
	}
	formatMap.(*sync.Map).Delete(format)
}

// LoadType loads registered go-type from type and format.
func LoadType(typ, format string) string {
	formatMap, ok := typeMap.Load(typ)
	if !ok {
		return ""
	}
	gotype, ok := formatMap.(*sync.Map).Load(format)
	if !ok {
		gotype, ok = formatMap.(*sync.Map).Load("")
		if !ok {
			return ""
		}
	}
	return gotype.(string)
}

func init() {
	RegisterType("integer", "", "int")
	RegisterType("integer", "int32", "int32")
	RegisterType("integer", "int64", "int64")
	RegisterType("number", "", "float64")
	RegisterType("number", "float", "float32")
	RegisterType("double", "double", "float64")
	RegisterType("string", "", "string")
	RegisterType("string", "byte", "[]byte")
	RegisterType("string", "binary", "[]byte")
	RegisterType("boolean", "", "bool")
	RegisterType("string", "date", "time.Time")
	RegisterType("string", "date-time", "time.Time")
	RegisterType("string", "password", "string")
}

// GoType returns go-type representation of given schema.
// If you need to use custom type and format, you can register using RegisterType function.
// (array and object cannot be overrided)
func (schema Schema) GoType() string {
	if schema.Ref != "" {
		ref := strings.Split(schema.Ref, "/")
		return ref[len(ref)-1]
	}
	switch schema.Type {
	case "object":
		var buf strings.Builder
		buf.WriteString("struct {")
		for name, prop := range schema.Properties {
			// call 4 times WriteString is faster than fmt.Fprintf
			buf.WriteString("\n")
			buf.WriteString(name)
			buf.WriteString(" ")
			buf.WriteString(prop.GoType())
		}
		buf.WriteString("\n}")
		return buf.String()
	case "array":
		return "[]" + schema.Items.GoType()
	default:
		return LoadType(schema.Type, schema.Format)
	}
}
