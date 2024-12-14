package servicecomponent

import (
	apperror "app/internal/app_error"
	"crypto/rand"
	"fmt"
	mathrand "math/rand"

	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
)

type CommonTools interface {
	GenerateUUID() (*uuid.UUID, *apperror.AppError)
	GenerateSalt(length int) (*[]byte, *apperror.AppError)
	HashPassword(password string, salt []byte) []byte
	Generate6DigitCode() string
}

type commonTools struct{}

func NewCommonTools() *commonTools {
	return &commonTools{}
}

func (ct *commonTools) GenerateUUID() (*uuid.UUID, *apperror.AppError) {
	newUUID, err := uuid.NewV7()
	if err != nil {
		generationError := apperror.AppError{
			StatusCode:      500,
			Message:         "Error occured while generating new UUID",
			StructAndMethod: "commonTools.GenerateUUID()",
			Argument:        nil,
			ChildAppError:   nil,
			ChildError:      &err,
		}
		return nil, &generationError
	}
	return &newUUID, nil
}

func (ct *commonTools) GenerateSalt(length int) (*[]byte, *apperror.AppError) {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		args := fmt.Sprintf("length: %d", length)
		generationError := apperror.AppError{
			StatusCode:      500,
			Message:         "Error occured while generating salt",
			StructAndMethod: "commonTools.GenerateSalt()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      &err,
		}
		return nil, &generationError
	}
	return &salt, nil
}

func (ct *commonTools) HashPassword(password string, salt []byte) []byte {
	timeCost := uint32(3)
	memoryCost := uint32(64 * 1024)
	threads := uint8(4)
	keyLength := uint32(32)

	hash := argon2.IDKey([]byte(password), salt, timeCost, memoryCost, threads, keyLength)

	return hash
}

func (ct *commonTools) Generate6DigitCode() string {
	return fmt.Sprintf("%06d", mathrand.Intn(1000000))
}
