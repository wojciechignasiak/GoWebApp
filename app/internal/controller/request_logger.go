package controller

import (
	"log"
)

type RequestLogger interface {
	LogRequest(statusCode int, controllerAddress string)
}

type requestLogger struct{}

func NewRequestLogger() *requestLogger {
	return &requestLogger{}
}
func (rl *requestLogger) LogRequest(statusCode int, controllerAddress string) {
	log.Printf("Status: %d Address: %s", statusCode, controllerAddress)
}
