package model

type (
	Response struct {
		Result interface{}   `json:"result,omitempty"`
		Error  ErrorResponse `json:"error,omitempty"`
	}

	ErrorResponse struct {
		Message string `json:"message"`
		Detail  string `json:"detail"`
	}
)

func NewErrorResponse(message string, err error) Response {
	return Response{
		Error: ErrorResponse{
			Message: message,
			Detail:  err.Error(),
		},
	}
}
