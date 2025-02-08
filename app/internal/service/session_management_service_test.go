package service

import (
	apperror "app/internal/app_error"
	"app/internal/model"
	"testing"

	"github.com/google/uuid"
)

type mockUuidGenerator struct {
	generatedUuid uuid.UUID
	customError   *apperror.AppError
}

func (m *mockUuidGenerator) GenerateUuid() (*uuid.UUID, *apperror.AppError) {
	if m.customError != nil {
		return nil, m.customError
	}
	return &m.generatedUuid, nil
}

func TestCreateSessionSuccess(t *testing.T) {
	mockUUID := uuid.New()
	ug := &mockUuidGenerator{generatedUuid: mockUUID}
	sms := NewSessionManagementService(ug)
	userSession := model.UserSession{
		Id:       uuid.New(),
		Email:    "john_doe@example.com",
		Username: "john_doe",
	}

	sessionId, err := sms.CreateSession(userSession)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if sessionId == nil || *sessionId != mockUUID {
		t.Errorf("expected sessionId %v, got %v", mockUUID, sessionId)
	}
}

func TestCreateSessionFailure(t *testing.T) {

	generationError := &apperror.AppError{
		StatusCode:      500,
		Message:         "UUID generation failed",
		StructAndMethod: "MockMethod",
		Argument:        nil,
		ChildAppError:   nil,
		ChildError:      nil,
	}

	ug := &mockUuidGenerator{
		generatedUuid: uuid.New(),
		customError:   generationError,
	}
	sms := NewSessionManagementService(ug)

	userSession := model.UserSession{
		Id:       uuid.New(),
		Email:    "john_doe@example.com",
		Username: "john_doe",
	}

	sessionId, createSessionError := sms.CreateSession(userSession)

	if sessionId != nil {
		t.Errorf("expected nil sessionId, got %v", sessionId)
	}

	if createSessionError == nil || createSessionError.Message != "UUID generation failed" {
		t.Errorf("expected error message %q, got %q", "UUID generation failed", createSessionError.Message)
	}
}

func TestGetUserSession(t *testing.T) {
	mockUUID := uuid.New()

	ug := &mockUuidGenerator{
		generatedUuid: mockUUID,
		customError:   nil,
	}

	sms := NewSessionManagementService(ug)

	userSession := model.UserSession{
		Id:       uuid.New(),
		Email:    "john_doe@example.com",
		Username: "john_doe",
	}

	sms.CreateSession(userSession)

	retrievedSession := sms.GetUserSession(mockUUID)

	if retrievedSession.Id != userSession.Id {
		t.Errorf("expected: %v, got %v", userSession.Id, retrievedSession.Id)
	}
}

func TestGetUserSessionNotFound(t *testing.T) {
	sms := NewSessionManagementService(&mockUuidGenerator{})
	mockUUID := uuid.New()

	retrievedSession := sms.GetUserSession(mockUUID)

	if retrievedSession != nil {
		t.Errorf("expected: nil session, got %v", retrievedSession)
	}
}

func TestDeleteSession(t *testing.T) {
	mockUUID := uuid.New()

	ug := &mockUuidGenerator{
		generatedUuid: mockUUID,
		customError:   nil,
	}

	sms := NewSessionManagementService(ug)

	userSession := model.UserSession{
		Id:       uuid.New(),
		Email:    "john_doe@example.com",
		Username: "john_doe",
	}

	sms.CreateSession(userSession)

	sms.DeleteSession(mockUUID)
	foundSession := sms.GetUserSession(mockUUID)
	if foundSession != nil {
		t.Errorf("expected: nil, got: %v", foundSession)
	}
}
