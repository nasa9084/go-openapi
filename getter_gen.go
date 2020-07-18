// Code generated by GetterGenerator. DO NOT EDIT.

package openapi

func (v *OpenAPI) OpenAPI() string {
	return v.openapi
}

func (v *OpenAPI) Info() *Info {
	if v.info == nil {
		return &Info{}
	}
	return v.info
}

func (v *OpenAPI) Servers() []*Server {
	return v.servers
}

func (v *OpenAPI) Paths() *Paths {
	if v.paths == nil {
		return &Paths{}
	}
	return v.paths
}

func (v *OpenAPI) Components() *Components {
	if v.components == nil {
		return &Components{}
	}
	return v.components
}

func (v *OpenAPI) Security() []*SecurityRequirement {
	return v.security
}

func (v *OpenAPI) Tags() []*Tag {
	return v.tags
}

func (v *OpenAPI) ExternalDocs() *ExternalDocumentation {
	if v.externalDocs == nil {
		return &ExternalDocumentation{}
	}
	return v.externalDocs
}

func (v *OpenAPI) Extension(key string) interface{} {
	return v.extension[key]
}

func (v *Info) Root() *OpenAPI {
	if v.root == nil {
		return &OpenAPI{}
	}
	return v.root
}

func (v *Info) Title() string {
	return v.title
}

func (v *Info) Description() string {
	return v.description
}

func (v *Info) TermsOfService() string {
	return v.termsOfService
}

func (v *Info) Contact() *Contact {
	if v.contact == nil {
		return &Contact{}
	}
	return v.contact
}

func (v *Info) License() *License {
	if v.license == nil {
		return &License{}
	}
	return v.license
}

func (v *Info) Version() string {
	return v.version
}

func (v *Info) Extension(key string) interface{} {
	return v.extension[key]
}

func (v *Contact) Root() *OpenAPI {
	if v.root == nil {
		return &OpenAPI{}
	}
	return v.root
}

func (v *Contact) Name() string {
	return v.name
}

func (v *Contact) URL() string {
	return v.url
}

func (v *Contact) Email() string {
	return v.email
}

func (v *Contact) Extension(key string) interface{} {
	return v.extension[key]
}

func (v *License) Root() *OpenAPI {
	if v.root == nil {
		return &OpenAPI{}
	}
	return v.root
}

func (v *License) Name() string {
	return v.name
}

func (v *License) URL() string {
	return v.url
}

func (v *License) Extension(key string) interface{} {
	return v.extension[key]
}

func (v *Server) Root() *OpenAPI {
	if v.root == nil {
		return &OpenAPI{}
	}
	return v.root
}

func (v *Server) URL() string {
	return v.url
}

func (v *Server) Description() string {
	return v.description
}

func (v *Server) Variables() map[string]*ServerVariable {
	return v.variables
}

func (v *Server) Extension(key string) interface{} {
	return v.extension[key]
}

func (v *ServerVariable) Root() *OpenAPI {
	if v.root == nil {
		return &OpenAPI{}
	}
	return v.root
}

func (v *ServerVariable) Enum() []string {
	return v.enum
}

func (v *ServerVariable) Default() string {
	return v.default_
}

func (v *ServerVariable) Description() string {
	return v.description
}

func (v *ServerVariable) Extension(key string) interface{} {
	return v.extension[key]
}

func (v *Components) Root() *OpenAPI {
	if v.root == nil {
		return &OpenAPI{}
	}
	return v.root
}

func (v *Components) Schemas() map[string]*Schema {
	return v.schemas
}

func (v *Components) Responses() map[string]*Response {
	return v.responses
}

func (v *Components) Parameters() map[string]*Parameter {
	return v.parameters
}

func (v *Components) Examples() map[string]*Example {
	return v.examples
}

func (v *Components) RequestBodies() map[string]*RequestBody {
	return v.requestBodies
}

func (v *Components) Headers() map[string]*Header {
	return v.headers
}

func (v *Components) SecuritySchemes() map[string]*SecurityScheme {
	return v.securitySchemes
}

func (v *Components) Links() map[string]*Link {
	return v.links
}

func (v *Components) Callbacks() map[string]*Callback {
	return v.callbacks
}

func (v *Components) Extension(key string) interface{} {
	return v.extension[key]
}

func (v *Paths) Root() *OpenAPI {
	if v.root == nil {
		return &OpenAPI{}
	}
	return v.root
}

