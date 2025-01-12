package controllercomponent

import (
	apperror "app/internal/app_error"
	"encoding/json"
	"net/http"
)

type ResponseHandler interface {
	HandleError(w http.ResponseWriter, err *apperror.AppError)
	SendResponse(w http.ResponseWriter, statusCode int, message string)
}

type responseHandler struct{}

func NewResponseHandler() *responseHandler {
	return &responseHandler{}
}

func (rh *responseHandler) HandleError(w http.ResponseWriter, err *apperror.AppError) {
	w.Header().Set("Content-Type", "application/json")

	response := map[string]string{"message": err.Message}
	if err.StatusCode == http.StatusInternalServerError {
		response["message"] = "internal server error"
	}

	jsonResponse, jsonErr := json.Marshal(response)
	if jsonErr != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(err.StatusCode)
	w.Write(jsonResponse)
}

func (rh *responseHandler) SendResponse(w http.ResponseWriter, statusCode int, message string) {
	response := map[string]string{"message": message}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(jsonResponse)
}
