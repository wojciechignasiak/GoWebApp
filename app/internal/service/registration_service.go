package service

import (
	apperror "app/internal/app_error"
	"app/internal/model"
	servicecomponent "app/internal/service_component"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type RegistrationService interface {
	Register(ctx context.Context, newUser model.NewUser) *apperror.AppError
}

type registrationService struct {
	us    UserService
	sgaph servicecomponent.SaltGeneratorAndPasswordHasher
	cv    servicecomponent.CredentialsValidator
	ug    servicecomponent.UuidGenerator
}

func NewRegistrationService(us UserService, sgaph servicecomponent.SaltGeneratorAndPasswordHasher, cv servicecomponent.CredentialsValidator, ug servicecomponent.UuidGenerator) RegistrationService {
	return &registrationService{
		us:    us,
		sgaph: sgaph,
		cv:    cv,
		ug:    ug,
	}
}

func (rs *registrationService) Register(ctx context.Context, newUser model.NewUser) *apperror.AppError {

	validationError := rs.validateCredentials(newUser.Email, newUser.ConfirmEmail, newUser.Password, newUser.ConfirmPassword, newUser.Username)
	if validationError != nil {
		newUser.Password = "anonymized"
		newUser.ConfirmPassword = "anonymized"
		args := fmt.Sprintf("newUser: %v", newUser)
		serviceError := apperror.AppError{
			StatusCode:      validationError.StatusCode,
			Message:         validationError.Message,
			StructAndMethod: "RegistrationService.RegisterUser()",
			Argument:        &args,
			ChildAppError:   validationError,
			ChildError:      nil,
		}
		return &serviceError
	}

	duplicateError := rs.checkisUsernameOrEmailDuplicate(ctx, newUser.Username, newUser.Email)
	if duplicateError != nil {
		newUser.Password = "anonymized"
		newUser.ConfirmPassword = "anonymized"
		args := fmt.Sprintf("newUser: %v", newUser)
		serviceError := apperror.AppError{
			StatusCode:      duplicateError.StatusCode,
			Message:         duplicateError.Message,
			StructAndMethod: "RegistrationService.RegisterUser()",
			Argument:        &args,
			ChildAppError:   duplicateError,
			ChildError:      nil,
		}
		return &serviceError
	}

	id, generationError := rs.ug.GenerateUuid()

	if generationError != nil {
		newUser.Password = "anonymized"
		newUser.ConfirmPassword = "anonymized"
		args := fmt.Sprintf("newUser: %v", newUser)
		serviceError := apperror.AppError{
			StatusCode:      generationError.StatusCode,
			Message:         generationError.Message,
			StructAndMethod: "RegistrationService.RegisterUser()",
			Argument:        &args,
			ChildAppError:   generationError,
			ChildError:      nil,
		}
		return &serviceError
	}

	salt, generationError := rs.sgaph.GenerateSalt()
	if generationError != nil {
		newUser.Password = "anonymized"
		newUser.ConfirmPassword = "anonymized"
		args := fmt.Sprintf("newUser: %v", newUser)
		serviceError := apperror.AppError{
			StatusCode:      generationError.StatusCode,
			Message:         generationError.Message,
			StructAndMethod: "RegistrationService.RegisterUser()",
			Argument:        &args,
			ChildAppError:   generationError,
			ChildError:      nil,
		}
		return &serviceError
	}

	hashed_password := rs.sgaph.HashPassword(newUser.Password, *salt)

	user := rs.convertNewUserModelToUserModel(newUser, *id, *salt, *hashed_password)
	createUserError := rs.us.CreateUser(ctx, *user)
	if createUserError != nil {
		newUser.Password = "anonymized"
		newUser.ConfirmPassword = "anonymized"
		args := fmt.Sprintf("newUser: %v", newUser)
		serviceError := apperror.AppError{
			StatusCode:      createUserError.StatusCode,
			Message:         createUserError.Message,
			StructAndMethod: "RegistrationService.RegisterUser()",
			Argument:        &args,
			ChildAppError:   createUserError,
			ChildError:      nil,
		}
		return &serviceError
	}

	return nil
}

func (rs *registrationService) checkisUsernameOrEmailDuplicate(ctx context.Context, username, email string) *apperror.AppError {
	user, err := rs.us.GetUserByUsername(ctx, username)
	if err != nil {
		args := fmt.Sprintf("username: %s, email: %s", username, email)
		serviceError := apperror.AppError{
			StatusCode:      err.StatusCode,
			Message:         err.Message,
			StructAndMethod: "RegistrationService.checkisUsernameOrEmailDuplicate()",
			Argument:        &args,
			ChildAppError:   err,
			ChildError:      err.ChildError,
		}
		return &serviceError
	}
	if user != nil {
		args := fmt.Sprintf("username: %s, email: %s", username, email)
		serviceError := apperror.AppError{
			StatusCode:      409,
			Message:         "username already in use",
			StructAndMethod: "RegistrationService.checkisUsernameOrEmailDuplicate()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      nil,
		}
		return &serviceError
	}

	user, err = rs.us.GetUserByEmail(ctx, email)
	if err != nil {
		args := fmt.Sprintf("username: %s, email: %s", username, email)
		serviceError := apperror.AppError{
			StatusCode:      err.StatusCode,
			Message:         err.Message,
			StructAndMethod: "RegistrationService.checkisUsernameOrEmailDuplicate()",
			Argument:        &args,
			ChildAppError:   err,
			ChildError:      err.ChildError,
		}
		return &serviceError
	}
	if user != nil {
		args := fmt.Sprintf("username: %s, email: %s", username, email)
		serviceError := apperror.AppError{
			StatusCode:      409,
			Message:         "email already in use",
			StructAndMethod: "RegistrationService.checkisUsernameOrEmailDuplicate()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      nil,
		}
		return &serviceError
	}

	return nil
}

func (rs *registrationService) validateCredentials(email, confirmEmail, password, confirmPassword, username string) *apperror.AppError {
	validationError := rs.cv.ValidateEmails(email, confirmEmail)
	if validationError != nil {
		args := fmt.Sprintf("email: %s, confirmEmail: %s, password: %s, confirmPassword: %s, username: %s", email, confirmEmail, "anonymized", "anonymized", username)
		serviceError := apperror.AppError{
			StatusCode:      validationError.StatusCode,
			Message:         validationError.Message,
			StructAndMethod: "RegistrationService.validateCredentials()",
			Argument:        &args,
			ChildAppError:   validationError,
			ChildError:      nil,
		}
		return &serviceError
	}
	validationError = rs.cv.ValidateUsername(username)
	if validationError != nil {
		args := fmt.Sprintf("email: %s, confirmEmail: %s, password: %s, confirmPassword: %s, username: %s", email, confirmEmail, "anonymized", "anonymized", username)
		serviceError := apperror.AppError{
			StatusCode:      validationError.StatusCode,
			Message:         validationError.Message,
			StructAndMethod: "RegistrationService.validateCredentials()",
			Argument:        &args,
			ChildAppError:   validationError,
			ChildError:      nil,
		}
		return &serviceError
	}
	validationError = rs.cv.ValidatePasswords(password, confirmPassword)
	if validationError != nil {
		args := fmt.Sprintf("email: %s, confirmEmail: %s, password: %s, confirmPassword: %s, username: %s", email, confirmEmail, "anonymized", "anonymized", username)
		serviceError := apperror.AppError{
			StatusCode:      validationError.StatusCode,
			Message:         validationError.Message,
			StructAndMethod: "RegistrationService.validateCredentials()",
			Argument:        &args,
			ChildAppError:   validationError,
			ChildError:      nil,
		}
		return &serviceError
	}

	return nil
}

func (rs *registrationService) convertNewUserModelToUserModel(newUser model.NewUser, id uuid.UUID, salt, hashedPassword []byte) *model.User {
	return &model.User{
		Id:                 id,
		Username:           newUser.Username,
		Email:              newUser.Email,
		Password:           hashedPassword,
		Salt:               salt,
		IsAccountConfirmed: false,
		IsAccountDeleted:   false,
	}
}
