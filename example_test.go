package openapi

import (
	"strconv"
	"testing"

	yaml "github.com/goccy/go-yaml"
)

func TestExampleExampleUnmarshalYAML(t *testing.T) {
	t.Run("in a request body", func(t *testing.T) {
		yml := `requestBody:
  content:
    'application/json':
      schema:
        $ref: '#/components/schemas/Address'
      examples:
        foo:
          summary: A foo example
          value: {"foo": "bar"}
        bar:
          summary: A bar example
          value: {"bar": "baz"}
    'application/xml':
      examples:
        xmlExample:
          summary: This is an example in XML
          externalValue: 'http://example.org/examples/address-example.xml'
    'text/plain':
      examples:
        textExample:
          summary: This is a text example
          externalValue: 'http://foo.bar/examples/address-example.txt'`

		var target struct {
			RequestBody RequestBody `yaml:"requestBody"`
		}
		if err := yaml.Unmarshal([]byte(yml), &target); err != nil {
			t.Fatal(err)
		}

		got := target.RequestBody
		want := RequestBody{
			content: map[string]*MediaType{
				"application/json": {
					schema: &Schema{
						reference: "#/components/schemas/Address",
					},
					examples: map[string]*Example{
						"foo": {
							summary: "A foo example",
							value: map[string]interface{}{
								"foo": "bar",
							},
						},
						"bar": {
							summary: "A bar example",
							value: map[string]interface{}{
								"bar": "baz",
							},
						},
					},
				},
				"application/xml": {
					examples: map[string]*Example{
						"xmlExample": {
							summary:       "This is an example in XML",
							externalValue: "http://example.org/examples/address-example.xml",
						},
					},
				},
				"text/plain": {
					examples: map[string]*Example{
						"textExample": {
							summary:       "This is a text example",
							externalValue: "http://foo.bar/examples/address-example.txt",
						},
					},
				},
			},
		}
		assertEqual(t, got, want)
	})
	t.Run("in a parameter", func(t *testing.T) {
		yml := `parameters:
  - name: 'zipCode'
    in: 'query'
    schema:
      type: 'string'
      format: 'zip-code'
    examples:
      zip-example:
        $ref: '#/components/examples/zip-example'`
		var target struct {
			Parameters []*Parameter
		}
		if err := yaml.Unmarshal([]byte(yml), &target); err != nil {
			t.Fatal(err)
		}
		got := target.Parameters
		want := []*Parameter{
			{
				name: "zipCode",
				in:   "query",
				schema: &Schema{
					type_:  "string",
					format: "zip-code",
				},
				examples: map[string]*Example{
					"zip-example": {
						reference: "#/components/examples/zip-example",
					},
				},
			},
		}
		assertEqual(t, got, want)
	})
	t.Run("in a response", func(t *testing.T) {
		yml := `responses:
  '200':
    description: your car appointment has been booked
    content:
      application/json:
        schema:
          $ref: '#/components/schemas/SuccessResponse'
        examples:
          confirmation-success:
            $ref: '#/components/examples/confirmation-success'`
		var target struct {
			Responses Responses
		}
		if err := yaml.Unmarshal([]byte(yml), &target); err != nil {
			t.Fatal(err)
		}
		got := target.Responses
		want := Responses{
			responses: map[string]*Response{
				"200": {
					description: "your car appointment has been booked",
					content: map[string]*MediaType{
						"application/json": {
							schema: &Schema{
								reference: "#/components/schemas/SuccessResponse",
							},
							examples: map[string]*Example{
								"confirmation-success": {
									reference: "#/components/examples/confirmation-success",
								},
							},
						},
					},
				},
			},
		}
		assertEqual(t, got, want)
	})
}

func TestExampleUnmarshalYAML(t *testing.T) {
	tests := []struct {
		yml  string
		want Example
	}{
		{
			yml: `description: foobar`,
			want: Example{
				description: "foobar",
			},
		},
		{
			yml: `x-foo: bar`,
			want: Example{
				extension: map[string]interface{}{
					"x-foo": "bar",
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got Example
			if err := yaml.Unmarshal([]byte(tt.yml), &got); err != nil {
				t.Fatal(err)
			}
			assertEqual(t, got, tt.want)
		})
	}
}

func TestExampleUnmarshalYAMLError(t *testing.T) {
	tests := []struct {
		yml  string
		want error
	}{
		{
			yml:  `foo: bar`,
			want: ErrUnknownKey("foo"),
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := yaml.Unmarshal([]byte(tt.yml), &Example{})
			assertSameError(t, got, tt.want)
		})
	}
}
