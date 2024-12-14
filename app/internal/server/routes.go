package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/user/register", s.userController.RegisterUser).Methods("POST")
	r.HandleFunc("/user/confirm-account/{confirmationCode}/{securityCode}", s.userController.ConfirmAccount).Methods("PUT")
	return r
}