func (v *Paths) Get(key string) *PathItem {
	if val, ok := v.paths[key]; ok {
		return val
	}
	return &PathItem{}
}

func (v *Paths) Extension(key string) interface{} {
	return v.extension[key]
}

func (v *PathItem) Root() *OpenAPI {
	if v.root == nil {
		return &OpenAPI{}
	}
	return v.root
}

func (v *PathItem) Summary() string {
	return v.summary
}

func (v *PathItem) Description() string {
	return v.description
}

func (v *PathItem) Get() *Operation {
	if v.get == nil {
		return &Operation{}
	}
	return v.get
}

func (v *PathItem) Put() *Operation {
	if v.put == nil {
		return &Operation{}
	}
	return v.put
}

func (v *PathItem) Post() *Operation {
	if v.post == nil {
		return &Operation{}
	}
	return v.post
}

func (v *PathItem) Delete() *Operation {
	if v.delete == nil {
		return &Operation{}
	}
	return v.delete
}

func (v *PathItem) Options() *Operation {
	if v.options == nil {
		return &Operation{}
	}
	return v.options
}

func (v *PathItem) Head() *Operation {
	if v.head == nil {
		return &Operation{}
	}
	return v.head
}

func (v *PathItem) Patch() *Operation {
	if v.patch == nil {
		return &Operation{}
	}
	return v.patch
}

func (v *PathItem) Trace() *Operation {
	if v.trace == nil {
		return &Operation{}
	}
	return v.trace
}

func (v *PathItem) Servers() []*Server {
	return v.servers
}

func (v *PathItem) Parameters() []*Parameter {
	return v.parameters
}

func (v *PathItem) Extension(key string) interface{} {
	return v.extension[key]
}

func (v *Operation) Root() *OpenAPI {
	if v.root == nil {
		return &OpenAPI{}
	}
	return v.root
}

func (v *Operation) Tags() []string {
	return v.tags
}

func (v *Operation) Summary() string {
	return v.summary
}

func (v *Operation) Description() string {
	return v.description
}

func (v *Operation) ExternalDocs() *ExternalDocumentation {
	if v.externalDocs == nil {
		return &ExternalDocumentation{}
	}
	return v.externalDocs
}

func (v *Operation) OperationID() string {
	return v.operationID
}

func (v *Operation) Parameters() []*Parameter {
	return v.parameters
}

func (v *Operation) RequestBody() *RequestBody {
	if v.requestBody == nil {
		return &RequestBody{}
	}
	return v.requestBody
}

func (v *Operation) Responses() *Responses {
	if v.responses == nil {
		return &Responses{}
	}
	return v.responses
}

func (v *Operation) Callbacks() map[string]*Callback {
	return v.callbacks
}

func (v *Operation) Deprecated() bool {
	return v.deprecated
}

func (v *Operation) Security() []*SecurityRequirement {
	return v.security
}

func (v *Operation) Servers() []*Server {
	return v.servers
}

func (v *Operation) Extension(key string) interface{} {
	return v.extension[key]
}

func (v *ExternalDocumentation) Root() *OpenAPI {
	if v.root == nil {
		return &OpenAPI{}
	}
	return v.root
}

func (v *ExternalDocumentation) Description() string {
	return v.description
}

func (v *ExternalDocumentation) URL() string {
	return v.url
}

func (v *ExternalDocumentation) Extension(key string) interface{} {
	return v.extension[key]
}

func (v *Parameter) Root() *OpenAPI {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	if v.root == nil {
		return &OpenAPI{}
	}
	return v.root
}

func (v *Parameter) Name() string {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.name
}

func (v *Parameter) In() string {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.in
}

func (v *Parameter) Description() string {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.description
}

func (v *Parameter) Required() bool {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.required
}

func (v *Parameter) Deprecated() bool {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.deprecated
}

func (v *Parameter) AllowEmptyValue() bool {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.allowEmptyValue
}

func (v *Parameter) Style() string {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.style
}

func (v *Parameter) Explode() bool {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.explode
}

func (v *Parameter) AllowReserved() bool {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.allowReserved
}

func (v *Parameter) Schema() *Schema {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	if v.schema == nil {
		return &Schema{}
	}
	return v.schema
}

func (v *Parameter) Example() interface{} {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.example
}

func (v *Parameter) Examples() map[string]*Example {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.examples
}

func (v *Parameter) Content() map[string]*MediaType {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.content
}

func (v *Parameter) Extension(key string) interface{} {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.extension[key]
}

