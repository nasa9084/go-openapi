package openapi

//+object
// OpenAPI is the root document object of the OpenAPI document.
type OpenAPI struct {
	openapi      string                 `required:"yes" format:"semver"`
	info         *Info                  `required:"yes"`
	servers      []*Server              `yaml:",omitempty"`
	paths        *Paths                 `required:"yes"`
	components   *Components            `yaml:",omitempty"`
	security     []*SecurityRequirement `yaml:",omitepmty"`
	tags         []*Tag                 `yaml:",omitempty"`
	externalDocs *ExternalDocumentation `yaml:",omitempty"`

	extension map[string]interface{} `yaml:",omitempty,inline" format:"prefix,x-"`
}

//+object
// Info provides metadata about the API.
type Info struct {
	root *OpenAPI `yaml:"-"`

	title          string   `required:"yes"`
	description    string   `yaml:",omitempty"`
	termsOfService string   `yaml:",omitempty" format:"url"`
	contact        *Contact `yaml:",omitempty"`
	license        *License `yaml:",omitempty"`
	version        string   `required:"yes"`

	extension map[string]interface{} `yaml:",omitempty,inline" format:"prefix,x-"`
}

//+object
// Contact information for the exposed API.
type Contact struct {
	root *OpenAPI `yaml:"-"`

	name  string `yaml:",omitempty"`
	url   string `yaml:",omitempty" format:"url"`
	email string `yaml:",omitempty" format:"email"`

	extension map[string]interface{} `yaml:",omitempty,inline" format:"prefix,x-"`
}

//+object
// License information for the exposed API.
type License struct {
	root *OpenAPI `yaml:"-"`

	name string `required:"yes"`
	url  string `yaml:",omitempty" format:"url"`

	extension map[string]interface{} `yaml:",omitempty,inline" format:"prefix,x-"`
}

//+object
// Server is an object representing a Server.
type Server struct {
	root *OpenAPI `yaml:"-"`

	url         string                     `required:"yes" format:"url,template"`
	description string                     `yaml:",omitempty"`
	variables   map[string]*ServerVariable `yaml:",omitempty"`

	extension map[string]interface{} `yaml:",omitempty,inline" format:"prefix,x-"`
}

//+object
// ServerVariable is an object representing a Server Variable for serverURL template substitution.
type ServerVariable struct {
	root *OpenAPI `yaml:"-"`

	enum        []string `yaml:",omitempty"`
	default_    string   `required:"yes" yaml:"default"` //nolint[golint]
	description string   `yaml:",omitempty"`

	extension map[string]interface{} `yaml:",omitempty,inline" format:"prefix,x-"`
}

//+object
// Components holds a set of reusable objects for different aspects of the OAS.
type Components struct {
	root *OpenAPI `yaml:"-"`

	schemas         map[string]*Schema         `yaml:",omitempty"`
	responses       map[string]*Response       `yaml:",omitempty"`
	parameters      map[string]*Parameter      `yaml:",omitempty"`
	examples        map[string]*Example        `yaml:",omitempty"`
	requestBodies   map[string]*RequestBody    `yaml:",omitempty"`
	headers         map[string]*Header         `yaml:",omitempty"`
	securitySchemes map[string]*SecurityScheme `yaml:",omitempty"`
	links           map[string]*Link           `yaml:",omitempty"`
	callbacks       map[string]*Callback       `yaml:",omitempty"`

	extension map[string]interface{} `yaml:",omitempty,inline" format:"prefix,x-"`
}

//+object
// Paths holds the relative paths to the individual endpoints and their operations.
type Paths struct {
	root *OpenAPI `yaml:"-"`

	paths map[string]*PathItem `yaml:",inline" format:"prefix,/"`

	extension map[string]interface{} `yaml:",omitempty,inline" format:"prefix,x-"`
}

//+object
// PathItem describes the operations available on a single path.
type PathItem struct {
	root *OpenAPI `yaml:"-"`

	summary     string       `yaml:",omitempty"`
	description string       `yaml:",omitempty"`
	get         *Operation   `yaml:",omitempty"`
	put         *Operation   `yaml:",omitempty"`
	post        *Operation   `yaml:",omitempty"`
	delete      *Operation   `yaml:",omitempty"`
	options     *Operation   `yaml:",omitempty"`
	head        *Operation   `yaml:",omitempty"`
	patch       *Operation   `yaml:",omitempty"`
	trace       *Operation   `yaml:",omitempty"`
	servers     []*Server    `yaml:",omitempty"`
	parameters  []*Parameter `yaml:",omitempty"`

	extension map[string]interface{} `yaml:",omitempty,inline" format:"prefix,x-"`
}

