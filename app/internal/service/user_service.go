package service

import (
	apperror "app/internal/app_error"
	"app/internal/database"
	"app/internal/model"
	servicecomponent "app/internal/service_component"
	"context"
	"fmt"
	"math/rand"

	"github.com/google/uuid"
)

type UserService interface {
	CreateUser(ctx context.Context, newUser model.User) *apperror.AppError
	ConfirmAccount(ctx context.Context, confirmAccount model.ConfirmAccount) *apperror.AppError
	GetUserByUsername(ctx context.Context, username string) (*model.User, *apperror.AppError)
	GetUserByEmail(ctx context.Context, email string) (*model.User, *apperror.AppError)
}

type userService struct {
	uowFactory func() (database.UnitOfWork, error)
	ug         servicecomponent.UuidGenerator
}

func NewUserService(uowFactory func() (database.UnitOfWork, error), ug servicecomponent.UuidGenerator) UserService {
	return &userService{
		uowFactory: uowFactory,
		ug:         ug,
	}
}

func (us *userService) CreateUser(ctx context.Context, user model.User) *apperror.AppError {

	uow, err := us.uowFactory()
	if err != nil {
		args := fmt.Sprintf("newUser: anonymized")
		serviceError := apperror.AppError{
			StatusCode:      500,
			Message:         "error occured while creating unit of work in user service",
			StructAndMethod: "UserService.RegisterUser()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      &err,
		}
		return &serviceError
	}

	uowError := uow.BeginTransaction()
	if uowError != nil {
		args := fmt.Sprintf("newUser: anonymized")
		serviceError := apperror.AppError{
			StatusCode:      uowError.StatusCode,
			Message:         uowError.Message,
			StructAndMethod: "UserService.RegisterUser()",
			Argument:        &args,
			ChildAppError:   uowError,
			ChildError:      nil,
		}
		return &serviceError
	}
	repositoryError := uow.UserRepository().CreateUser(ctx, user)
	if repositoryError != nil {
		args := fmt.Sprintf("newUser: anonymized")
		serviceError := apperror.AppError{
			StatusCode:      repositoryError.StatusCode,
			Message:         repositoryError.Message,
			StructAndMethod: "UserService.createUser()",
			Argument:        &args,
			ChildAppError:   repositoryError,
			ChildError:      nil,
		}
		return &serviceError
	}

	accountConfirmationError := us.createAccountConfirmation(ctx, uow, user.Id)
	if accountConfirmationError != nil {
		uow.Rollback()
		args := fmt.Sprintf("newUser: anonymized")
		serviceError := apperror.AppError{
			StatusCode:      accountConfirmationError.StatusCode,
			Message:         accountConfirmationError.Message,
			StructAndMethod: "UserService.RegisterUser()",
			Argument:        &args,
			ChildAppError:   accountConfirmationError,
			ChildError:      nil,
		}
		return &serviceError
	}

	uowError = uow.Commit()
	if uowError != nil {
		args := fmt.Sprintf("newUser: anonymized")
		serviceError := apperror.AppError{
			StatusCode:      uowError.StatusCode,
			Message:         uowError.Message,
			StructAndMethod: "UserService.RegisterUser()",
			Argument:        &args,
			ChildAppError:   uowError,
			ChildError:      nil,
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
			StructAndMethod: "UserService.GetUserByEmail()",
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
			StructAndMethod: "UserService.GetUserByEmail()",
			Argument:        &args,
			ChildAppError:   repositoryError,
			ChildError:      nil,
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
			StructAndMethod: "UserService.GetUserByUsername()",
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
			StructAndMethod: "UserService.GetUserByUsername()",
			Argument:        &args,
			ChildAppError:   repositoryError,
			ChildError:      nil,
		}

		return nil, &serviceError
	}

	return user, nil
}

