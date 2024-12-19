package service

import (
	apperror "app/internal/app_error"
	"app/internal/database"
	"app/internal/model"
	servicecomponent "app/internal/service_component"
	"context"
	"fmt"
	"regexp"
	"unicode"

	"github.com/google/uuid"
)

type UserService interface {
	RegisterUser(ctx context.Context, newUser model.CreateUser) *apperror.AppError
	ConfirmAccount(ctx context.Context, confirmAccount model.ConfirmAccount) *apperror.AppError
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

func (us *userService) RegisterUser(ctx context.Context, newUser model.CreateUser) *apperror.AppError {
	validationError := us.validateNewUserData(newUser)
	if validationError != nil {
		newUser.Password = "anonimized"
		newUser.ConfirmPassword = "anonimized"
		args := fmt.Sprintf("newUser: %v", newUser)
		serviceError := apperror.AppError{
			StatusCode:      validationError.StatusCode,
			Message:         validationError.Message,
			StructAndMethod: "UserService.RegisterUser()",
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
			StructAndMethod: "UserService.RegisterUser()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      &err,
		}
		return &serviceError
	}

	duplicateError := us.checkisUsernameOrEmailDuplicate(ctx, newUser.Username, newUser.Email)
	if duplicateError != nil {
		newUser.Password = "anonimized"
		newUser.ConfirmPassword = "anonimized"
		args := fmt.Sprintf("newUser: %v", newUser)
		serviceError := apperror.AppError{
			StatusCode:      duplicateError.StatusCode,
			Message:         duplicateError.Message,
			StructAndMethod: "UserService.RegisterUser()",
			Argument:        &args,
			ChildAppError:   duplicateError,
			ChildError:      duplicateError.ChildError,
		}
		return &serviceError
	}

	uowError := uow.BeginTransaction()
	if uowError != nil {
		newUser.Password = "anonimized"
		newUser.ConfirmPassword = "anonimized"
		args := fmt.Sprintf("newUser: %v", newUser)
		serviceError := apperror.AppError{
			StatusCode:      uowError.StatusCode,
			Message:         uowError.Message,
			StructAndMethod: "UserService.RegisterUser()",
			Argument:        &args,
			ChildAppError:   uowError,
			ChildError:      uowError.ChildError,
		}
		return &serviceError
	}

	userID, createUserError := us.createUser(ctx, uow, newUser)
	if createUserError != nil {
		uow.Rollback()
		newUser.Password = "anonimized"
		newUser.ConfirmPassword = "anonimized"
		args := fmt.Sprintf("newUser: %v", newUser)
		serviceError := apperror.AppError{
			StatusCode:      createUserError.StatusCode,
			Message:         createUserError.Message,
			StructAndMethod: "UserService.RegisterUser()",
			Argument:        &args,
			ChildAppError:   createUserError,
			ChildError:      createUserError.ChildError,
		}
		return &serviceError
	}

	accountConfirmationError := us.createAccountConfirmation(ctx, uow, *userID)
	if accountConfirmationError != nil {
		uow.Rollback()
		newUser.Password = "anonimized"
		newUser.ConfirmPassword = "anonimized"
		args := fmt.Sprintf("newUser: %v", newUser)
		serviceError := apperror.AppError{
			StatusCode:      accountConfirmationError.StatusCode,
			Message:         accountConfirmationError.Message,
			StructAndMethod: "UserService.RegisterUser()",
			Argument:        &args,
			ChildAppError:   accountConfirmationError,
			ChildError:      accountConfirmationError.ChildError,
		}
		return &serviceError
	}

	uow.Commit()
	return nil
}

func (us *userService) createUser(ctx context.Context, uow database.UnitOfWork, newUser model.CreateUser) (*uuid.UUID, *apperror.AppError) {

	userUuid, generationError := us.ct.GenerateUUID()
	if generationError != nil {
		newUser.Password = "anonimized"
		newUser.ConfirmPassword = "anonimized"
		args := fmt.Sprintf("newUser: %v", newUser)
		serviceError := apperror.AppError{
			StatusCode:      generationError.StatusCode,
			Message:         generationError.Message,
			StructAndMethod: "UserService.createUser()",
			Argument:        &args,
			ChildAppError:   generationError,
			ChildError:      generationError.ChildError,
		}
		return nil, &serviceError
	}

	salt, generationError := us.ct.GenerateSalt(16)
	if generationError != nil {
		newUser.Password = "anonimized"
		newUser.ConfirmPassword = "anonimized"
		args := fmt.Sprintf("newUser: %v", newUser)
		serviceError := apperror.AppError{
			StatusCode:      generationError.StatusCode,
			Message:         generationError.Message,
			StructAndMethod: "UserService.createUser()",
			Argument:        &args,
			ChildAppError:   generationError,
			ChildError:      generationError.ChildError,
		}

		return nil, &serviceError
	}

	hashedPassword := us.ct.HashPassword(newUser.Password, *salt)

	user := model.User{
		Id:       *userUuid,
		Username: newUser.Username,
		Email:    newUser.Email,
		Password: hashedPassword,
		Salt:     *salt,
	}

	repositoryError := uow.UserRepository().CreateUser(ctx, user)
	if repositoryError != nil {
		newUser.Password = "anonimized"
		newUser.ConfirmPassword = "anonimized"
		args := fmt.Sprintf("newUser: %v", newUser)
		serviceError := apperror.AppError{
			StatusCode:      repositoryError.StatusCode,
			Message:         repositoryError.Message,
			StructAndMethod: "UserService.createUser()",
			Argument:        &args,
			ChildAppError:   repositoryError,
			ChildError:      repositoryError.ChildError,
		}
		return nil, &serviceError
	}

	return userUuid, nil
}

func (us *userService) validateNewUserData(newUser model.CreateUser) *apperror.AppError {

	if err := us.validateUsername(newUser.Username); err != nil {
		args := fmt.Sprintf("newUser: %v", newUser)
		serviceError := apperror.AppError{
			StatusCode:      err.StatusCode,
			Message:         err.Message,
			StructAndMethod: "UserService.validateNewUserData()",
			Argument:        &args,
			ChildAppError:   err,
			ChildError:      err.ChildError,
		}
		return &serviceError
	}

	if err := us.validateEmails(newUser.Email, newUser.ConfirmEmail); err != nil {
		args := fmt.Sprintf("newUser: %v", newUser)
		serviceError := apperror.AppError{
			StatusCode:      err.StatusCode,
			Message:         err.Message,
			StructAndMethod: "UserService.validateNewUserData()",
			Argument:        &args,
			ChildAppError:   err,
			ChildError:      err.ChildError,
		}
		return &serviceError
	}

	if err := us.validatePasswords(newUser.Password, newUser.ConfirmPassword); err != nil {
		args := fmt.Sprintf("newUser: %v", newUser)
		serviceError := apperror.AppError{
			StatusCode:      err.StatusCode,
			Message:         err.Message,
			StructAndMethod: "UserService.validateNewUserData()",
			Argument:        &args,
			ChildAppError:   err,
			ChildError:      err.ChildError,
		}
		return &serviceError
	}

	return nil
}

func (us *userService) validateUsername(username string) *apperror.AppError {
	if len(username) < 5 || len(username) > 20 {
		args := fmt.Sprintf("username: %s", username)
		validationError := apperror.AppError{
			StatusCode:      400,
			Message:         "Username must contain between 5 and 20 characters",
			StructAndMethod: "UserService.validateUsername()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      nil,
		}
		return &validationError
	}
	return nil
}

func (us *userService) validateEmails(email, confirmEmail string) *apperror.AppError {
	if !us.areEmailsTheSame(email, confirmEmail) {
		args := fmt.Sprintf("email: %s, confirmEmail: %s", email, confirmEmail)
		validationError := apperror.AppError{
			StatusCode:      400,
			Message:         "Provided emails do not match",
			StructAndMethod: "UserService.validateEmails()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      nil,
		}
		return &validationError
	}
	if !us.isValidEmail(email) {
		args := fmt.Sprintf("email: %s, confirmEmail: %s", email, confirmEmail)
		validationError := apperror.AppError{
			StatusCode:      400,
			Message:         "Invalid email format",
			StructAndMethod: "UserService.validateEmails()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      nil,
		}
		return &validationError
	}
	return nil
}

func (us *userService) areEmailsTheSame(email, confirmEmail string) bool {
	if email != confirmEmail {
		return false
	}
	return true
}

func (us *userService) isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func (us *userService) validatePasswords(password, confirmPassword string) *apperror.AppError {
	if !us.arePasswordsTheSame(password, confirmPassword) {
		validationError := apperror.AppError{
			StatusCode:      400,
			Message:         "Provided passwords are not the same",
			StructAndMethod: "UserService.validatePasswords()",
			Argument:        nil,
			ChildAppError:   nil,
			ChildError:      nil,
		}
		return &validationError

	}
	if !us.isPasswordLongEnough(password) {
		validationError := apperror.AppError{
			StatusCode:      400,
			Message:         "Password must contain at least 8 characters",
			StructAndMethod: "UserService.validatePasswords()",
			Argument:        nil,
			ChildAppError:   nil,
			ChildError:      nil,
		}
		return &validationError
	}

	if !us.doesPasswordContainsSpecialCharacters(password) {
		validationError := apperror.AppError{
			StatusCode:      403,
			Message:         "Password must contain at least one digit and one special character",
			StructAndMethod: "UserService.validatePasswords()",
			Argument:        nil,
			ChildAppError:   nil,
			ChildError:      nil,
		}
		return &validationError
	}

	return nil
}

func (us *userService) arePasswordsTheSame(password, confirmPassword string) bool {
	if password != confirmPassword {
		return false
	}
	return true
}

func (us *userService) isPasswordLongEnough(password string) bool {
	if len(password) < 8 {
		return false
	}
	return true
}

func (us *userService) doesPasswordContainsSpecialCharacters(password string) bool {
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

func (us *userService) checkisUsernameOrEmailDuplicate(ctx context.Context, username, email string) *apperror.AppError {
	user, err := us.getUserByUsername(ctx, username)
	if err != nil {
		args := fmt.Sprintf("username: %s, email: %s", username, email)
		serviceError := apperror.AppError{
			StatusCode:      err.StatusCode,
			Message:         err.Message,
			StructAndMethod: "UserService.checkisUsernameOrEmailDuplicate()",
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
			StructAndMethod: "UserService.checkisUsernameOrEmailDuplicate()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      nil,
		}
		return &serviceError
	}

	user, err = us.getUserByEmail(ctx, email)
	if err != nil {
		args := fmt.Sprintf("username: %s, email: %s", username, email)
		serviceError := apperror.AppError{
			StatusCode:      err.StatusCode,
			Message:         err.Message,
			StructAndMethod: "UserService.checkisUsernameOrEmailDuplicate()",
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
			StructAndMethod: "UserService.checkisUsernameOrEmailDuplicate()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      nil,
		}
		return &serviceError
	}

	return nil
}

func (us *userService) getUserByEmail(ctx context.Context, email string) (*model.User, *apperror.AppError) {
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
			ChildError:      repositoryError.ChildError,
		}

		return nil, &serviceError
	}

	return user, nil
}

