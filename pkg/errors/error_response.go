package errors

type ErrorMessage struct {
	Message string
}

func BadRequest(msg string) *ErrorMessage {
	return &ErrorMessage{
		Message: msg,
	}
}

func NotFound(msg string) *ErrorMessage {
	return &ErrorMessage{
		Message: msg,
	}
}

func MethodNotAllowed(msg string) *ErrorMessage {
	return &ErrorMessage{
		Message: msg,
	}
}

func UnprocessableEntity(msg string) *ErrorMessage {
	return &ErrorMessage{
		Message: msg,
	}
}

func Conflict(msg string) *ErrorMessage {
	return &ErrorMessage{
		Message: msg,
	}
}

func Unauthorized(msg string) *ErrorMessage {
	return &ErrorMessage{
		Message: msg,
	}
}

func Forbidden(msg string) *ErrorMessage {
	return &ErrorMessage{
		Message: msg,
	}
}

func InternalServerError(msg string) *ErrorMessage {
	return &ErrorMessage{
		Message: msg,
	}
}