func (us *userService) createAccountConfirmation(ctx context.Context, uow database.UnitOfWork, userId uuid.UUID) *apperror.AppError {
	accountConfirmationUuid, generationError := us.ug.GenerateUuid()
	if generationError != nil {
		args := fmt.Sprintf("userId: %v", userId)
		serviceError := apperror.AppError{
			StatusCode:      generationError.StatusCode,
			Message:         generationError.Message,
			StructAndMethod: "UserService.createAccountConfirmation()",
			Argument:        &args,
			ChildAppError:   generationError,
			ChildError:      nil,
		}
		return &serviceError
	}

	securityCode := us.generate6DigitCodeForAccountConfirmation()
	accountConfirmation := model.AccountConfirmation{
		UserId:           userId,
		ConfirmationCode: *accountConfirmationUuid,
		SecurityCode:     securityCode,
	}

	repositoryError := uow.UserRepository().CreateAccountConfirmation(ctx, accountConfirmation)
	if repositoryError != nil {
		args := fmt.Sprintf("userId: %v", userId)
		serviceError := apperror.AppError{
			StatusCode:      repositoryError.StatusCode,
			Message:         repositoryError.Message,
			StructAndMethod: "UserService.createAccountConfirmation()",
			Argument:        &args,
			ChildAppError:   repositoryError,
			ChildError:      nil,
		}
		return &serviceError
	}

	return nil
}

func (us *userService) generate6DigitCodeForAccountConfirmation() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func (us *userService) ConfirmAccount(ctx context.Context, confirmAccount model.ConfirmAccount) *apperror.AppError {

	accountConfirmation, getAccountConfirmationError := us.getAccountConfirmationByConfirmationCode(ctx, confirmAccount.ConfirmationCode)
	if getAccountConfirmationError != nil {
		args := fmt.Sprintf("confirmAccount: %v", confirmAccount)
		serviceError := apperror.AppError{
			StatusCode:      getAccountConfirmationError.StatusCode,
			Message:         getAccountConfirmationError.Message,
			StructAndMethod: "UserService.ConfirmAccount()",
			Argument:        &args,
			ChildAppError:   getAccountConfirmationError,
			ChildError:      nil,
		}
		return &serviceError
	}
	if accountConfirmation == nil {
		args := fmt.Sprintf("confirmAccount: %v", confirmAccount)
		serviceError := apperror.AppError{
			StatusCode:      404,
			Message:         "content not found",
			StructAndMethod: "UserService.ConfirmAccount()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      nil,
		}
		return &serviceError
	}

	if confirmAccount.SecurityCode != accountConfirmation.SecurityCode {
		args := fmt.Sprintf("confirmAccount: %v", confirmAccount)
		serviceError := apperror.AppError{
			StatusCode:      404,
			Message:         "content not found",
			StructAndMethod: "UserService.ConfirmAccount()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      nil,
		}
		return &serviceError
	}

	user, getUserError := us.getUserById(ctx, accountConfirmation.UserId)
	if getUserError != nil {
		args := fmt.Sprintf("confirmAccount: %v", confirmAccount)
		serviceError := apperror.AppError{
			StatusCode:      getUserError.StatusCode,
			Message:         getUserError.Message,
			StructAndMethod: "UserService.ConfirmAccount()",
			Argument:        &args,
			ChildAppError:   getUserError,
			ChildError:      nil,
		}
		return &serviceError
	}

	if user.IsAccountDeleted {
		args := fmt.Sprintf("confirmAccount: %v", confirmAccount)
		serviceError := apperror.AppError{
			StatusCode:      404,
			Message:         "content not found",
			StructAndMethod: "UserService.ConfirmAccount()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      nil,
		}
		return &serviceError
	}

	if user.IsAccountConfirmed {
		args := fmt.Sprintf("confirmAccount: %v", confirmAccount)
		serviceError := apperror.AppError{
			StatusCode:      200,
			Message:         "account already confirmed",
			StructAndMethod: "UserService.ConfirmAccount()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      nil,
		}
		return &serviceError
	}

	uow, err := us.uowFactory()
	if err != nil {
		args := fmt.Sprintf("confirmAccount: %v", confirmAccount)
		serviceError := apperror.AppError{
			StatusCode:      500,
			Message:         "error occured while creating unit of work in user service",
			StructAndMethod: "UserService.ConfirmAccount()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      &err,
		}
		return &serviceError
	}

	uowError := uow.BeginTransaction()
	if uowError != nil {
		args := fmt.Sprintf("confirmAccount: %v", confirmAccount)
		serviceError := apperror.AppError{
			StatusCode:      uowError.StatusCode,
			Message:         uowError.Message,
			StructAndMethod: "UserService.ConfirmAccount()",
			Argument:        &args,
			ChildAppError:   uowError,
			ChildError:      nil,
		}
		return &serviceError
	}

	confirmUserAccountError := us.setUserIsConfirmedStatusToTrue(ctx, uow, user.Id)
	if confirmUserAccountError != nil {
		uow.Rollback()
		args := fmt.Sprintf("confirmAccount: %v", confirmAccount)
		serviceError := apperror.AppError{
			StatusCode:      confirmUserAccountError.StatusCode,
			Message:         confirmUserAccountError.Message,
			StructAndMethod: "UserService.ConfirmAccount()",
			Argument:        &args,
			ChildAppError:   uowError,
			ChildError:      nil,
		}
		return &serviceError
	}

	uow.Commit()
	return nil
}

