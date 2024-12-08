package model

import "github.com/google/uuid"

type AccountConfirmation struct {
	UserId           uuid.UUID `json:"user_id" db:"user_id"`
	ConfirmationCode uuid.UUID `json:"confirmation_code" db:"confirmation_code"`
}
