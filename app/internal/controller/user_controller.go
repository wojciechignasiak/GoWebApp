package controller

import (
	"app/internal/logs"
	"app/internal/model"
	"app/internal/service"
	"encoding/json"
	"log"
	"net/http"
)

type UserController struct {
	logger      logs.Logger
	userService service.UserService
}

func NewUserController(userService service.UserService, logger logs.Logger) *UserController {
	return &UserController{
		logger:      logger,
		userService: userService,
	}
}

func (uc *UserController) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var newUser model.CreateUser
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&newUser)
	if err != nil {
		http.Error(w, "message: invalid JSON input", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	serviceError := uc.userService.RegisterUser(ctx, newUser)

	if serviceError != nil {

		switch serviceError.StatusCode {
		case 400:
			uc.logger.LogRequest(400, "/user/register")
			http.Error(w, "message: "+serviceError.Message, http.StatusBadRequest)
		case 409:
			uc.logger.LogRequest(409, "/user/register")
			http.Error(w, "message: "+serviceError.Message, http.StatusConflict)
		case 500:
			uc.logger.LogRequest(500, "/user/register")
			uc.logger.LogAppError(serviceError)
			http.Error(w, "message: Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	response := map[string]string{"message": "user registered successfully"}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResponse)
	uc.logger.LogRequest(200, "/user/register")
}