func (v *Parameter) Reference() string {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.reference
}

func (v *Parameter) Resolved() *Parameter {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	if v.resolved == nil {
		return &Parameter{}
	}
	return v.resolved
}

func (v *RequestBody) Root() *OpenAPI {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	if v.root == nil {
		return &OpenAPI{}
	}
	return v.root
}

func (v *RequestBody) Description() string {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.description
}

func (v *RequestBody) Content() map[string]*MediaType {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.content
}

func (v *RequestBody) Required() bool {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.required
}

func (v *RequestBody) Extension(key string) interface{} {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.extension[key]
}

func (v *RequestBody) Reference() string {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.reference
}

func (v *RequestBody) Resolved() *RequestBody {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	if v.resolved == nil {
		return &RequestBody{}
	}
	return v.resolved
}

func (v *MediaType) Root() *OpenAPI {
	if v.root == nil {
		return &OpenAPI{}
	}
	return v.root
}

func (v *MediaType) Schema() *Schema {
	if v.schema == nil {
		return &Schema{}
	}
	return v.schema
}

func (v *MediaType) Example() interface{} {
	return v.example
}

func (v *MediaType) Examples() map[string]*Example {
	return v.examples
}

func (v *MediaType) Encoding() map[string]*Encoding {
	return v.encoding
}

func (v *MediaType) Extension(key string) interface{} {
	return v.extension[key]
}

func (v *Encoding) Root() *OpenAPI {
	if v.root == nil {
		return &OpenAPI{}
	}
	return v.root
}

func (v *Encoding) ContentType() string {
	return v.contentType
}

func (v *Encoding) Headers() map[string]*Header {
	return v.headers
}

func (v *Encoding) Style() string {
	return v.style
}

func (v *Encoding) Explode() bool {
	return v.explode
}

func (v *Encoding) AllowReserved() bool {
	return v.allowReserved
}

func (v *Encoding) Extension(key string) interface{} {
	return v.extension[key]
}

func (v *Responses) Root() *OpenAPI {
	if v.root == nil {
		return &OpenAPI{}
	}
	return v.root
}

func (v *Responses) Get(key string) *Response {
	if val, ok := v.responses[key]; ok {
		return val
	}
	return &Response{}
}

func (v *Responses) Extension(key string) interface{} {
	return v.extension[key]
}

func (v *Response) Root() *OpenAPI {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	if v.root == nil {
		return &OpenAPI{}
	}
	return v.root
}

func (v *Response) Description() string {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.description
}

func (v *Response) Headers() map[string]*Header {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.headers
}

func (v *Response) Content() map[string]*MediaType {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.content
}

func (v *Response) Links() map[string]*Link {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.links
}

func (v *Response) Extension(key string) interface{} {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.extension[key]
}

func (v *Response) Reference() string {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.reference
}

func (v *Response) Resolved() *Response {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	if v.resolved == nil {
		return &Response{}
	}
	return v.resolved
}

func (v *Callback) Root() *OpenAPI {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	if v.root == nil {
		return &OpenAPI{}
	}
	return v.root
}

func (v *Callback) Get(key string) *PathItem {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	if val, ok := v.callback[key]; ok {
		return val
	}
	return &PathItem{}
}

func (v *Callback) Extension(key string) interface{} {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.extension[key]
}

func (v *Callback) Reference() string {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.reference
}

func (v *Callback) Resolved() *Callback {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	if v.resolved == nil {
		return &Callback{}
	}
	return v.resolved
}

func (v *Example) Root() *OpenAPI {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	if v.root == nil {
		return &OpenAPI{}
	}
	return v.root
}

func (v *Example) Summary() string {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.summary
}

func (v *Example) Description() string {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.description
}

func (v *Example) Value() interface{} {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.value
}

func (v *Example) ExternalValue() string {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.externalValue
}

func (v *Example) Extension(key string) interface{} {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.extension[key]
}

func (v *Example) Reference() string {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.reference
}

func (v *Example) Resolved() *Example {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	if v.resolved == nil {
		return &Example{}
	}
	return v.resolved
}

func (v *Link) Root() *OpenAPI {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	if v.root == nil {
		return &OpenAPI{}
	}
	return v.root
}

func (v *Link) OperationRef() string {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.operationRef
}

func (v *Link) OperationID() string {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.operationID
}

func (v *Link) Parameters() map[string]interface{} {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.parameters
}

func (v *Link) RequestBody() interface{} {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.requestBody
}

