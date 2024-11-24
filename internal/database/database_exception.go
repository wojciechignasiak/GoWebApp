package database

import "fmt"

type HTTPError struct {
	StatusCode int
	Message    string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("%d: %s", e.StatusCode, e.Message)
}

func BadRequest(message string) *HTTPError {
	return &HTTPError{
		StatusCode: 400,
		Message:    message,
	}
}

func NotFound(message string) *HTTPError {
	return &HTTPError{
		StatusCode: 404,
		Message:    message,
	}
}

func InternalServerError(message string) *HTTPError {
	return &HTTPError{
		StatusCode: 500,
		Message:    message,
	}
}
