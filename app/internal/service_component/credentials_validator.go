package servicecomponent

import (
	apperror "app/internal/app_error"
	"fmt"
	"regexp"
	"unicode"
)

type CredentialsValidator interface {
	ValidateUsername(username string) *apperror.AppError
	ValidateEmails(email, confirmEmail string) *apperror.AppError
	ValidatePasswords(password, confirmPassword string) *apperror.AppError
}

type credentialsValidator struct{}

func NewCredentialsValidator() *credentialsValidator {
	return &credentialsValidator{}
}

func (cv *credentialsValidator) ValidateUsername(username string) *apperror.AppError {
	if len(username) < 5 || len(username) > 20 {
		args := fmt.Sprintf("username: %s", username)
		validationError := apperror.AppError{
			StatusCode:      400,
			Message:         "username must contain between 5 and 20 characters",
			StructAndMethod: "CredentialsValidator.validateUsername()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      nil,
		}
		return &validationError
	}
	return nil
}

func (cv *credentialsValidator) ValidateEmails(email, confirmEmail string) *apperror.AppError {
	if !cv.areEmailsTheSame(email, confirmEmail) {
		args := fmt.Sprintf("email: %s, confirmEmail: %s", email, confirmEmail)
		validationError := apperror.AppError{
			StatusCode:      400,
			Message:         "provided emails do not match",
			StructAndMethod: "CredentialsValidator.validateEmails()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      nil,
		}
		return &validationError
	}
	if !cv.isValidEmail(email) {
		args := fmt.Sprintf("email: %s, confirmEmail: %s", email, confirmEmail)
		validationError := apperror.AppError{
			StatusCode:      400,
			Message:         "invalid email format",
			StructAndMethod: "CredentialsValidator.validateEmails()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      nil,
		}
		return &validationError
	}
	return nil
}

func (cv *credentialsValidator) areEmailsTheSame(email, confirmEmail string) bool {
	if email != confirmEmail {
		return false
	}
	return true
}

func (cv *credentialsValidator) isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func (cv *credentialsValidator) ValidatePasswords(password, confirmPassword string) *apperror.AppError {
	if !cv.arePasswordsTheSame(password, confirmPassword) {
		validationError := apperror.AppError{
			StatusCode:      400,
			Message:         "provided passwords are not the same",
			StructAndMethod: "CredentialsValidator.validatePasswords()",
			Argument:        nil,
			ChildAppError:   nil,
			ChildError:      nil,
		}
		return &validationError

	}
	if !cv.isPasswordLongEnough(password) {
		validationError := apperror.AppError{
			StatusCode:      400,
			Message:         "password must contain at least 8 characters",
			StructAndMethod: "CredentialsValidator.validatePasswords()",
			Argument:        nil,
			ChildAppError:   nil,
			ChildError:      nil,
		}
		return &validationError
	}

	if !cv.doesPasswordContainsSpecialCharacters(password) {
		validationError := apperror.AppError{
			StatusCode:      403,
			Message:         "password must contain at least one digit and one special character",
			StructAndMethod: "CredentialsValidator.validatePasswords()",
			Argument:        nil,
			ChildAppError:   nil,
			ChildError:      nil,
		}
		return &validationError
	}

	return nil
}

func (cv *credentialsValidator) arePasswordsTheSame(password, confirmPassword string) bool {
	if password != confirmPassword {
		return false
	}
	return true
}

func (cv *credentialsValidator) isPasswordLongEnough(password string) bool {
	if len(password) < 8 {
		return false
	}
	return true
}

func (cv *credentialsValidator) doesPasswordContainsSpecialCharacters(password string) bool {
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
		return false
	}
	return true
}
