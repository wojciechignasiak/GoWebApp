package server

import (
	"app/internal/controller"
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	address        string
	port           int
	userController *controller.UserController
}

func NewServer(address string, port int, userController *controller.UserController) *http.Server {
	NewServer := &Server{
		address:        address,
		port:           port,
		userController: userController,
	}
	server := http.Server{
		Addr:         fmt.Sprintf("%s:%d", NewServer.address, NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return &server
}
