package service

import (
	apperror "app/internal/app_error"
	"app/internal/model"
	servicecomponent "app/internal/service_component"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type SessionManagementService interface {
	CreateSession(userSession model.UserSession) (*uuid.UUID, *apperror.AppError)
	GetUserSession(sessionId uuid.UUID) (*model.UserSession, bool)
	DeleteSession(sessionId uuid.UUID)
}

type sessionManagementService struct {
	sessions map[uuid.UUID]model.UserSession
	mu       sync.RWMutex
	ug       servicecomponent.UuidGenerator
}

func NewSessionManagementService(ug servicecomponent.UuidGenerator) SessionManagementService {
	return &sessionManagementService{
		sessions: make(map[uuid.UUID]model.UserSession),
		ug:       ug,
	}
}

func (sms *sessionManagementService) CreateSession(userSession model.UserSession) (*uuid.UUID, *apperror.AppError) {
	sessionId, generationError := sms.ug.GenerateUuid()

	if generationError != nil {
		args := fmt.Sprintf("userSession: %v", userSession)
		serviceError := apperror.AppError{
			StatusCode:      generationError.StatusCode,
			Message:         generationError.Message,
			StructAndMethod: "SessionManagementService.CreateSession()",
			Argument:        &args,
			ChildAppError:   generationError,
			ChildError:      nil,
		}
		return nil, &serviceError
	}

	sms.saveSession(*sessionId, userSession)

	return sessionId, nil
}

func (sms *sessionManagementService) saveSession(sessionId uuid.UUID, userSession model.UserSession) {
	sms.mu.Lock()
	defer sms.mu.Unlock()
	sms.sessions[sessionId] = userSession
}

func (sms *sessionManagementService) GetUserSession(sessionId uuid.UUID) (*model.UserSession, bool) {
	sms.mu.RLock()
	defer sms.mu.RUnlock()
	userSession, exists := sms.sessions[sessionId]
	return &userSession, exists
}

func (sms *sessionManagementService) DeleteSession(sessionId uuid.UUID) {
	sms.mu.Lock()
	defer sms.mu.Unlock()
	delete(sms.sessions, sessionId)
}
