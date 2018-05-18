package openapi

// codebeat:disable[TOO_MANY_IVARS]

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

// Validate the values of Link object.
func (link Link) Validate() error {
	validaters := []validater{}
	for _, i := range link.Parameters {
		if v, ok := i.(validater); ok {
			validaters = append(validaters, v)
		}
	}
	if v, ok := link.RequestBody.(validater); ok {
		validaters = append(validaters, v)
	}
	validaters = append(validaters, link.Server)
	return validateAll(validaters)
}
