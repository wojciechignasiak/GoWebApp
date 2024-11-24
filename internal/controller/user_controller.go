package controller

import (
	"app/internal/model"
	"app/internal/service"
	"encoding/json"
	"log"
	"net/http"
)

type UserController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (uc *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser model.CreateUser
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&newUser)
	if err != nil {
		http.Error(w, "message: invalid JSON input", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	serviceError := uc.userService.CreateUser(ctx, newUser)

	if serviceError != nil {
		if serviceError.StatusCode == 400 {
			http.Error(w, "message: "+serviceError.Message, http.StatusBadRequest)
		}
		if serviceError.StatusCode == 409 {
			http.Error(w, "message: "+serviceError.Message, http.StatusConflict)
		}
		if serviceError.StatusCode == 500 {
			http.Error(w, "message: Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	response := map[string]string{"message": "user created successfully"}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResponse)
}
