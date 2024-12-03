package model

import (
	apperror "app/internal/app_error"
	"fmt"
	"unicode/utf8"

	"github.com/google/uuid"
)

type Chat struct {
	Id   uuid.UUID `json:"id" db:"id"`
	Name string    `json:"name" db:"name"`
}

func (c *Chat) ValidateCreateChat() *apperror.AppError {
	err := validateName(c.Name)

	if err != nil {
		return err
	}
	return nil
}

func validateName(name string) *apperror.AppError {
	numberOfCharacters := utf8.RuneCountInString(name)
	if numberOfCharacters == 0 {
		args := fmt.Sprintf("name: %s", name)
		err := apperror.AppError{
			StatusCode:      400,
			Message:         "Chat name can't be empty.",
			StructAndMethod: "Chat.validateName()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      nil,
		}
		return &err
	} else if numberOfCharacters > 100 {
		args := fmt.Sprintf("name: %s", name)
		err := apperror.AppError{
			StatusCode:      400,
			Message:         "Chat name must be no more than 100 characters long.",
			StructAndMethod: "Chat.validateName()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      nil,
		}
		return &err
	} else {
		return nil
	}
}
