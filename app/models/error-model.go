package models

type ErrorResponse struct {
	HttpStatus int
	Code       string
	Message    string
}

func (e ErrorResponse) Error() string {
	return e.Message
}

func NewErrorResponse(httpStatus int, code string, message string) ErrorResponse {
	return ErrorResponse{HttpStatus: httpStatus, Code: code, Message: message}
}

