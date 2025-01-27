package controller

import (
	apperror "app/internal/app_error"
	controllercomponent "app/internal/controller_component"
	"app/internal/logs"
	"app/internal/model"
	"app/internal/service"
	"encoding/json"
	"net/http"
)

type AuthController struct {
	authService         service.AuthService
	registrationService service.RegistrationService
	responseHandler     controllercomponent.ResponseHandler
	logger              logs.Logger
}

func NewAuthController(authService service.AuthService, registrationService service.RegistrationService, responseHandler controllercomponent.ResponseHandler, logger logs.Logger) *AuthController {
	return &AuthController{
		authService:         authService,
		registrationService: registrationService,
		responseHandler:     responseHandler,
		logger:              logger,
	}
}

func (ac *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var newUser model.NewUser
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&newUser)
	if err != nil {
		invalidJson := apperror.AppError{
			StatusCode:      400,
			Message:         "invalid JSON input",
			StructAndMethod: "AuthController.RegisterUser()",
			Argument:        nil,
			ChildAppError:   nil,
			ChildError:      nil,
		}
		ac.responseHandler.HandleError(w, &invalidJson)
		return
	}

	ctx := r.Context()

	serviceError := ac.registrationService.Register(ctx, newUser)

	if serviceError != nil {
		ac.logger.LogRequest(serviceError.StatusCode, "/auth/register")
		if serviceError.StatusCode == http.StatusInternalServerError {
			ac.logger.LogAppError(serviceError)
		}
		ac.responseHandler.HandleError(w, serviceError)
		return
	}
	ac.logger.LogRequest(http.StatusCreated, "/auth/register")
	ac.responseHandler.SendResponse(w, http.StatusCreated, "registered successfully")
}

func (ac *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var credentials model.Credentials
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&credentials)
	if err != nil {
		invalidJson := apperror.AppError{
			StatusCode:      400,
			Message:         "invalid JSON input",
			StructAndMethod: "AuthController.Login()",
			Argument:        nil,
			ChildAppError:   nil,
			ChildError:      nil,
		}
		ac.responseHandler.HandleError(w, &invalidJson)
		return
	}

	ctx := r.Context()

	sessionId, serviceError := ac.authService.Login(ctx, credentials)

	if serviceError != nil {
		ac.logger.LogRequest(serviceError.StatusCode, "/auth/login")
		if serviceError.StatusCode == http.StatusInternalServerError {
			ac.logger.LogAppError(serviceError)
		}
		ac.responseHandler.HandleError(w, serviceError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "sessionId",
		Value:    sessionId.String(),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
	ac.logger.LogRequest(http.StatusOK, "/auth/login")
	ac.responseHandler.SendResponse(w, http.StatusOK, "Login successful")
}
