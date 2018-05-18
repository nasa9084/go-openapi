package openapi

// codebeat:disable[TOO_MANY_IVARS]

// SecurityRequirement Object
type SecurityRequirement []map[string][]string

// Validate the values of SecurityRequirement object.
// TODO: implement varidation
func (secReq SecurityRequirement) Validate() error {
	return nil
}
