package servicecomponent

import (
	apperror "app/internal/app_error"

	"github.com/google/uuid"
)

type UuidGenerator interface {
	GenerateUuid() (*uuid.UUID, *apperror.AppError)
}

type uuidGenerator struct {
	uuidFunc func() (uuid.UUID, error)
}

func NewUuidGenerator() UuidGenerator {
	return &uuidGenerator{
		uuidFunc: uuid.NewRandom,
	}
}

func (ug *uuidGenerator) GenerateUuid() (*uuid.UUID, *apperror.AppError) {
	newUUID, err := ug.uuidFunc()
	if err != nil {
		generationError := apperror.AppError{
			StatusCode:      500,
			Message:         "Error occurred while generating new UUID",
			StructAndMethod: "UuidGenerator.GenerateUUID()",
			Argument:        nil,
			ChildAppError:   nil,
			ChildError:      &err,
		}
		return nil, &generationError
	}
	return &newUUID, nil
}
