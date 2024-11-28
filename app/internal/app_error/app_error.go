package apperror

type AppError struct {
	StatusCode      int
	Message         string
	StructAndMethod string
	Argument        *string
	ChildAppError   *AppError
	ChildError      *error
}