func (us *userService) getUserByUsername(ctx context.Context, username string) (*model.User, *apperror.AppError) {
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
			ChildError:      repositoryError.ChildError,
		}

		return nil, &serviceError
	}

	return user, nil
}

func (us *userService) createAccountConfirmation(ctx context.Context, uow database.UnitOfWork, userId uuid.UUID) *apperror.AppError {
	accountConfirmationUuid, generationError := us.ct.GenerateUUID()
	if generationError != nil {
		args := fmt.Sprintf("userId: %v", userId)
		serviceError := apperror.AppError{
			StatusCode:      generationError.StatusCode,
			Message:         generationError.Message,
			StructAndMethod: "UserService.createAccountConfirmation()",
			Argument:        &args,
			ChildAppError:   generationError,
			ChildError:      generationError.ChildError,
		}
		return &serviceError
	}

	securityCode := us.ct.Generate6DigitCode()
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
			ChildError:      repositoryError.ChildError,
		}
		return &serviceError
	}

	return nil
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
			ChildError:      getAccountConfirmationError.ChildError,
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
			ChildError:      getUserError.ChildError,
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
			ChildError:      uowError.ChildError,
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
			ChildError:      confirmUserAccountError.ChildError,
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
			ChildError:      repositoryError.ChildError,
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
			ChildError:      repositoryError.ChildError,
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
			ChildError:      repositoryError.ChildError,
		}
		return &serviceError
	}
	return nil
}