//+object
// Operation describes a single API operation on a path.
type Operation struct {
	root *OpenAPI `yaml:"-"`

	tags         []string               `yaml:",omitempty"`
	summary      string                 `yaml:",omitempty"`
	description  string                 `yaml:",omitempty"`
	externalDocs *ExternalDocumentation `yaml:",omitempty"`
	operationID  string                 `yaml:"operationId,omitempty"`
	parameters   []*Parameter           `yaml:",omitempty"`
	requestBody  *RequestBody           `yaml:",omitempty"`
	responses    *Responses             `required:"yes"`
	callbacks    map[string]*Callback   `yaml:",omitempty"`
	deprecated   bool                   `yaml:",omitempty"`
	security     []*SecurityRequirement `yaml:",omitempty"`
	servers      []*Server              `yaml:",omitempty"`

	extension map[string]interface{} `yaml:",omitempty,inline" format:"prefix,x-"`
}

//+object
// ExternalDocumentation allows referencing an external resource for extended documentation.
type ExternalDocumentation struct {
	root *OpenAPI `yaml:"-"`

	description string `yaml:",omitempty"`
	url         string `required:"yes" format:"url"`

	extension map[string]interface{} `yaml:",omitempty,inline" format:"prefix,x-"`
}

//+object
// Parameter describes a single operation parameter.
type Parameter struct {
	root *OpenAPI `yaml:"-"`

	name            string                `required:"yes"`
	in              string                `required:"yes" oneof:"query,header,path,cookie"`
	description     string                `yaml:",omitempty"`
	required        bool                  `yaml:",omitempty"`
	deprecated      bool                  `yaml:",omitempty"`
	allowEmptyValue bool                  `yaml:",omitempty"`
	style           string                `yaml:",omitempty"`
	explode         bool                  `yaml:",omitempty"`
	allowReserved   bool                  `yaml:",omitempty"`
	schema          *Schema               `yaml:",omitempty"`
	example         interface{}           `yaml:",omitempty"`
	examples        map[string]*Example   `yaml:",omitempty"`
	content         map[string]*MediaType `yaml:",omitempty"`

	extension map[string]interface{} `yaml:",omitempty,inline" format:"prefix,x-"`

	reference string     `yaml:"$ref,omitempty"`
	resolved  *Parameter `yaml:"-"`
}

//+object
// RequestBody describes a single request body.
type RequestBody struct {
	root *OpenAPI `yaml:"-"`

	description string                `yaml:",omitempty"`
	content     map[string]*MediaType `required:"yes"`
	required    bool                  `yaml:",omitempty"`

	extension map[string]interface{} `yaml:",omitempty,inline" format:"prefix,x-"`

	reference string       `yaml:"$ref,omitempty"`
	resolved  *RequestBody `yaml:"-"`
}

//+object
// MediaType provides schema and examples for the media type identified by its key.
type MediaType struct {
	root *OpenAPI `yaml:"-"`

	schema   *Schema              `yaml:",omitempty"`
	example  interface{}          `yaml:",omitempty"`
	examples map[string]*Example  `yaml:",omitempty"`
	encoding map[string]*Encoding `yaml:",omitempty"`

	extension map[string]interface{} `yaml:",omitempty,inline" format:"prefix,x-"`
}

//+object
// Encoding is a single encoding definition applied to a single schema property.
type Encoding struct {
	root *OpenAPI `yaml:"-"`

	contentType   string             `yaml:",omitempty"`
	headers       map[string]*Header `yaml:",omitempty"`
	style         string             `yaml:",omitempty"`
	explode       bool               `yaml:",omitempty"`
	allowReserved bool               `yaml:",omitempty"`

	extension map[string]interface{} `yaml:",omitempty,inline" format:"prefix,x-"`
}

//+object
// Responses is a container for the expected responses of an operation.
type Responses struct {
	root *OpenAPI `yaml:"-"`

	responses map[string]*Response `yaml:",omitempty,inline" format:"regexp,^[1-5]([0-9][0-9]|XX)|default$"`

	extension map[string]interface{} `yaml:",omitempty,inline" format:"prefix,x-"`
}

