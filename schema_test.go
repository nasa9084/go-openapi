package openapi

import (
	"errors"
	"strconv"
	"testing"

	yaml "github.com/goccy/go-yaml"
)

func TestSchemaExampleUnmarshalYAML(t *testing.T) {
	t.Run("primitive", testSchemaExampleUnmarshalYAMLPrimitive)
	t.Run("simple model", testSchemaExampleUnmarshalYAMLSimpleModel)
	t.Run("simple string to string map", testSchemaExampleUnmarshalYAMLStringToStringMap)
	t.Run("string to model map", testSchemaExampleUnmarshalYAMLStringToModelMap)
	t.Run("model with example", testSchemaExampleUnmarshalYAMLModelExample)
	t.Run("models with composition", testSchemaExampleUnmarshalYAMLComposition)
	t.Run("models with polymorphism support", testSchemaExampleUnmarshalYAMLPolymorphism)
}

func testSchemaExampleUnmarshalYAMLPrimitive(t *testing.T) {
	yml := `type: string
format: email`

	var got Schema
	if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
		t.Fatal(err)
	}

	want := Schema{
		type_:  "string",
		format: "email",
	}
	assertEqual(t, got, want)
}

func testSchemaExampleUnmarshalYAMLSimpleModel(t *testing.T) {
	yml := `type: object
required:
- name
properties:
  name:
    type: string
  address:
    $ref: '#/components/schemas/Address'
  age:
    type: integer
    format: int32
    minimum: 0`

	var got Schema
	if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
		t.Fatal(err)
	}

	want := Schema{
		type_:    "object",
		required: []string{"name"},
		properties: map[string]*Schema{
			"name":    {type_: "string"},
			"address": {reference: "#/components/schemas/Address"},
			"age": {
				type_:   "integer",
				format:  "int32",
				minimum: 0,
			},
		},
	}
	assertEqual(t, got, want)
}

func testSchemaExampleUnmarshalYAMLStringToStringMap(t *testing.T) {
	yml := `type: object
additionalProperties:
  type: string`

	var got Schema
	if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
		t.Fatal(err)
	}

	want := Schema{
		type_: "object",
		additionalProperties: &Schema{
			type_: "string",
		},
	}
	assertEqual(t, got, want)
}

func testSchemaExampleUnmarshalYAMLStringToModelMap(t *testing.T) {
	yml := `type: object
additionalProperties:
  $ref: '#/components/schemas/ComplexModel'`

	var got Schema
	if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
		t.Fatal(err)
	}

	want := Schema{
		type_: "object",
		additionalProperties: &Schema{
			reference: "#/components/schemas/ComplexModel",
		},
	}
	assertEqual(t, got, want)
}

func testSchemaExampleUnmarshalYAMLModelExample(t *testing.T) {
	yml := `type: object
properties:
  id:
    type: integer
    format: int64
  name:
    type: string
required:
- name
example:
  name: Puma
  id: 1`

	var got Schema
	if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
		t.Fatal(err)
	}

	want := Schema{
		type_: "object",
		properties: map[string]*Schema{
			"id": {
				type_:  "integer",
				format: "int64",
			},
			"name": {
				type_: "string",
			},
		},
		required: []string{"name"},
		example: map[string]interface{}{
			"name": "Puma",
			"id":   uint64(1),
		},
	}
	assertEqual(t, got, want)
}

func testSchemaExampleUnmarshalYAMLComposition(t *testing.T) {
	yml := `components:
  schemas:
    ErrorModel:
      type: object
      required:
      - message
      - code
      properties:
        message:
          type: string
        code:
          type: integer
          minimum: 100
          maximum: 600
    ExtendedErrorModel:
      allOf:
      - $ref: '#/components/schemas/ErrorModel'
      - type: object
        required:
        - rootCause
        properties:
          rootCause:
            type: string`

	var target struct {
		Components Components
	}

	if err := yaml.Unmarshal([]byte(yml), &target); err != nil {
		t.Fatal(err)
	}

	got := target.Components
	want := Components{
		schemas: map[string]*Schema{
			"ErrorModel": {
				type_:    "object",
				required: []string{"message", "code"},
				properties: map[string]*Schema{
					"message": {
						type_: "string",
					},
					"code": {
						type_:   "integer",
						minimum: 100,
						maximum: 600,
					},
				},
			},
			"ExtendedErrorModel": {
				allOf: []*Schema{
					{
						reference: "#/components/schemas/ErrorModel",
					},
					{
						type_:    "object",
						required: []string{"rootCause"},
						properties: map[string]*Schema{
							"rootCause": {type_: "string"},
						},
					},
				},
			},
		},
	}
	assertEqual(t, got, want)
}