func (v *Link) Description() string {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.description
}

func (v *Link) Server() *Server {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	if v.server == nil {
		return &Server{}
	}
	return v.server
}

func (v *Link) Extension(key string) interface{} {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.extension[key]
}

func (v *Link) Reference() string {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.reference
}

func (v *Link) Resolved() *Link {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	if v.resolved == nil {
		return &Link{}
	}
	return v.resolved
}

func (v *Header) Root() *OpenAPI {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	if v.root == nil {
		return &OpenAPI{}
	}
	return v.root
}

func (v *Header) Description() string {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.description
}

func (v *Header) Required() bool {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.required
}

func (v *Header) Deprecated() bool {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.deprecated
}

func (v *Header) AllowEmptyValue() bool {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.allowEmptyValue
}

func (v *Header) Style() string {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.style
}

func (v *Header) Explode() bool {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.explode
}

func (v *Header) AllowReserved() bool {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.allowReserved
}

func (v *Header) Schema() *Schema {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	if v.schema == nil {
		return &Schema{}
	}
	return v.schema
}

func (v *Header) Example() interface{} {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.example
}

func (v *Header) Examples() map[string]*Example {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.examples
}

func (v *Header) Content() map[string]*MediaType {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.content
}

func (v *Header) Extension(key string) interface{} {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.extension[key]
}

func (v *Header) Reference() string {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.reference
}

func (v *Header) Resolved() *Header {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	if v.resolved == nil {
		return &Header{}
	}
	return v.resolved
}

func (v *Tag) Root() *OpenAPI {
	if v.root == nil {
		return &OpenAPI{}
	}
	return v.root
}

func (v *Tag) Name() string {
	return v.name
}

func (v *Tag) Description() string {
	return v.description
}

func (v *Tag) ExternalDocs() *ExternalDocumentation {
	if v.externalDocs == nil {
		return &ExternalDocumentation{}
	}
	return v.externalDocs
}

func (v *Tag) Extension(key string) interface{} {
	return v.extension[key]
}

func (v *Schema) Root() *OpenAPI {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	if v.root == nil {
		return &OpenAPI{}
	}
	return v.root
}

func (v *Schema) Title() string {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.title
}

func (v *Schema) MultipleOf() int {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.multipleOf
}

func (v *Schema) Maximum() int {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.maximum
}

func (v *Schema) ExclusiveMaximum() bool {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.exclusiveMaximum
}

func (v *Schema) Minimum() int {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.minimum
}

func (v *Schema) ExclusiveMinimum() bool {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.exclusiveMinimum
}

func (v *Schema) MaxLength() int {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.maxLength
}

func (v *Schema) MinLength() int {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.minLength
}

func (v *Schema) Pattern() string {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.pattern
}

func (v *Schema) MaxItems() int {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.maxItems
}

func (v *Schema) MinItems() int {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.minItems
}

func (v *Schema) MaxProperties() int {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.maxProperties
}

func (v *Schema) MinProperties() int {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.minProperties
}

func (v *Schema) Required() []string {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.required
}

func (v *Schema) Enum() []string {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.enum
}

func (v *Schema) Type() string {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.type_
}

func (v *Schema) AllOf() []*Schema {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.allOf
}

func (v *Schema) OneOf() []*Schema {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.oneOf
}

func (v *Schema) AnyOf() []*Schema {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.anyOf
}

func (v *Schema) Not() *Schema {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	if v.not == nil {
		return &Schema{}
	}
	return v.not
}

func (v *Schema) Items() *Schema {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	if v.items == nil {
		return &Schema{}
	}
	return v.items
}

func (v *Schema) Properties() map[string]*Schema {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.properties
}

func (v *Schema) AdditionalProperties() *Schema {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	if v.additionalProperties == nil {
		return &Schema{}
	}
	return v.additionalProperties
}

func (v *Schema) Description() string {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.description
}

func (v *Schema) Format() string {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.format
}

func (v *Schema) Default() string {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.default_
}

func (v *Schema) Nullable() bool {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.nullable
}

func (v *Schema) Discriminator() *Discriminator {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	if v.discriminator == nil {
		return &Discriminator{}
	}
	return v.discriminator
}

func (v *Schema) ReadOnly() bool {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.readOnly
}

func (v *Schema) WriteOnly() bool {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.writeOnly
}

func (v *Schema) XML() *XML {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	if v.xml == nil {
		return &XML{}
	}
	return v.xml
}

