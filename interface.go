package openapi

// Document represents a OpenAPI Specification document.
type Document struct {
	Version      string `yaml:"openapi"`
	Info         *Info
	Servers      []*Server
	Paths        Paths
	Components   *Components
	Security     *SecurityRequirement
	Tags         []*Tag
	ExternalDocs *ExternalDocumentation `yaml:"externalDocs"`
}

// Info Object
type Info struct {
	Title          string
	Description    string
	TermsOfService string `yaml:"termsOfService"`
	Contact        *Contact
	License        *License
	Version        string
}

// Contact Object
type Contact struct {
	Name  string
	URL   string
	Email string
}

// License Object
type License struct {
	Name string
	URL  string
}

// Server Object
type Server struct {
	URL         string
	Description string
	Variables   map[string]*ServerVariable
}

// ServerVariable Object
type ServerVariable struct {
	Enum        []string
	Default     string
	Description string
}

// Paths Object
type Paths map[string]*PathItem

// PathItem Object
type PathItem struct {
	Ref         string
	Summary     string
	Description string
	Get         *Operation
	Put         *Operation
	Post        *Operation
	Delete      *Operation
	Options     *Operation
	Head        *Operation
	Patch       *Operation
	Trace       *Operation
	Servers     []*Server
	Parameters  []*Parameter
}

// Operation Object
type Operation struct {
	Tags         []string
	Summary      string
	Description  string
	ExternalDocs *ExternalDocumentation `yaml:"externalDocs"`
	OperationID  string                 `yaml:"operationId"`
	Parameters   []*Parameter
	RequestBody  *RequestBody `yaml:"requestBody"`
	Responses    Responses
	Callbacks    map[string]*Callback
	Deprecated   bool
	Security     *SecurityRequirement
	Servers      []*Server
}

// Parameter Object
type Parameter struct {
	Name            string
	In              string
	Description     string
	Required        bool
	Deprecated      string
	AllowEmptyValue bool `yaml:"allowEmptyValue"`

	Style         string
	Explode       bool
	AllowReserved bool `yaml:"allowReserved"`
	Schema        *Schema
	Example       interface{}
	Examples      map[string]*Example

	Content map[string]*MediaType

	Ref string `yaml:"$ref"`
}

// RequestBody Object
type RequestBody struct {
	Description string
	Content     map[string]MediaType
	Required    bool

	Ref string `yaml:"$ref"`
}

// Responses Object
type Responses map[string]*Response

// Response Object
type Response struct {
	Description string
	Headers     map[string]*Header
	Content     map[string]*MediaType
	Links       map[string]*Link

	Ref string `yaml:"$ref"`
}

// Callbacks Object
type Callbacks map[string]*Callback

// Callback Object
type Callback map[string]*PathItem

// Schema Object
type Schema struct {
	Title            string
	MultipleOf       int `yaml:"multipleOf"`
	Maximum          int
	ExclusiveMaximum bool `yaml:"exclusiveMaximum"`
	Minimum          int
	ExclusiveMinimum bool `yaml:"exclusiveMinimum"`
	MaxLength        int  `yaml:"maxLength"`
	MinLength        int  `yaml:"minLength"`
	Pattern          string
	MaxItems         int `yaml:"maxItems"`
	MinItems         int `yaml:"minItems"`
	MaxProperties    int `yaml:"maxProperties"`
	MinProperties    int `yaml:"minProperties"`
	Required         []string
	Enum             []string

	Type                       string
	AllOf                      *Schema `yaml:"allOf"`
	OneOf                      *Schema `yaml:"oneOf"`
	AnyOf                      *Schema `yaml:"anyOf"`
	Not                        *Schema
	Items                      *Schema
	Properties                 map[string]*Schema
	EnableAdditionalProperties bool `yaml:"additionalProperties"`
	Description                string
	Format                     string
	Default                    string

	Nullable      bool
	Discriminator *Discriminator
	ReadOnly      bool `yaml:"readOnly"`
	WriteOnly     bool `yaml:"writeOnly"`
	XML           *XML
	ExternalDocs  *ExternalDocumentation `yaml:"externalDocs"`
	Example       interface{}
	Deprecated    bool

	Ref string `yaml:"$ref"`
}

// Example Object
type Example struct {
	Summary       string
	Description   string
	Value         interface{}
	ExternalValue interface{} `yaml:"externalValue"`

	Ref string `yaml:"$ref"`
}

// MediaType Object
type MediaType struct {
	Schema   *Schema
	Example  interface{}
	Examples map[string]*Example
	Encoding map[string]*Encoding

	Ref string `yaml:"$ref"`
}

// Header Object
type Header struct {
	Description     string
	Required        bool
	Deprecated      string
	AllowEmptyValue bool `yaml:"allowEmptyValue"`

	Style         string
	Explode       bool
	AllowReserved bool `yaml:"allowReserved"`
	Schema        Schema
	Example       interface{}
	Examples      map[string]*Example

	Content map[string]*MediaType

	Ref string `yaml:"$ref"`
}

// Link Object
type Link struct {
	OperationRef string `yaml:"operationRef"`
	OperationID  string `yaml:"operationId"`
	Parameters   map[string]interface{}
	RequestBody  interface{} `yaml:"requestBody"`
	Description  string
	Server       *Server

	Ref string `yaml:"$ref"`
}

// Encoding Object
type Encoding struct {
	ContentType   string `yaml:"contentType"`
	Headers       map[string]*Header
	Style         string
	Explode       bool
	AllowReserved bool `yaml:"allowReserved"`
}

// Discriminator Object
type Discriminator struct {
	PropertyName string `yaml:"propertyName"`
	Mapping      map[string]string
}

// XML Object
type XML struct {
	Name      string
	Namespace string
	Prefix    string
	Attribute bool
	Wrapped   bool
}

// Components Object
type Components struct {
	Schemas         map[string]*Schema
	Responses       Responses
	Parameters      map[string]*Parameter
	Examples        map[string]*Example
	RequestBodies   map[string]*RequestBody `yaml:"requestBodies"`
	Headers         map[string]*Header
	SecuritySchemes map[string]*SecurityScheme `yaml:"securitySchemes"`
	Links           map[string]*Link
	Callbacks       Callbacks
}

// SecurityScheme Object
type SecurityScheme struct {
	Type             string
	Description      string
	Name             string
	In               string
	Scheme           string
	BearerFormat     string `yaml:"bearerFormat"`
	Flows            *OAuthFlows
	OpenIDConnectURL string `yaml:"openIdConnectUrl"`

	Ref string `yaml:"$ref"`
}

// OAuthFlows Object
type OAuthFlows struct {
	Implicit          *OAuthFlow
	Password          *OAuthFlow
	ClientCredentials *OAuthFlow `yaml:"clientCredentials"`
	AuthorizationCode *OAuthFlow `yaml:"authorizationCode"`
}

// OAuthFlow Object
type OAuthFlow struct {
	AuthorizationURL string `yaml:"authorizationUrl"`
	TokenURL         string `yaml:"tokenUrl"`
	RefreshURL       string `yaml:"refreshUrl"`
	Scopes           map[string]string
}

// SecurityRequirement Object
type SecurityRequirement []map[string][]string

// Tag Object
type Tag struct {
	Name         string
	Description  string
	ExternalDocs *ExternalDocumentation `yaml:"externalDocs"`
}

// ExternalDocumentation Object
type ExternalDocumentation struct {
	Description string
	URL         string
}
