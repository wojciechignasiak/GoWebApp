package server

import (
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {

	mux := http.NewServeMux()
	mux.HandleFunc("/user/create", s.userController.CreateUser)
	return mux
}
