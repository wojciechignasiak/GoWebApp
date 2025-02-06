package servicecomponent

import (
	apperror "app/internal/app_error"

	"github.com/google/uuid"
)

type UuidGenerator interface {
	GenerateUuid() (*uuid.UUID, *apperror.AppError)
}

type uuidGenerator struct{}

func NewUuidGenerator() UuidGenerator {
	return &uuidGenerator{}
}

func (ug *uuidGenerator) GenerateUuid() (*uuid.UUID, *apperror.AppError) {
	newUUID, err := uuid.NewRandom()
	if err != nil {
		generationError := apperror.AppError{
			StatusCode:      500,
			Message:         "Error occured while generating new UUID",
			StructAndMethod: "UuidGenerator.GenerateUUID()",
			Argument:        nil,
			ChildAppError:   nil,
			ChildError:      &err,
		}
		return nil, &generationError
	}
	return &newUUID, nil
}
