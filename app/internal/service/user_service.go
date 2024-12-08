package service

import (
	apperror "app/internal/app_error"
	"app/internal/database"
	"app/internal/model"
	servicecomponent "app/internal/service_component"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type UserService interface {
	CreateUser(ctx context.Context, newUser model.CreateUser) *apperror.AppError
}

type userService struct {
	uowFactory func() (database.UnitOfWork, error)
	ct         servicecomponent.CommonTools
}

func NewUserService(uowFactory func() (database.UnitOfWork, error), ct servicecomponent.CommonTools) *userService {
	return &userService{
		uowFactory: uowFactory,
		ct:         ct,
	}
}

func (us *userService) CreateUser(ctx context.Context, newUser model.CreateUser) *apperror.AppError {

	validationError := newUser.ValidateCreateUser()
	if validationError != nil {
		newUser.Password = "anonimized"
		newUser.ConfirmPassword = "anonimized"
		args := fmt.Sprintf("newUser: %v", newUser)
		serviceError := apperror.AppError{
			StatusCode:      validationError.StatusCode,
			Message:         validationError.Message,
			StructAndMethod: "userService.CreateUser()",
			Argument:        &args,
			ChildAppError:   validationError,
			ChildError:      validationError.ChildError,
		}

		return &serviceError
	}

	uow, err := us.uowFactory()
	if err != nil {
		newUser.Password = "anonimized"
		newUser.ConfirmPassword = "anonimized"
		args := fmt.Sprintf("newUser: %v", newUser)
		serviceError := apperror.AppError{
			StatusCode:      500,
			Message:         "error occured while creating unit of work in user service",
			StructAndMethod: "userService.CreateUser()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      &err,
		}

		return &serviceError
	}

	uowError := uow.BeginTransaction()
	if uowError != nil {
		uow.Rollback()
		newUser.Password = "anonimized"
		newUser.ConfirmPassword = "anonimized"
		args := fmt.Sprintf("newUser: %v", newUser)
		serviceError := apperror.AppError{
			StatusCode:      uowError.StatusCode,
			Message:         uowError.Message,
			StructAndMethod: "userService.CreateUser()",
			Argument:        &args,
			ChildAppError:   uowError,
			ChildError:      uowError.ChildError,
		}

		return &serviceError
	}

	userWithTheSameEmail, repositoryError := uow.UserRepository().GetUserByEmail(ctx, newUser.Email)
	if repositoryError != nil {
		uow.Rollback()
		args := fmt.Sprintf("newUser: %v", newUser)
		serviceError := apperror.AppError{
			StatusCode:      repositoryError.StatusCode,
			Message:         repositoryError.Message,
			StructAndMethod: "userService.CreateUser()",
			Argument:        &args,
			ChildAppError:   repositoryError,
			ChildError:      repositoryError.ChildError,
		}

		return &serviceError
	}
	if userWithTheSameEmail != nil {
		uow.Rollback()
		newUser.Password = "anonimized"
		newUser.ConfirmPassword = "anonimized"
		args := fmt.Sprintf("newUser: %v", newUser)
		serviceError := apperror.AppError{
			StatusCode:      409,
			Message:         "email adress already in use",
			StructAndMethod: "userService.CreateUser()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      nil,
		}

		return &serviceError
	}

	userWithTheSameUsername, repositoryError := uow.UserRepository().GetUserByUsername(ctx, newUser.Username)
	if repositoryError != nil {
		newUser.Password = "anonimized"
		newUser.ConfirmPassword = "anonimized"
		args := fmt.Sprintf("newUser: %v", newUser)
		serviceError := apperror.AppError{
			StatusCode:      repositoryError.StatusCode,
			Message:         repositoryError.Message,
			StructAndMethod: "userService.CreateUser()",
			Argument:        &args,
			ChildAppError:   repositoryError,
			ChildError:      repositoryError.ChildError,
		}

		return &serviceError
	}
	if userWithTheSameUsername != nil {
		newUser.Password = "anonimized"
		newUser.ConfirmPassword = "anonimized"
		args := fmt.Sprintf("newUser: %v", newUser)
		serviceError := apperror.AppError{
			StatusCode:      409,
			Message:         "Username adress already in use",
			StructAndMethod: "userService.CreateUser()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      nil,
		}
		return &serviceError

	}

	userUuid, generationError := us.ct.GenerateUUID()
	if generationError != nil {
		newUser.Password = "anonimized"
		newUser.ConfirmPassword = "anonimized"
		args := fmt.Sprintf("newUser: %v", newUser)
		serviceError := apperror.AppError{
			StatusCode:      generationError.StatusCode,
			Message:         generationError.Message,
			StructAndMethod: "userService.CreateUser()",
			Argument:        &args,
			ChildAppError:   generationError,
			ChildError:      generationError.ChildError,
		}
		return &serviceError
	}

	salt, generationError := us.ct.GenerateSalt(16)
	if generationError != nil {
		newUser.Password = "anonimized"
		newUser.ConfirmPassword = "anonimized"
		args := fmt.Sprintf("newUser: %v", newUser)
		serviceError := apperror.AppError{
			StatusCode:      generationError.StatusCode,
			Message:         generationError.Message,
			StructAndMethod: "userService.CreateUser()",
			Argument:        &args,
			ChildAppError:   generationError,
			ChildError:      generationError.ChildError,
		}

		return &serviceError
	}

	hashedPassword := us.ct.HashPassword(newUser.Password, *salt)

	user := model.User{
		Id:          *userUuid,
		Username:    newUser.Username,
		Email:       newUser.Email,
		Password:    hashedPassword,
		Salt:        *salt,
		PhoneNumber: newUser.PhoneNumber,
	}

	repositoryError = uow.UserRepository().CreateUser(ctx, user)
	if repositoryError != nil {
		uow.Rollback()
		newUser.Password = "anonimized"
		newUser.ConfirmPassword = "anonimized"
		args := fmt.Sprintf("newUser: %v", newUser)
		serviceError := apperror.AppError{
			StatusCode:      repositoryError.StatusCode,
			Message:         repositoryError.Message,
			StructAndMethod: "userService.CreateUser()",
			Argument:        &args,
			ChildAppError:   repositoryError,
			ChildError:      repositoryError.ChildError,
		}
		return &serviceError
	}

	serviceError := us.createAccountConfirmation(ctx, uow, user.Id)
	if serviceError != nil {
		uow.Rollback()
		newUser.Password = "anonimized"
		newUser.ConfirmPassword = "anonimized"
		args := fmt.Sprintf("newUser: %v", newUser)
		serviceError := apperror.AppError{
			StatusCode:      serviceError.StatusCode,
			Message:         serviceError.Message,
			StructAndMethod: "userService.CreateUser()",
			Argument:        &args,
			ChildAppError:   serviceError,
			ChildError:      serviceError.ChildError,
		}
		return &serviceError
	}

	uow.Commit()
	return nil
}

func (us *userService) createAccountConfirmation(ctx context.Context, uow database.UnitOfWork, userId uuid.UUID) *apperror.AppError {
	accountConfirmationUuid, generationError := us.ct.GenerateUUID()
	if generationError != nil {
		args := fmt.Sprintf("userId: %v", userId)
		serviceError := apperror.AppError{
			StatusCode:      generationError.StatusCode,
			Message:         generationError.Message,
			StructAndMethod: "userService.createAccountConfirmation()",
			Argument:        &args,
			ChildAppError:   generationError,
			ChildError:      generationError.ChildError,
		}
		return &serviceError
	}
	accountConfirmation := model.AccountConfirmation{
		UserId:           userId,
		ConfirmationCode: *accountConfirmationUuid,
	}

	repositoryError := uow.UserRepository().CreateAccountConfirmation(ctx, accountConfirmation)
	if repositoryError != nil {
		args := fmt.Sprintf("userId: %v", userId)
		serviceError := apperror.AppError{
			StatusCode:      repositoryError.StatusCode,
			Message:         repositoryError.Message,
			StructAndMethod: "userService.createAccountConfirmation()",
			Argument:        &args,
			ChildAppError:   repositoryError,
			ChildError:      repositoryError.ChildError,
		}
		return &serviceError
	}

	return nil
}

func (us *userService) GetUserByEmail(ctx context.Context, email string) (*model.User, *apperror.AppError) {
	uow, err := us.uowFactory()

	if err != nil {
		args := fmt.Sprintf("newUser: %s", email)
		serviceError := apperror.AppError{
			StatusCode:      500,
			Message:         "Error occured while creating unit of work in user service",
			StructAndMethod: "userService.GetUserByEmail()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      &err,
		}

		return nil, &serviceError
	}

	user, repositoryError := uow.UserRepository().GetUserByEmail(ctx, email)

	if repositoryError != nil {
		args := fmt.Sprintf("newUser: %s", email)
		serviceError := apperror.AppError{
			StatusCode:      repositoryError.StatusCode,
			Message:         repositoryError.Message,
			StructAndMethod: "userService.GetUserByEmail()",
			Argument:        &args,
			ChildAppError:   repositoryError,
			ChildError:      repositoryError.ChildError,
		}

		return nil, &serviceError
	}

	return user, nil
}

func (us *userService) GetUserByUsername(ctx context.Context, username string) (*model.User, *apperror.AppError) {
	uow, err := us.uowFactory()

	if err != nil {
		args := fmt.Sprintf("newUser: %s", username)
		serviceError := apperror.AppError{
			StatusCode:      500,
			Message:         "Error occured while creating unit of work in user service",
			StructAndMethod: "userService.GetUserByUsername()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      &err,
		}

		return nil, &serviceError
	}

	user, repositoryError := uow.UserRepository().GetUserByUsername(ctx, username)

	if repositoryError != nil {
		args := fmt.Sprintf("newUser: %s", username)
		serviceError := apperror.AppError{
			StatusCode:      repositoryError.StatusCode,
			Message:         repositoryError.Message,
			StructAndMethod: "userService.GetUserByUsername()",
			Argument:        &args,
			ChildAppError:   repositoryError,
			ChildError:      repositoryError.ChildError,
		}

		return nil, &serviceError
	}

	return user, nil
}