func (us *userService) getAccountConfirmationByConfirmationCode(ctx context.Context, confirmationCode uuid.UUID) (*model.AccountConfirmation, *apperror.AppError) {

	uow, err := us.uowFactory()

	if err != nil {
		args := fmt.Sprintf("confirmationCode: %s", confirmationCode)
		serviceError := apperror.AppError{
			StatusCode:      500,
			Message:         "error occured while creating unit of work in user service",
			StructAndMethod: "UserService.getAccountConfirmationByConfirmationCode()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      &err,
		}
		return nil, &serviceError
	}

	accountConfirmation, repositoryError := uow.UserRepository().GetAccountConfirmationByConfirmationCode(ctx, confirmationCode)
	if repositoryError != nil {
		args := fmt.Sprintf("confirmationCode: %v", confirmationCode)
		serviceError := apperror.AppError{
			StatusCode:      repositoryError.StatusCode,
			Message:         repositoryError.Message,
			StructAndMethod: "UserService.getAccountConfirmationByConfirmationCode()",
			Argument:        &args,
			ChildAppError:   repositoryError,
			ChildError:      nil,
		}

		return nil, &serviceError
	}

	return accountConfirmation, nil
}

func (us *userService) getUserById(ctx context.Context, userId uuid.UUID) (*model.User, *apperror.AppError) {
	uow, err := us.uowFactory()

	if err != nil {
		args := fmt.Sprintf("userId: %s", userId)
		serviceError := apperror.AppError{
			StatusCode:      500,
			Message:         "Error occured while creating unit of work in user service",
			StructAndMethod: "UserService.getUserById()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      &err,
		}

		return nil, &serviceError
	}

	user, repositoryError := uow.UserRepository().GetUserById(ctx, userId)

	if repositoryError != nil {
		args := fmt.Sprintf("userId: %v", userId)
		serviceError := apperror.AppError{
			StatusCode:      repositoryError.StatusCode,
			Message:         repositoryError.Message,
			StructAndMethod: "UserService.GetUserById()",
			Argument:        &args,
			ChildAppError:   repositoryError,
			ChildError:      nil,
		}

		return nil, &serviceError
	}

	return user, nil
}

func (us *userService) setUserIsConfirmedStatusToTrue(ctx context.Context, uow database.UnitOfWork, userId uuid.UUID) *apperror.AppError {
	repositoryError := uow.UserRepository().ConfirmUserAccount(ctx, userId)
	if repositoryError != nil {
		args := fmt.Sprintf("userId: %v", userId)
		serviceError := apperror.AppError{
			StatusCode:      repositoryError.StatusCode,
			Message:         repositoryError.Message,
			StructAndMethod: "UserService.setUserIsConfirmedStatusToTrue()",
			Argument:        &args,
			ChildAppError:   repositoryError,
			ChildError:      nil,
		}
		return &serviceError
	}
	return nil
}
