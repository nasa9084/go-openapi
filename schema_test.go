package openapi_test

import (
	"reflect"
	"strconv"
	"testing"

	openapi "github.com/nasa9084/go-openapi"
	yaml "gopkg.in/yaml.v2"
)

func TestSchema_Validate(t *testing.T) {
	candidates := []candidate{
		{"empty", openapi.Schema{}, nil},
	}
	testValidater(t, candidates)
}

func TestSchameUnmarshal(t *testing.T) {
	tests := []struct {
		data string
		want *openapi.Schema
	}{
		{
			data: `---
type: string
`,
			want: &openapi.Schema{
				Type: "string",
			},
		},
		{
			data: `---
type: string
x-extension: val`,
			want: &openapi.Schema{
				Type: "string",
				Extension: map[string]interface{}{
					"x-extension": "val",
				},
			},
		},
		{
			data: `---
type: array
items:
  type: string
default:
  - foo
  - bar`,
			want: &openapi.Schema{
				Type: "array",
				Items: &openapi.Schema{
					Type: "string",
				},
			},
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got openapi.Schema
			err := yaml.Unmarshal([]byte(tt.data), &got)
			if err != nil {
				t.Fatal(err)
			}

			if got.Type != tt.want.Type {
				t.Errorf("%s != %s", got.Type, tt.want.Type)
				return
			}
			if !reflect.DeepEqual(got.Extension, tt.want.Extension) {
				t.Errorf("%v != %v", got.Extension, tt.want.Extension)
				return
			}
		})
	}
}
