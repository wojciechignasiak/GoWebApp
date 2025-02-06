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

type saltGeneratorAndPasswordHasher struct {
	argon2Func   func(password []byte, salt []byte, time uint32, memory uint32, threads uint8, keyLen uint32) []byte
	randReadFunc func(b []byte) (n int, err error)
	timeCost     uint32
	memoryCost   uint32
	threads      uint8
	keyLength    uint32
}

func NewSaltGeneratorAndPasswordHasher() SaltGeneratorAndPasswordHasher {
	return &saltGeneratorAndPasswordHasher{
		argon2Func:   argon2.IDKey,
		randReadFunc: rand.Read,
		timeCost:     uint32(3),
		memoryCost:   uint32(64 * 1024),
		threads:      uint8(4),
		keyLength:    uint32(32),
	}
}

func (sgaph *saltGeneratorAndPasswordHasher) HashPassword(password string, salt []byte) *[]byte {
	hash := sgaph.argon2Func([]byte(password), salt, sgaph.timeCost, sgaph.memoryCost, sgaph.threads, sgaph.keyLength)
	return &hash
}

func (sgaph *saltGeneratorAndPasswordHasher) GenerateSalt() (*[]byte, *apperror.AppError) {
	salt := make([]byte, 16)
	_, err := sgaph.randReadFunc(salt)
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
