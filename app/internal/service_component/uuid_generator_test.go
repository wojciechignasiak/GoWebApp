package servicecomponent

import (
	apperror "app/internal/app_error"
	"errors"
	"testing"

	"github.com/google/uuid"
)

// Test poprawnego generowania UUID
func TestGenerateUuid_Success(t *testing.T) {
	mockUUID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")

	generator := &uuidGenerator{
		uuidFunc: func() (uuid.UUID, error) {
			return mockUUID, nil
		},
	}

	result, err := generator.GenerateUuid()

	if err != nil {
		t.Fatalf("Oczekiwano braku błędu, otrzymano: %v", err)
	}
	if result == nil {
		t.Fatal("Oczekiwano UUID, otrzymano nil")
	}
	if *result != mockUUID {
		t.Errorf("Oczekiwano UUID %v, otrzymano %v", mockUUID, *result)
	}
}

func TestGenerateUuid_Error(t *testing.T) {
	expectedErr := errors.New("mock error")

	generator := &uuidGenerator{
		uuidFunc: func() (uuid.UUID, error) {
			return uuid.UUID{}, expectedErr
		},
	}

	result, generationError := generator.GenerateUuid()

	if result != nil {
		t.Errorf("Expected nil for UUID, got: %v", result)
	}
	if generationError == nil {
		t.Fatal("Expected error, got nil")
	}
	if _, ok := interface{}(generationError).(*apperror.AppError); !ok {
		t.Errorf("Expected *app_error.AppError, got: %T", generationError)
	}
	if generationError.ChildError == nil || *generationError.ChildError != expectedErr {
		t.Errorf("Expected error: %v, got: %v", expectedErr, generationError.ChildError)
	}
}
