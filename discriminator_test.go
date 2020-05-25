package openapi

import (
	"errors"
	"strconv"
	"testing"

	yaml "github.com/goccy/go-yaml"
)

func TestDiscriminatorExampleUnmarshalYAML(t *testing.T) {
	t.Run("example", func(t *testing.T) {
		yml := `MyResponseType:
  oneOf:
  - $ref: '#/components/schemas/Cat'
  - $ref: '#/components/schemas/Dog'
  - $ref: '#/components/schemas/Lizard'
  discriminator:
    propertyName: petType`
		var got map[string]*Schema
		if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
			t.Fatal(err)
		}
		want := map[string]*Schema{
			"MyResponseType": {
				oneOf: []*Schema{
					{reference: "#/components/schemas/Cat"},
					{reference: "#/components/schemas/Dog"},
					{reference: "#/components/schemas/Lizard"},
				},
				discriminator: &Discriminator{
					propertyName: "petType",
				},
			},
		}
		assertEqual(t, got, want)
	})
	t.Run("mapping", func(t *testing.T) {
		yml := `MyResponseType:
  oneOf:
  - $ref: '#/components/schemas/Cat'
  - $ref: '#/components/schemas/Dog'
  - $ref: '#/components/schemas/Lizard'
  - $ref: 'https://gigantic-server.com/schemas/Monster/schema.json'
  discriminator:
    propertyName: petType
    mapping:
      dog: '#/components/schemas/Dog'
      monster: 'https://gigantic-server.com/schemas/Monster/schema.json'`
		var got map[string]*Schema
		if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
			t.Fatal(err)
		}
		want := map[string]*Schema{
			"MyResponseType": {
				oneOf: []*Schema{
					{reference: "#/components/schemas/Cat"},
					{reference: "#/components/schemas/Dog"},
					{reference: "#/components/schemas/Lizard"},
					{reference: "https://gigantic-server.com/schemas/Monster/schema.json"},
				},
				discriminator: &Discriminator{
					propertyName: "petType",
					mapping: map[string]string{
						"dog":     "#/components/schemas/Dog",
						"monster": "https://gigantic-server.com/schemas/Monster/schema.json",
					},
				},
			},
		}
		assertEqual(t, got, want)
	})
}

func TestDiscriminatorUnmarshalYAML(t *testing.T) {
	tests := []struct {
		yml  string
		want Discriminator
	}{}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got Discriminator
			if err := yaml.Unmarshal([]byte(tt.yml), &got); err != nil {
				t.Fatal(err)
			}
			assertEqual(t, got, tt.want)
		})
	}
}

func TestDiscriminatorUnmarshalYAMLError(t *testing.T) {
	tests := []struct {
		yml  string
		want error
	}{
		{
			yml:  `mapping: foobar`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml:  `foo: bar`,
			want: ErrUnknownKey("foo"),
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := yaml.Unmarshal([]byte(tt.yml), &Discriminator{})
			assertSameError(t, got, tt.want)
		})
	}
}