//+object
// Response describes a single response from an API Operation, including design-time,
//+object
// static links to operations based on the response.
type Response struct {
	root *OpenAPI `yaml:"-"`

	description string                `required:"yes"`
	headers     map[string]*Header    `yaml:",omitempty"`
	content     map[string]*MediaType `yaml:",omitempty"`
	links       map[string]*Link      `yaml:",omitempty"`

	extension map[string]interface{} `yaml:",omitempty,inline" format:"prefix,x-"`

	reference string    `yaml:"$ref,omitempty"`
	resolved  *Response `yaml:"-"`
}

//+object
// Callback is a map of possible out-of band callbacks relatedd to the parent operation.
type Callback struct {
	root *OpenAPI `yaml:"-"`

	callback map[string]*PathItem `yaml:",omitempty,inline" format:"runtime"`

	extension map[string]interface{} `yaml:",omitempty,inline" format:"prefix,x-"`

	reference string    `yaml:"$ref,omitempty"`
	resolved  *Callback `yaml:"-"`
}

//+object
// Example object represents an example.
type Example struct {
	root *OpenAPI `yaml:"-"`

	summary       string      `yaml:",omitempty"`
	description   string      `yaml:",omitempty"`
	value         interface{} `yaml:",omitempty"`
	externalValue string      `yaml:",omitempty"`

	extension map[string]interface{} `yaml:",omitempty,inline" format:"prefix,x-"`

	reference string   `yaml:"$ref,omitempty"`
	resolved  *Example `yaml:"-"`
}

//+object
// Link represents a possible design-time link for a response.
type Link struct {
	root *OpenAPI `yaml:"-"`

	operationRef string                 `yaml:",omitempty"`
	operationID  string                 `yaml:"operationId,omitempty"`
	parameters   map[string]interface{} `yaml:",omitempty"`
	requestBody  interface{}            `yaml:",omitempty"`
	description  string                 `yaml:",omitempty"`
	server       *Server                `yaml:",omitempty"`

	extension map[string]interface{} `yaml:",omitempty,inline" format:"prefix,x-"`

	reference string `yaml:"$ref,omitempty"`
	resolved  *Link  `yaml:"-"`
}

//+object
// Header object
type Header struct {
	root *OpenAPI `yaml:"-"`

	description     string                `yaml:",omitempty"`
	required        bool                  `yaml:",omitempty"`
	deprecated      bool                  `yaml:",omitempty"`
	allowEmptyValue bool                  `yaml:",omitempty"`
	style           string                `yaml:",omitempty"`
	explode         bool                  `yaml:",omitempty"`
	allowReserved   bool                  `yaml:",omitempty"`
	schema          *Schema               `yaml:",omitempty"`
	example         interface{}           `yaml:",omitempty"`
	examples        map[string]*Example   `yaml:",omitempty"`
	content         map[string]*MediaType `yaml:",omitempty"`

	extension map[string]interface{} `yaml:",omitempty,inline" format:"prefix,x-"`

	reference string  `yaml:"$ref,omitempty"`
	resolved  *Header `yaml:"-"`
}

//+object
// Tag adds metadata to a single tag that is used by the Operation Object.
type Tag struct {
	root *OpenAPI `yaml:"-"`

	name         string                 `required:"yes"`
	description  string                 `yaml:",omitempty"`
	externalDocs *ExternalDocumentation `yaml:",omitempty"`

	extension map[string]interface{} `yaml:",omitempty,inline" format:"prefix,x-"`
}