func testSchemaExampleUnmarshalYAMLPolymorphism(t *testing.T) {
	yml := `components:
  schemas:
    Pet:
      type: object
      discriminator:
        propertyName: petType
      properties:
        name:
          type: string
        petType:
          type: string
      required:
      - name
      - petType
    Cat:  ## "Cat" will be used as the discriminator value
      description: A representation of a cat
      allOf:
      - $ref: '#/components/schemas/Pet'
      - type: object
        properties:
          huntingSkill:
            type: string
            description: The measured skill for hunting
            enum:
            - clueless
            - lazy
            - adventurous
            - aggressive
        required:
        - huntingSkill
    Dog:  ## "Dog" will be used as the discriminator value
      description: A representation of a dog
      allOf:
      - $ref: '#/components/schemas/Pet'
      - type: object
        properties:
          packSize:
            type: integer
            format: int32
            description: the size of the pack the dog is from
            default: 0
            minimum: 0
        required:
        - packSize`

	var target struct {
		Components Components
	}

	if err := yaml.Unmarshal([]byte(yml), &target); err != nil {
		t.Fatal(err)
	}

	got := target.Components
	want := Components{
		schemas: map[string]*Schema{
			"Pet": {
				type_: "object",
				discriminator: &Discriminator{
					propertyName: "petType",
				},
				properties: map[string]*Schema{
					"name":    {type_: "string"},
					"petType": {type_: "string"},
				},
				required: []string{"name", "petType"},
			},
			"Cat": {
				description: "A representation of a cat",
				allOf: []*Schema{
					{reference: "#/components/schemas/Pet"},
					{
						type_: "object",
						properties: map[string]*Schema{
							"huntingSkill": {
								type_:       "string",
								description: "The measured skill for hunting",
								enum:        []string{"clueless", "lazy", "adventurous", "aggressive"},
							},
						},
						required: []string{"huntingSkill"},
					},
				},
			},
			"Dog": {
				description: "A representation of a dog",
				allOf: []*Schema{
					{reference: "#/components/schemas/Pet"},
					{
						type_: "object",
						properties: map[string]*Schema{
							"packSize": {
								type_:       "integer",
								format:      "int32",
								description: "the size of the pack the dog is from",
								default_:    "0",
								minimum:     0,
							},
						},
						required: []string{"packSize"},
					},
				},
			},
		},
	}
	assertEqual(t, got, want)
}

func TestSchemaUnmarshalYAML(t *testing.T) {
	tests := []struct {
		yml  string
		want Schema
	}{
		{
			yml: `title: Foo API`,
			want: Schema{
				title: "Foo API",
			},
		},
		{
			yml: `multipleOf: 3`,
			want: Schema{
				multipleOf: 3,
			},
		},
		{
			yml: `exclusiveMaximum: true`,
			want: Schema{
				exclusiveMaximum: true,
			},
		},
		{
			yml: `exclusiveMinimum: true`,
			want: Schema{
				exclusiveMinimum: true,
			},
		},
		{
			yml: `minLength: 3
maxLength: 6`,
			want: Schema{
				minLength: 3,
				maxLength: 6,
			},
		},
		{
			yml: `pattern: ^foo.+$`,
			want: Schema{
				pattern: "^foo.+$",
			},
		},
		{
			yml: `minItems: 3
maxItems: 6`,
			want: Schema{
				minItems: 3,
				maxItems: 6,
			},
		},
		{
			yml: `minProperties: 3
maxProperties: 6`,
			want: Schema{
				minProperties: 3,
				maxProperties: 6,
			},
		},
		{
			yml: `anyOf:
- description: foo`,
			want: Schema{
				anyOf: []*Schema{
					{
						description: "foo",
					},
				},
			},
		},
		{
			yml: `not:
  description: foo`,
			want: Schema{
				not: &Schema{
					description: "foo",
				},
			},
		},
		{
			yml: `nullable: true`,
			want: Schema{
				nullable: true,
			},
		},
		{
			yml: `writeOnly: true
readOnly: true`,
			want: Schema{
				writeOnly: true,
				readOnly:  true,
			},
		},
		{
			yml: `externalDocs:
  url: https://example.com`,
			want: Schema{
				externalDocs: &ExternalDocumentation{
					url: "https://example.com",
				},
			},
		},
		{
			yml: `deprecated: true`,
			want: Schema{
				deprecated: true,
			},
		},
		{
			yml: `x-foo: bar`,
			want: Schema{
				extension: map[string]interface{}{
					"x-foo": "bar",
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got Schema
			if err := yaml.Unmarshal([]byte(tt.yml), &got); err != nil {
				t.Fatal(err)
			}
			assertEqual(t, got, tt.want)
		})
	}
}

func TestSchemaUnmarshalYAMLError(t *testing.T) {
	tests := []struct {
		yml  string
		want error
	}{
		{
			yml:  `enum: foobar`,
			want: errors.New("String node doesn't ArrayNode"),
		},
		{
			yml:  `allOf: foobar`,
			want: errors.New("String node doesn't ArrayNode"),
		},
		{
			yml:  `oneOf: foobar`,
			want: errors.New("String node doesn't ArrayNode"),
		},
		{
			yml:  `anyOf: foobar`,
			want: errors.New("String node doesn't ArrayNode"),
		},
		{
			yml:  `not: foobar`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml:  `properties: foobar`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml:  `additionalProperties: foobar`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml:  `discriminator: foobar`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml:  `xml: foobar`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml:  `externalDocs: foobar`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml:  `foo: bar`,
			want: ErrUnknownKey("foo"),
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := yaml.Unmarshal([]byte(tt.yml), &Schema{})
			assertSameError(t, got, tt.want)
		})
	}
}