func (v *Schema) ExternalDocs() *ExternalDocumentation {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	if v.externalDocs == nil {
		return &ExternalDocumentation{}
	}
	return v.externalDocs
}

func (v *Schema) Example() interface{} {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.example
}

func (v *Schema) Deprecated() bool {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.deprecated
}

func (v *Schema) Extension(key string) interface{} {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.extension[key]
}

func (v *Schema) Reference() string {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.reference
}

func (v *Schema) Resolved() *Schema {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	if v.resolved == nil {
		return &Schema{}
	}
	return v.resolved
}

func (v *Discriminator) Root() *OpenAPI {
	if v.root == nil {
		return &OpenAPI{}
	}
	return v.root
}

func (v *Discriminator) PropertyName() string {
	return v.propertyName
}

func (v *Discriminator) Mapping() map[string]string {
	return v.mapping
}

func (v *XML) Root() *OpenAPI {
	if v.root == nil {
		return &OpenAPI{}
	}
	return v.root
}

func (v *XML) Name() string {
	return v.name
}

func (v *XML) Namespace() string {
	return v.namespace
}

func (v *XML) Prefix() string {
	return v.prefix
}

func (v *XML) Attribute() bool {
	return v.attribute
}

func (v *XML) Wrapped() bool {
	return v.wrapped
}

func (v *XML) Extension(key string) interface{} {
	return v.extension[key]
}

func (v *SecurityScheme) Root() *OpenAPI {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	if v.root == nil {
		return &OpenAPI{}
	}
	return v.root
}

func (v *SecurityScheme) Type() string {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.type_
}

func (v *SecurityScheme) Description() string {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.description
}

func (v *SecurityScheme) Name() string {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.name
}

func (v *SecurityScheme) In() string {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.in
}

func (v *SecurityScheme) Scheme() string {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.scheme
}

func (v *SecurityScheme) BearerFormat() string {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.bearerFormat
}

func (v *SecurityScheme) Flows() *OAuthFlows {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	if v.flows == nil {
		return &OAuthFlows{}
	}
	return v.flows
}

func (v *SecurityScheme) OpenIDConnectURL() string {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.openIDConnectURL
}

func (v *SecurityScheme) Extension(key string) interface{} {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.extension[key]
}

func (v *SecurityScheme) Reference() string {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	return v.reference
}

func (v *SecurityScheme) Resolved() *SecurityScheme {
	if v.reference != "" {
		resolved, err := v.resolve()
		if err != nil {
			panic(err)
		}
		v = resolved
	}
	if v.resolved == nil {
		return &SecurityScheme{}
	}
	return v.resolved
}

func (v *OAuthFlows) Root() *OpenAPI {
	if v.root == nil {
		return &OpenAPI{}
	}
	return v.root
}

func (v *OAuthFlows) Implicit() *OAuthFlow {
	if v.implicit == nil {
		return &OAuthFlow{}
	}
	return v.implicit
}

func (v *OAuthFlows) Password() *OAuthFlow {
	if v.password == nil {
		return &OAuthFlow{}
	}
	return v.password
}

func (v *OAuthFlows) ClientCredentials() *OAuthFlow {
	if v.clientCredentials == nil {
		return &OAuthFlow{}
	}
	return v.clientCredentials
}

func (v *OAuthFlows) AuthorizationCode() *OAuthFlow {
	if v.authorizationCode == nil {
		return &OAuthFlow{}
	}
	return v.authorizationCode
}

func (v *OAuthFlows) Extension(key string) interface{} {
	return v.extension[key]
}

func (v *OAuthFlow) Root() *OpenAPI {
	if v.root == nil {
		return &OpenAPI{}
	}
	return v.root
}

func (v *OAuthFlow) AuthorizationURL() string {
	return v.authorizationURL
}

func (v *OAuthFlow) TokenURL() string {
	return v.tokenURL
}

func (v *OAuthFlow) RefreshURL() string {
	return v.refreshURL
}

func (v *OAuthFlow) Scopes() map[string]string {
	return v.scopes
}

func (v *OAuthFlow) Extension(key string) interface{} {
	return v.extension[key]
}

func (v *SecurityRequirement) Root() *OpenAPI {
	if v.root == nil {
		return &OpenAPI{}
	}
	return v.root
}

func (v *SecurityRequirement) Get(key string) []string {
	if val, ok := v.securityRequirement[key]; ok {
		return val
	}
	return []string{}
}
