package openapi

import (
	"strconv"
	"testing"

	yaml "github.com/goccy/go-yaml"
)

func TestPathsExampleUnmarshalYAML(t *testing.T) {
	yml := `/pets:
  get:
    description: Returns all pets from the system that the user has access to
    responses:
      '200':
        description: A list of pets.
        content:
          application/json:
            schema:
              type: array
              items:
                $ref: '#/components/schemas/pet'`

	var got Paths
	if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
		t.Fatal(err)
	}

	want := Paths{
		paths: map[string]*PathItem{
			"/pets": {
				get: &Operation{
					description: "Returns all pets from the system that the user has access to",
					responses: &Responses{
						responses: map[string]*Response{
							"200": {
								description: "A list of pets.",
								content: map[string]*MediaType{
									"application/json": {
										schema: &Schema{
											type_: "array",
											items: &Schema{
												reference: "#/components/schemas/pet",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	assertEqual(t, got, want)
}

func TestPathsUnmarshalYAML(t *testing.T) {
	tests := []struct {
		yml  string
		want Paths
	}{
		{
			yml: `x-foo: bar`,
			want: Paths{
				extension: map[string]interface{}{
					"x-foo": "bar",
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got Paths
			if err := yaml.Unmarshal([]byte(tt.yml), &got); err != nil {
				t.Fatal(err)
			}
			assertEqual(t, got, tt.want)
		})
	}
}

func TestPathsUnmarshalYAMLError(t *testing.T) {
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
			got := yaml.Unmarshal([]byte(tt.yml), &Paths{})
			assertSameError(t, got, tt.want)
		})
	}
}
