package openapi

import (
	"errors"
	"strconv"
	"testing"

	yaml "github.com/goccy/go-yaml"
)

func TestCallbackExampleUnmarshalYAML(t *testing.T) {
	yml := `myWebhook:
  'http://notificationServer.com?transactionId={$request.body#/id}&email={$request.body#/email}':
    post:
      requestBody:
        description: Callback payload
        content:
          'application/json':
            schema:
              $ref: '#/components/schemas/SomePayload'
      responses:
        '200':
          description: webhook successfully processed and no retries will be performed`

	var got map[string]*Callback
	if err := yaml.Unmarshal([]byte(yml), &got); err != nil {
		t.Fatal(err)
	}

	want := map[string]*Callback{
		"myWebhook": {
			callback: map[string]*PathItem{
				"http://notificationServer.com?transactionId={$request.body#/id}&email={$request.body#/email}": {
					post: &Operation{
						requestBody: &RequestBody{
							description: "Callback payload",
							content: map[string]*MediaType{
								"application/json": {
									schema: &Schema{
										reference: "#/components/schemas/SomePayload",
									},
								},
							},
						},
						responses: &Responses{
							responses: map[string]*Response{
								"200": {
									description: "webhook successfully processed and no retries will be performed",
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

func TestCallbackUnmarshalYAML(t *testing.T) {
	tests := []struct {
		yml  string
		want Callback
	}{
		{
			yml: `x-foo: bar`,
			want: Callback{
				extension: map[string]interface{}{
					"x-foo": "bar",
				},
			},
		},
		{
			yml: `$ref: "#/components/callbacks/foo"`,
			want: Callback{
				reference: "#/components/callbacks/foo",
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got Callback
			if err := yaml.Unmarshal([]byte(tt.yml), &got); err != nil {
				t.Fatal(err)
			}
			assertEqual(t, got, tt.want)
		})
	}
}

func TestCallbackUnmarshalYAMLError(t *testing.T) {
	tests := []struct {
		yml  string
		want error
	}{
		{
			yml:  `$url: foobar`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml:  `foo: bar`,
			want: ErrUnknownKey("foo"),
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := yaml.Unmarshal([]byte(tt.yml), &Callback{})
			assertSameError(t, got, tt.want)
		})
	}
}
