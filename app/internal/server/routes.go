package server

import (
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {

	mux := http.NewServeMux()
	mux.HandleFunc("/user/register", s.userController.RegisterUser)
	return mux
}
