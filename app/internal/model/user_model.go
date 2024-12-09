package model

import (
	"time"

	"github.com/google/uuid"
)

type CreateUser struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	ConfirmEmail    string `json:"confirm_email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type User struct {
	Id                 uuid.UUID  `json:"id" db:"id"`
	Username           string     `json:"username" db:"username"`
	Email              string     `json:"email" db:"email"`
	Password           []byte     `json:"password" db:"password"`
	Salt               []byte     `json:"salt" db:"salt"`
	RegistrationDate   *time.Time `json:"registration_date" db:"registration_date"`
	IsAccountConfirmed bool       `json:"is_account_confirmed" db:"is_account_confirmed"`
	IsAccountDeleted   bool       `json:"is_account_deleted" db:"is_account_deleted"`
}

type ReturnUser struct {
	Id                 uuid.UUID `json:"id"`
	Username           string    `json:"username"`
	Email              string    `json:"email"`
	RegistrationDate   time.Time `json:"registration_date" db:"registration_date"`
	IsAccountConfirmed bool      `json:"is_account_confirmed" db:"is_account_confirmed"`
	IsAccountDeleted   bool      `json:"is_account_deleted" db:"is_account_deleted"`
}
