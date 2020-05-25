package openapi

import (
	"errors"
	"strconv"
	"testing"

	yaml "github.com/goccy/go-yaml"
)

func TestEncodingExampleUnmarshalYAML(t *testing.T) {
	yml := `requestBody:
  content:
    multipart/mixed:
      schema:
        type: object
        properties:
          id:
            # default is text/plain
            type: string
            format: uuid
          address:
            # default is application/json
            type: object
            properties: {}
          historyMetadata:
            # need to declare XML format!
            description: metadata in XML format
            type: object
            properties: {}
          profileImage:
            # default is application/octet-stream, need to declare an image type only!
            type: string
            format: binary
      encoding:
        historyMetadata:
          # require XML Content-Type in utf-8 encoding
          contentType: application/xml; charset=utf-8
        profileImage:
          # only accept png/jpeg
          contentType: image/png, image/jpeg
          headers:
            X-Rate-Limit-Limit:
              description: The number of allowed requests in the current period
              schema:
                type: integer`

	var target struct {
		RequestBody RequestBody `yaml:"requestBody"`
	}

	if err := yaml.Unmarshal([]byte(yml), &target); err != nil {
		t.Fatal(err)
	}

	got := target.RequestBody
	want := RequestBody{
		content: map[string]*MediaType{
			"multipart/mixed": {
				schema: &Schema{
					type_: "object",
					properties: map[string]*Schema{
						"id": {
							type_:  "string",
							format: "uuid",
						},
						"address": {
							type_:      "object",
							properties: map[string]*Schema{},
						},
						"historyMetadata": {
							description: "metadata in XML format",
							type_:       "object",
							properties:  map[string]*Schema{},
						},
						"profileImage": {
							type_:  "string",
							format: "binary",
						},
					},
				},
				encoding: map[string]*Encoding{
					"historyMetadata": {
						contentType: "application/xml; charset=utf-8",
					},
					"profileImage": {
						contentType: "image/png, image/jpeg",
						headers: map[string]*Header{
							"X-Rate-Limit-Limit": {
								description: "The number of allowed requests in the current period",
								schema: &Schema{
									type_: "integer",
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

func TestEncodingUnmarshalYAML(t *testing.T) {
	tests := []struct {
		yml  string
		want Encoding
	}{
		{
			yml: `style: foobar`,
			want: Encoding{
				style: "foobar",
			},
		},
		{
			yml: `explode: true`,
			want: Encoding{
				explode: true,
			},
		},
		{
			yml: `allowReserved: true`,
			want: Encoding{
				allowReserved: true,
			},
		},
		{
			yml: `x-foo: bar`,
			want: Encoding{
				extension: map[string]interface{}{
					"x-foo": "bar",
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got Encoding
			if err := yaml.Unmarshal([]byte(tt.yml), &got); err != nil {
				t.Fatal(err)
			}
			assertEqual(t, got, tt.want)
		})
	}
}

func TestEncodingUnmarshalYAMLError(t *testing.T) {
	tests := []struct {
		yml  string
		want error
	}{
		{
			yml:  `headers: foo`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml:  `foo: bar`,
			want: ErrUnknownKey("foo"),
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := yaml.Unmarshal([]byte(tt.yml), &Encoding{})
			assertSameError(t, got, tt.want)
		})
	}
}