//+object
// Schema allows the definition of input and output data types.
type Schema struct {
	root *OpenAPI `yaml:"-"`

	title            string   `yaml:",omitempty"`
	multipleOf       int      `yaml:",omitempty"`
	maximum          int      `yaml:",omitempty"`
	exclusiveMaximum bool     `yaml:",omitempty"`
	minimum          int      `yaml:",omitempty"`
	exclusiveMinimum bool     `yaml:",omitempty"`
	maxLength        int      `yaml:",omitempty"`
	minLength        int      `yaml:",omitempty"`
	pattern          string   `yaml:",omitempty"`
	maxItems         int      `yaml:",omitempty"`
	minItems         int      `yaml:",omitempty"`
	maxProperties    int      `yaml:",omitempty"`
	minProperties    int      `yaml:",omitempty"`
	required         []string `yaml:",omitempty"`
	enum             []string `yaml:",omitempty"`

	type_                string             `yaml:"type,omitempty"` //nolint[golint]
	allOf                []*Schema          `yaml:",omitempty"`
	oneOf                []*Schema          `yaml:",omitempty"`
	anyOf                []*Schema          `yaml:",omitempty"`
	not                  *Schema            `yaml:",omitempty"`
	items                *Schema            `yaml:",omitempty"`
	properties           map[string]*Schema `yaml:",omitempty"`
	additionalProperties *Schema            `yaml:",omitempty"`
	description          string             `yaml:",omitempty"`
	format               string             `yaml:",omitempty"`
	default_             string             `yaml:"default,omitempty"` //nolint[golint]

	nullable      bool                   `yaml:",omitempty"`
	discriminator *Discriminator         `yaml:",omitempty"`
	readOnly      bool                   `yaml:",omitempty"`
	writeOnly     bool                   `yaml:",omitempty"`
	xml           *XML                   `yaml:",omitempty"`
	externalDocs  *ExternalDocumentation `yaml:",omitempty"`
	example       interface{}            `yaml:",omitempty"`
	deprecated    bool                   `yaml:",omitempty"`

	extension map[string]interface{} `yaml:",omitempty,inline" format:"prefix,x-"`

	reference string  `yaml:"$ref,omitempty"`
	resolved  *Schema `yaml:"-"`
}

//+object
// Discriminator object.
type Discriminator struct {
	root *OpenAPI `yaml:"-"`

	propertyName string            `yaml:",omitempty"`
	mapping      map[string]string `yaml:",omitempty"`
}

//+object
// XML is a metadata object that allows for more fine-tuned XML model definitions.
type XML struct {
	root *OpenAPI `yaml:"-"`

	name      string `yaml:",omitempty"`
	namespace string `yaml:",omitempty"`
	prefix    string `yaml:",omitempty"`
	attribute bool   `yaml:",omitempty"`
	wrapped   bool   `yaml:",omitempty"`

	extension map[string]interface{} `yaml:",omitempty,inline" format:"prefix,x-"`
}

//+object
// SecuritySchema defines a security scheme that can be used by the operations.
type SecurityScheme struct {
	root *OpenAPI `yaml:"-"`

	type_            string      `yaml:"type,omitempty" oneof:"apiKey,http,oauth2,openIdConnect"` //nolint[golint]
	description      string      `yaml:",omitempty"`
	name             string      `yaml:",omitempty"`
	in               string      `yaml:",omitempty" oneof:"query,header,cookie"`
	scheme           string      `yaml:",omitempty"`
	bearerFormat     string      `yaml:",omitempty"`
	flows            *OAuthFlows `yaml:",omitempty"`
	openIDConnectURL string      `yaml:"openIdConnectUrl,omitempty" format:"url"`

	extension map[string]interface{} `yaml:",omitempty,inline" format:"prefix,x-"`

	reference string          `yaml:"$ref,omitempty"`
	resolved  *SecurityScheme `yaml:"-"`
}

//+object
// OAuthFlows allows configuration of the supported OAuthFlows.
type OAuthFlows struct {
	root *OpenAPI `yaml:"-"`

	implicit          *OAuthFlow `yaml:",omitempty"`
	password          *OAuthFlow `yaml:",omitempty"`
	clientCredentials *OAuthFlow `yaml:",omitempty"`
	authorizationCode *OAuthFlow `yaml:",omitempty"`

	extension map[string]interface{} `yaml:",omitempty,inline" format:"prefix,x-"`
}

//+object
// OAuthFlow is configuration details for a supported OAuth Flow.
type OAuthFlow struct {
	root *OpenAPI `yaml:"-"`

	authorizationURL string            `yaml:"authorizationUrl,omitempty" format:"url"`
	tokenURL         string            `yaml:"tokenUrl,omitempty" format:"url"`
	refreshURL       string            `yaml:"refreshUrl,omitempty" format:"url"`
	scopes           map[string]string `yaml:",omitempty"`

	extension map[string]interface{} `yaml:",omitempty,inline" format:"prefix,x-"`
}

//+object
// SecurityRequirements is lists the required security schemes to execute this operation.
type SecurityRequirement struct {
	root *OpenAPI `yaml:"-"`

	securityRequirement map[string][]string `yaml:",inline"`
}
