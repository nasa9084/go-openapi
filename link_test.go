package openapi

import (
	"errors"
	"strconv"
	"testing"

	yaml "github.com/goccy/go-yaml"
)

func TestLinkExampleUnmarshalYAML(t *testing.T) {
	t.Run("$request.path.id", func(t *testing.T) {
		yml := `paths:
  /users/{id}:
    parameters:
    - name: id
      in: path
      required: true
      description: the user identifier, as userId
      schema:
        type: string
    get:
      responses:
        '200':
          description: the user being returned
          content:
            application/json:
              schema:
                type: object
                properties:
                  uuid: # the unique user id
                    type: string
                    format: uuid
          links:
            address:
              # the target link operationId
              operationId: getUserAddress
              parameters:
                # get the "id" field from the request path parameter named "id"
                userId: $request.path.id
  # the path item of the linked operation
  /users/{userid}/address:
    parameters:
    - name: userid
      in: path
      required: true
      description: the user identifier, as userId
      schema:
        type: string
    # linked operation
    get:
      operationId: getUserAddress
      responses:
        '200':
          description: the user's address`
		var target struct {
			Paths Paths
		}
		if err := yaml.Unmarshal([]byte(yml), &target); err != nil {
			t.Fatal(err)
		}
		got := target.Paths
		want := Paths{
			paths: map[string]*PathItem{
				"/users/{id}": {
					parameters: []*Parameter{
						{
							name:        "id",
							in:          "path",
							required:    true,
							description: "the user identifier, as userId",
							schema:      &Schema{type_: "string"},
						},
					},
					get: &Operation{
						responses: &Responses{
							responses: map[string]*Response{
								"200": {
									description: "the user being returned",
									content: map[string]*MediaType{
										"application/json": {
											schema: &Schema{
												type_: "object",
												properties: map[string]*Schema{
													"uuid": {
														type_:  "string",
														format: "uuid",
													},
												},
											},
										},
									},
									links: map[string]*Link{
										"address": {
											operationID: "getUserAddress",
											parameters: map[string]interface{}{
												"userId": "$request.path.id",
											},
										},
									},
								},
							},
						},
					},
				},
				"/users/{userid}/address": {
					parameters: []*Parameter{
						{
							name:        "userid",
							in:          "path",
							required:    true,
							description: "the user identifier, as userId",
							schema:      &Schema{type_: "string"},
						},
					},
					get: &Operation{
						operationID: "getUserAddress",
						responses: &Responses{
							responses: map[string]*Response{
								"200": {description: "the user's address"},
							},
						},
					},
				},
			},
		}
		assertEqual(t, got, want)
	})
	t.Run("can use values from the response body", func(t *testing.T) {
		yml := `links:
  address:
    operationId: getUserAddressByUUID
    parameters:
      # get the "uuid" field from the "uuid" field in the response body
      userUuid: $response.body#/uuid`
		var target struct {
			Links map[string]*Link
		}
		if err := yaml.Unmarshal([]byte(yml), &target); err != nil {
			t.Fatal(err)
		}
		got := target.Links
		want := map[string]*Link{
			"address": {
				operationID: "getUserAddressByUUID",
				parameters: map[string]interface{}{
					"userUuid": "$response.body#/uuid",
				},
			},
		}
		assertEqual(t, got, want)
	})
	t.Run("relative operationRef", func(t *testing.T) {
		yml := `links:
  UserRepositories:
    # returns array of '#/components/schemas/repository'
    operationRef: '#/paths/~12.0~1repositories~1{username}/get'
    parameters:
      username: $response.body#/username`
		var target struct {
			Links map[string]*Link
		}
		if err := yaml.Unmarshal([]byte(yml), &target); err != nil {
			t.Fatal(err)
		}
		got := target.Links
		want := map[string]*Link{
			"UserRepositories": {
				operationRef: "#/paths/~12.0~1repositories~1{username}/get",
				parameters: map[string]interface{}{
					"username": "$response.body#/username",
				},
			},
		}
		assertEqual(t, got, want)
	})
	t.Run("absolute operationRef", func(t *testing.T) {
		yml := `links:
  UserRepositories:
    # returns array of '#/components/schemas/repository'
    operationRef: 'https://na2.gigantic-server.com/#/paths/~12.0~1repositories~1{username}/get'
    parameters:
      username: $response.body#/username`
		var target struct {
			Links map[string]*Link
		}
		if err := yaml.Unmarshal([]byte(yml), &target); err != nil {
			t.Fatal(err)
		}
		got := target.Links
		want := map[string]*Link{
			"UserRepositories": {
				operationRef: "https://na2.gigantic-server.com/#/paths/~12.0~1repositories~1{username}/get",
				parameters: map[string]interface{}{
					"username": "$response.body#/username",
				},
			},
		}
		assertEqual(t, got, want)
	})
}

func TestLinkUnmarshalYAML(t *testing.T) {
	tests := []struct {
		yml  string
		want Link
	}{
		{
			yml: `requestBody: {}`,
			want: Link{
				requestBody: map[string]interface{}{},
			},
		},
		{
			yml: `description: foo`,
			want: Link{
				description: "foo",
			},
		},
		{
			yml: `server:
  url: example.com`,
			want: Link{
				server: &Server{
					url: "example.com",
				},
			},
		},
		{
			yml: `x-foo: bar`,
			want: Link{
				extension: map[string]interface{}{
					"x-foo": "bar",
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var got Link
			if err := yaml.Unmarshal([]byte(tt.yml), &got); err != nil {
				t.Fatal(err)
			}
			assertEqual(t, got, tt.want)
		})
	}
}

func TestLinkUnmarshalYAMLError(t *testing.T) {
	tests := []struct {
		yml  string
		want error
	}{
		{
			yml:  `parameters: foobar`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml:  `server: foobar`,
			want: errors.New("String node doesn't MapNode"),
		},
		{
			yml:  `foo: bar`,
			want: ErrUnknownKey("foo"),
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := yaml.Unmarshal([]byte(tt.yml), &Link{})
			assertSameError(t, got, tt.want)
		})
	}
}
