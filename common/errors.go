package common

import (
	"fmt"
)

type GinError struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Errors  any    `json:"errors"`
}

// NewGinError Creates a new GinError with the given status, message and detail.
func NewGinError(status string, message string, detail any) *GinError {
	return &GinError{
		Status:  status,
		Message: message,
		Errors:  detail,
	}
}

func (e *GinError) Error() string {
	return fmt.Sprintf("%d: %s", e.Status, e.Message)
}
