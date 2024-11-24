package apperror

type AppError struct {
	StatusCode    int
	Message       string
	ChildAppError *AppError
	ChildError    *error
	Logging       bool
}

// func (e *CustomError) Error(message string, statusCode int, childError *CustomError, logging bool) *CustomError {
// 	return &CustomError{
// 		Message:    message,
// 		StatusCode: statusCode,
// 		ChildError: childError,
// 		Logging:    logging,
// 	}
// }
