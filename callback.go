package openapi

// codebeat:disable[TOO_MANY_IVARS]

// Callback Object
type Callback map[string]*PathItem

// Validate the values of Callback object.
func (callback Callback) Validate() error {
	for _, pathItem := range callback {
		if err := pathItem.Validate(); err != nil {
			return err
		}
	}
	return nil
}
