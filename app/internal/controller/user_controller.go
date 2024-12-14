package controller

import (
	controllercomponent "app/internal/controller_component"
	"app/internal/logs"
	"app/internal/model"
	"app/internal/service"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type UserController struct {
	userService     service.UserService
	responseHandler controllercomponent.ResponseHandler
	logger          logs.Logger
}

func NewUserController(userService service.UserService, responseHandler controllercomponent.ResponseHandler, logger logs.Logger) *UserController {
	return &UserController{
		userService:     userService,
		responseHandler: responseHandler,
		logger:          logger,
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
		uc.logger.LogRequest(serviceError.StatusCode, "/user/register")
		if serviceError.StatusCode == http.StatusInternalServerError {
			uc.logger.LogAppError(serviceError)
		}
		uc.responseHandler.HandleError(w, serviceError)
		return
	}
	uc.logger.LogRequest(http.StatusCreated, "/user/register")
	uc.responseHandler.SendResponse(w, http.StatusCreated, "user registered successfully")
}

func (uc *UserController) ConfirmAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	confirmationCode := vars["confirmationCode"]
	securityCode := vars["securityCode"]
	confirmationCodeUUID, err := uuid.Parse(confirmationCode)
	if err != nil {
		http.Error(w, "message: invalid confirmation code format", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	confirmAccount := model.ConfirmAccount{
		ConfirmationCode: confirmationCodeUUID,
		SecurityCode:     securityCode,
	}
	serviceError := uc.userService.ConfirmAccount(ctx, confirmAccount)

	if serviceError != nil {
		uc.logger.LogRequest(serviceError.StatusCode, "/user/confirm-account")
		if serviceError.StatusCode == http.StatusInternalServerError {
			uc.logger.LogAppError(serviceError)
		}
		uc.responseHandler.HandleError(w, serviceError)
		return
	}

	uc.logger.LogRequest(http.StatusOK, "/user/confirm-account")
	uc.responseHandler.SendResponse(w, http.StatusOK, "account confirmed successfully")
}
