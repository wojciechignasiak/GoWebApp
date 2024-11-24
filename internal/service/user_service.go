package service

import (
	apperror "app/internal/app_error"
	"app/internal/database"
	"app/internal/model"
	servicecomponent "app/internal/service_component"
	"context"
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
		serviceError := apperror.AppError{
			StatusCode:    validationError.StatusCode,
			Message:       validationError.Message,
			ChildAppError: validationError,
			ChildError:    validationError.ChildError,
			Logging:       validationError.Logging,
		}

		return &serviceError
	}

	uow, err := us.uowFactory()
	if err != nil {
		serviceError := apperror.AppError{
			StatusCode:    500,
			Message:       "error occured while creating unit of work in user service",
			ChildAppError: nil,
			ChildError:    &err,
			Logging:       true,
		}

		return &serviceError
	}

	uowError := uow.BeginTransaction()
	if uowError != nil {
		serviceError := apperror.AppError{
			StatusCode:    uowError.StatusCode,
			Message:       uowError.Message,
			ChildAppError: uowError,
			ChildError:    uowError.ChildError,
			Logging:       uowError.Logging,
		}

		return &serviceError
	}

	userWithTheSameEmail, repositoryError := uow.UserRepository().GetUserByEmail(ctx, newUser.Email)
	if repositoryError != nil {
		serviceError := apperror.AppError{
			StatusCode:    repositoryError.StatusCode,
			Message:       repositoryError.Message,
			ChildAppError: repositoryError,
			ChildError:    repositoryError.ChildError,
			Logging:       repositoryError.Logging,
		}

		return &serviceError
	}
	if userWithTheSameEmail != nil {
		serviceError := apperror.AppError{
			StatusCode: 409,
			Message:    "email adress already in use",
			ChildError: nil,
			Logging:    false,
		}

		return &serviceError
	}

	userWithTheSameUsername, repositoryError := uow.UserRepository().GetUserByUsername(ctx, newUser.Username)
	if repositoryError != nil {
		serviceError := apperror.AppError{
			StatusCode:    repositoryError.StatusCode,
			Message:       repositoryError.Message,
			ChildAppError: repositoryError,
			ChildError:    repositoryError.ChildError,
			Logging:       repositoryError.Logging,
		}

		return &serviceError
	}
	if userWithTheSameUsername != nil {
		serviceError := apperror.AppError{
			StatusCode: 409,
			Message:    "username adress already in use",
			ChildError: nil,
			Logging:    false,
		}
		return &serviceError

	}

	uuid, generationError := us.ct.GenerateUUID()
	if generationError != nil {
		serviceError := apperror.AppError{
			StatusCode:    generationError.StatusCode,
			Message:       generationError.Message,
			ChildAppError: generationError,
			ChildError:    generationError.ChildError,
			Logging:       generationError.Logging,
		}

		return &serviceError
	}

	salt, generationError := us.ct.GenerateSalt(16)
	if generationError != nil {
		serviceError := apperror.AppError{
			StatusCode:    generationError.StatusCode,
			Message:       generationError.Message,
			ChildAppError: generationError,
			ChildError:    generationError.ChildError,
			Logging:       generationError.Logging,
		}

		return &serviceError
	}

	hashedPassword := us.ct.HashPassword(newUser.Password, *salt)

	user := model.User{
		Id:          *uuid,
		Username:    newUser.Username,
		Email:       newUser.Email,
		Password:    hashedPassword,
		Salt:        *salt,
		PhoneNumber: newUser.PhoneNumber,
	}

	defer func() {
		if p := recover(); p != nil {
			_ = uow.Rollback()
			panic(p)
		} else if err != nil {
			_ = uow.Rollback()
		} else {
			_ = uow.Commit()
		}
	}()

	repositoryError = uow.UserRepository().CreateUser(ctx, user)
	if repositoryError != nil {
		serviceError := apperror.AppError{
			StatusCode:    repositoryError.StatusCode,
			Message:       repositoryError.Message,
			ChildAppError: repositoryError,
			ChildError:    repositoryError.ChildError,
			Logging:       repositoryError.Logging,
		}

		return &serviceError
	}
	return nil
}

func (us *userService) GetUserByEmail(ctx context.Context, email string) (*model.User, *apperror.AppError) {
	uow, err := us.uowFactory()

	if err != nil {
		serviceError := apperror.AppError{
			StatusCode:    500,
			Message:       "error occured while creating unit of work in user service",
			ChildAppError: nil,
			ChildError:    &err,
			Logging:       true,
		}

		return nil, &serviceError
	}

	user, repositoryError := uow.UserRepository().GetUserByEmail(ctx, email)

	if repositoryError != nil {
		serviceError := apperror.AppError{
			StatusCode:    repositoryError.StatusCode,
			Message:       repositoryError.Message,
			ChildAppError: repositoryError,
			ChildError:    repositoryError.ChildError,
			Logging:       repositoryError.Logging,
		}

		return nil, &serviceError
	}

	return user, nil
}

func (us *userService) GetUserByUsername(ctx context.Context, username string) (*model.User, *apperror.AppError) {
	uow, err := us.uowFactory()

	if err != nil {
		serviceError := apperror.AppError{
			StatusCode:    500,
			Message:       "error occured while creating unit of work in user service",
			ChildAppError: nil,
			ChildError:    &err,
			Logging:       true,
		}

		return nil, &serviceError
	}

	user, repositoryError := uow.UserRepository().GetUserByUsername(ctx, username)

	if repositoryError != nil {
		serviceError := apperror.AppError{
			StatusCode:    repositoryError.StatusCode,
			Message:       repositoryError.Message,
			ChildAppError: repositoryError,
			ChildError:    repositoryError.ChildError,
			Logging:       repositoryError.Logging,
		}

		return nil, &serviceError
	}

	return user, nil
}
