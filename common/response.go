package common

type GinResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// NewGinResponse creates a new GinResponse with the given status, message, and detail.
func NewGinResponse(status int, message string, data any) *GinResponse {
	return &GinResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
}
