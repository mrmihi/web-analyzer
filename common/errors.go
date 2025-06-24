package common

import (
	"fmt"
)

// TODO: Make errors of my own and use standards

type GinError struct {
	Code    int
	Message string
	Detail  any `json:"detail"`
}

// NewGinError Creates a new GinError with the given status code, message and detail.
func NewGinError(code int, message string, detail any) *GinError {
	return &GinError{
		Code:    code,
		Message: message,
		Detail:  detail,
	}
}

// Error implements the error interface for GinError. Do not rename this method.
func (e *GinError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

// Status returns the HTTP status code.
func (e *GinError) Status() int {
	return e.Code
}
