package openapi

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

	Type                       string
	AllOf                      *Schema `yaml:"allOf"`
	OneOf                      *Schema `yaml:"oneOf"`
	AnyOf                      *Schema `yaml:"anyOf"`
	Not                        *Schema
	Items                      *Schema
	Properties                 map[string]*Schema
	EnableAdditionalProperties bool `yaml:"additionalProperties"`
	Description                string
	Format                     string
	Default                    string

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
	if schema.AllOf != nil {
		validaters = append(validaters, schema.AllOf)
	}
	if schema.OneOf != nil {
		validaters = append(validaters, schema.OneOf)
	}
	if schema.AnyOf != nil {
		validaters = append(validaters, schema.AnyOf)
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
