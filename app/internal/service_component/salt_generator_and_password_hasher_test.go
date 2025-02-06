package servicecomponent

import (
	"errors"
	"testing"

	apperror "app/internal/app_error"

	"golang.org/x/crypto/argon2"
)

func TestGenerateSalt_Success(t *testing.T) {
	generator := &saltGeneratorAndPasswordHasher{
		randReadFunc: func(b []byte) (int, error) {
			for i := range b {
				b[i] = byte(i)
			}
			return len(b), nil
		},
	}

	salt, err := generator.GenerateSalt()

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if salt == nil {
		t.Fatal("Expected generated salt, got nil")
	}
	if len(*salt) != 16 {
		t.Errorf("Expected salt length of 16, got %d", len(*salt))
	}
}

func TestGenerateSalt_Error(t *testing.T) {
	expectedErr := errors.New("mock read error")

	generator := &saltGeneratorAndPasswordHasher{
		randReadFunc: func(b []byte) (int, error) {
			return 0, expectedErr
		},
	}

	salt, err := generator.GenerateSalt()

	if salt != nil {
		t.Errorf("Expected nil for salt, got: %v", salt)
	}
	if err == nil {
		t.Fatal("Expected an error, got nil")
	}
	if _, ok := interface{}(err).(*apperror.AppError); !ok {
		t.Errorf("Expected *app_error.AppError, got: %T", err)
	}
	if err.ChildError == nil || *err.ChildError != expectedErr {
		t.Errorf("Expected error %v, got %v", expectedErr, err.ChildError)
	}
}

func TestHashPassword(t *testing.T) {
	password := "securepassword"
	salt := []byte("1234567890123456")
	expectedHash := argon2.IDKey([]byte(password), salt, 3, 64*1024, 4, 32)

	generator := &saltGeneratorAndPasswordHasher{
		argon2Func: func(password []byte, salt []byte, time uint32, memory uint32, threads uint8, keyLen uint32) []byte {
			return expectedHash
		},
	}

	hash := generator.HashPassword(password, salt)

	if hash == nil {
		t.Fatal("Expected a result, got nil")
	}
	if len(*hash) != 32 {
		t.Errorf("Expected hash length of 32, got %d", len(*hash))
	}
	if string(*hash) != string(expectedHash) {
		t.Errorf("Hash does not match the expected result")
	}
}
