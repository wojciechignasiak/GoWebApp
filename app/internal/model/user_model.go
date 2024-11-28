package model

import (
	apperror "app/internal/app_error"
	"fmt"
	"regexp"
	"unicode"

	"github.com/google/uuid"
)

type CreateUser struct {
	Username        string  `json:"username"`
	Email           string  `json:"email"`
	ConfirmEmail    string  `json:"confirm_email"`
	Password        string  `json:"password"`
	ConfirmPassword string  `json:"confirm_password"`
	PhoneNumber     *string `json:"phone_number"`
}

func (cu *CreateUser) ValidateCreateUser() *apperror.AppError {

	if err := validateUsername(cu.Username); err != nil {
		return err
	}

	if err := validateEmails(cu.Email, cu.ConfirmEmail); err != nil {
		return err
	}

	if err := validatePasswords(cu.Password, cu.ConfirmPassword); err != nil {
		return err
	}

	if cu.PhoneNumber != nil {
		if err := validatePhoneNumber(*cu.PhoneNumber); err != nil {
			return err
		}
	}

	return nil
}

func validateUsername(username string) *apperror.AppError {
	if len(username) < 5 || len(username) > 20 {
		args := fmt.Sprintf("username: %s", username)
		validationError := apperror.AppError{
			StatusCode:      400,
			Message:         "Username must contain between 5 and 20 characters",
			StructAndMethod: "CreateUser.validateUsername()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      nil,
		}
		return &validationError
	}
	return nil
}

func validateEmails(email, confirmEmail string) *apperror.AppError {
	if email != confirmEmail {
		args := fmt.Sprintf("confirmEmail: %s", confirmEmail)
		validationError := apperror.AppError{
			StatusCode:      400,
			Message:         "Provided emails do not match",
			StructAndMethod: "CreateUser.validateEmails()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      nil,
		}
		return &validationError
	}
	if !isValidEmail(email) {
		args := fmt.Sprintf("confirmEmail: %s", confirmEmail)
		validationError := apperror.AppError{
			StatusCode:      400,
			Message:         "Invalid email format",
			StructAndMethod: "CreateUser.validateEmails()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      nil,
		}
		return &validationError
	}
	return nil
}

func validatePasswords(password, confirmPassword string) *apperror.AppError {
	if password != confirmPassword {
		validationError := apperror.AppError{
			StatusCode:      400,
			Message:         "Provided passwords are not the same",
			StructAndMethod: "CreateUser.validatePasswords()",
			Argument:        nil,
			ChildAppError:   nil,
			ChildError:      nil,
		}
		return &validationError

	}
	if len(password) < 8 {
		validationError := apperror.AppError{
			StatusCode:      400,
			Message:         "Password must contain at least 8 characters",
			StructAndMethod: "CreateUser.validatePasswords()",
			Argument:        nil,
			ChildAppError:   nil,
			ChildError:      nil,
		}
		return &validationError
	}

	var hasDigit, hasSpecial bool
	for _, char := range password {
		if unicode.IsDigit(char) {
			hasDigit = true
		} else if unicode.IsPunct(char) || unicode.IsSymbol(char) {
			hasSpecial = true
		}
		if hasDigit && hasSpecial {
			break
		}
	}

	if !hasDigit || !hasSpecial {
		validationError := apperror.AppError{
			StatusCode:      403,
			Message:         "Password must contain at least one digit and one special character",
			StructAndMethod: "CreateUser.validatePasswords()",
			Argument:        nil,
			ChildAppError:   nil,
			ChildError:      nil,
		}
		return &validationError
	}

	return nil
}

func validatePhoneNumber(phoneNumber string) *apperror.AppError {
	if !isValidPhoneNumber(phoneNumber) {
		validationError := apperror.AppError{
			StatusCode:      403,
			Message:         "Invalid phone number format",
			StructAndMethod: "CreateUser.validatePhoneNumber()",
			Argument:        nil,
			ChildAppError:   nil,
			ChildError:      nil,
		}
		return &validationError
	}
	return nil
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func isValidPhoneNumber(phoneNumber string) bool {
	re := regexp.MustCompile(`^\+?[0-9]{7,15}$`)
	return re.MatchString(phoneNumber)
}

type User struct {
	Id          uuid.UUID `json:"id" db:"id"`
	Username    string    `json:"username" db:"username"`
	Email       string    `json:"email" db:"email"`
	Password    []byte    `json:"password" db:"password"`
	Salt        []byte    `json:"salt" db:"salt"`
	PhoneNumber *string   `json:"phone_number" db:"phone_number"`
}

type ReturnUser struct {
	Id          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	PhoneNumber *string   `json:"phone_number"`
}
