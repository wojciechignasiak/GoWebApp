package apperror

type AppError struct {
	StatusCode      int
	Message         string
	StructAndMethod string
	Argument        *string
	ChildAppError   *AppError
	ChildError      *error
}

func (e *AppError) Error() string {
	if e == nil {
		return "<nil>"
	}
	return e.Message
}

func (e *AppError) FullError() *AppError {
	return e
}
