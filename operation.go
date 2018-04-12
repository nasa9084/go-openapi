package openapi

import "strconv"

// SuccessResponse returns a success response object.
// If there are 2 or more success responses (like created and ok),
// it's not sure which is returned.
func (op *Operation) SuccessResponse() (*Response, int, bool) {
	if op.Responses == nil {
		return nil, 0, false
	}

	for statusStr, resp := range op.Responses {
		statusInt, err := strconv.Atoi(statusStr)
		if err != nil {
			continue
		}
		if statusInt/100 == 2 {
			return resp, statusInt, true
		}
	}
	return op.Responses["Default"], 0, false
}
