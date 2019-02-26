package openapi

import (
	"strings"
)

// codebeat:disable[TOO_MANY_IVARS]

// Paths Object
type Paths map[string]*PathItem

// Validate the values of Paths object.
func (paths Paths) Validate() error {
	for path, pathItem := range paths {
		if !strings.HasPrefix(path, "/") {
			return ErrPathFormat
		}
		if err := pathItem.Validate(); err != nil {
			return err
		}
	}
	if hasDuplicatedOperationID(paths) {
		return ErrOperationIDDuplicated
	}
	return nil
}

func hasDuplicatedOperationID(paths Paths) bool {
	opIDs := map[string]struct{}{}
	for _, pathItem := range paths {
		for _, op := range pathItem.Operations() {
			if _, ok := opIDs[op.OperationID]; ok {
				return true
			}
			opIDs[op.OperationID] = struct{}{}
		}
	}

	return false
}

// GetOperationByID returns an operation by operationId.
// If the paths object has two or more operations which matches
// given operationId, this function returns the operation
// matched first. So you should call Validate() before using this
// function.
func (paths Paths) GetOperationByID(operationID string) *Operation {
	for _, pathItem := range paths {
		for _, op := range pathItem.Operations() {
			if op.OperationID == operationID {
				return op
			}
		}
	}
	return nil
}
