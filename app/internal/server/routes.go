package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/auth/register", s.authController.Register).Methods("POST")
	r.HandleFunc("/auth/login", s.authController.Login).Methods("POST")
	r.HandleFunc("/auth/logout", s.authController.Logout).Methods("DELETE")
	r.HandleFunc("/user/confirm-account/{confirmationCode}/{securityCode}", s.userController.ConfirmAccount).Methods("PUT")
	return r
}
