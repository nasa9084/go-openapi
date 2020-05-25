package openapi

import (
	"errors"
	"strconv"
	"testing"

	yaml "github.com/goccy/go-yaml"
)

func TestMediaTypeExampleUnmarshalYAML(t *testing.T) {
	yml := `application/json:
  schema:
    $ref: "#/components/schemas/Pet"
  examples:
    cat:
      summary: An example of a cat
      value:
        name: Fluffy
        petType: Cat
        color: White
        gender: male
        breed: Persian
    dog:
      summary: An example of a dog with a cat's name
      value:
        name: Puma
        petType: Dog
        color: Black
        gender: Female
        breed: Mixed
    frog:
      $ref: "#/components/examples/frog-example"`

	var got map[string]*MediaType
	if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
		t.Fatal(err)
	}

	want := map[string]*MediaType{
		"application/json": {
			schema: &Schema{
				reference: "#/components/schemas/Pet",
			},
			examples: map[string]*Example{
				"cat": {
					summary: "An example of a cat",
					value: map[string]interface{}{
						"name":    "Fluffy",
						"petType": "Cat",
						"color":   "White",
						"gender":  "male",
						"breed":   "Persian",
					},
				},
				"dog": {
					summary: "An example of a dog with a cat's name",
					value: map[string]interface{}{
						"name":    "Puma",
						"petType": "Dog",
						"color":   "Black",
						"gender":  "Female",
						"breed":   "Mixed",
					},
				},
				"frog": {
					reference: "#/components/examples/frog-example",
				},
			},
		},
	}
	assertEqual(t, got, want)
}

func TestMediaTypeUnmarshalYAML(t *testing.T) {
	tests := []struct {
		yml  string
		want MediaType
	}{
		{
			yml: `x-foo: bar`,
			want: MediaType{
				extension: map[string]interface{}{
					"x-foo": "bar",
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got MediaType
			if err := yaml.Unmarshal([]byte(tt.yml), &got); err != nil {
				t.Fatal(err)
			}
			assertEqual(t, got, tt.want)
		})
	}
}

func TestMediaTypeUnmarshalYAMLError(t *testing.T) {
	tests := []struct {
		yml  string
		want error
	}{
		{
			yml:  `examples: foo`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml:  `encoding: foo`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml:  `foo: bar`,
			want: ErrUnknownKey("foo"),
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := yaml.Unmarshal([]byte(tt.yml), &MediaType{})
			assertSameError(t, got, tt.want)
		})
	}
}
