package openapi

import "strconv"

// SuccessResponse returns a success response object.
// If there are 2 or more success responses (like created and ok),
// it's not sure which is returned.
// If only match the default response or 2XX response, returned status code will be 0.
func (op *Operation) SuccessResponse() (*Response, int, bool) {
	if op == nil || op.Responses == nil {
		return nil, -1, false
	}
	var defaultResponse *Response
	for statusStr, resp := range op.Responses {
		switch statusStr {
		case "default":
			defaultResponse = resp
		case "2XX":
			defaultResponse = resp
		case "1XX", "3XX", "4XX", "5XX":
			continue
		}
		statusInt, err := strconv.Atoi(statusStr)
		if err != nil {
			continue
		}
		if statusInt/100 == 2 {
			return resp, statusInt, true
		}
	}
	return defaultResponse, 0, (defaultResponse != nil)
}
