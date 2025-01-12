package logs

import (
	apperror "app/internal/app_error"
	"log"
)

type Logger interface {
	LogRequest(statusCode int, controllerAddress string)
	LogAppError(err *apperror.AppError)
}

type logger struct{}

func NewLogger() *logger {
	return &logger{}
}
func (l *logger) LogRequest(statusCode int, controllerAddress string) {
	log.Printf("Status Code: %d Address: %s", statusCode, controllerAddress)
}
func (l *logger) LogAppError(err *apperror.AppError) {
	if err == nil {
		return
	}

	log.Printf("Status Code: %d, Message: %s, Struct and Method: %s", err.StatusCode, err.Message, err.StructAndMethod)

	if err.Argument != nil {
		log.Printf("Argument: %s", *err.Argument)
	}

	if err.ChildError != nil && *err.ChildError != nil {
		log.Printf("Child Error: %s", (*err.ChildError).Error())
	}

	if err.ChildAppError != nil {
		log.Println("Nested Error:")
		l.LogAppError(err.ChildAppError)
	}
}
