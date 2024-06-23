package dto

type ErrorResponse struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
}

func NewInternalServerError(err error) *ErrorResponse {
	return &ErrorResponse{
		Message: "Internal Server Error",
		Code:    500,
		Data:    err.Error(),
	}
}

func NewValidationError(err error) *ErrorResponse {
	return &ErrorResponse{
		Message: "Validation Error",
		Code:    422,
		Data:    err.Error(),
	}
}
