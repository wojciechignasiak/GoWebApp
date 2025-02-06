package servicecomponent

import (
	apperror "app/internal/app_error"
	"crypto/rand"

	"golang.org/x/crypto/argon2"
)

type SaltGeneratorAndPasswordHasher interface {
	HashPassword(password string, salt []byte) *[]byte
	GenerateSalt() (*[]byte, *apperror.AppError)
}

type saltGeneratorAndPasswordHasher struct{}

func NewSaltGeneratorAndPasswordHasher() SaltGeneratorAndPasswordHasher {
	return &saltGeneratorAndPasswordHasher{}
}

func (sgaph *saltGeneratorAndPasswordHasher) HashPassword(password string, salt []byte) *[]byte {
	timeCost := uint32(3)
	memoryCost := uint32(64 * 1024)
	threads := uint8(4)
	keyLength := uint32(32)
	hash := argon2.IDKey([]byte(password), salt, timeCost, memoryCost, threads, keyLength)
	return &hash
}

func (sgaph *saltGeneratorAndPasswordHasher) GenerateSalt() (*[]byte, *apperror.AppError) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		generationError := apperror.AppError{
			StatusCode:      500,
			Message:         "Error occured while generating salt",
			StructAndMethod: "SaltGeneratorAndPasswordHasher.GenerateSalt()",
			Argument:        nil,
			ChildAppError:   nil,
			ChildError:      &err,
		}
		return nil, &generationError
	}
	return &salt, nil
}
