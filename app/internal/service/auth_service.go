package service

import (
	apperror "app/internal/app_error"
	"app/internal/model"
	"context"
	"crypto/rand"
	"crypto/subtle"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
)

type AuthService interface {
	Login(ctx context.Context, credentials model.Credentials) (*uuid.UUID, *apperror.AppError)
	HashPassword(password string, salt []byte) *[]byte
	GenerateSalt(length int) (*[]byte, *apperror.AppError)
}

type authService struct {
	us  UserService
	sms SessionManagementService
}

func NewAuthService(us UserService, sms SessionManagementService) AuthService {
	return &authService{
		us:  us,
		sms: sms,
	}
}

func (as *authService) Login(ctx context.Context, credentials model.Credentials) (*uuid.UUID, *apperror.AppError) {
	user, userServiceError := as.us.GetUserByUsername(ctx, credentials.Username)
	if userServiceError != nil {
		credentials.Password = "anonymized"
		args := fmt.Sprintf("credentials %v", credentials)
		serviceError := apperror.AppError{
			StatusCode:      userServiceError.StatusCode,
			Message:         userServiceError.Message,
			StructAndMethod: "AuthService.Login()",
			Argument:        &args,
			ChildAppError:   userServiceError,
			ChildError:      nil,
		}
		return nil, &serviceError
	}

	if user == nil {
		credentials.Password = "anonymized"
		args := fmt.Sprintf("credentials %v", credentials)
		serviceError := apperror.AppError{
			StatusCode:      401,
			Message:         "Invalid username or password",
			StructAndMethod: "AuthService.Login()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      nil,
		}
		return nil, &serviceError
	}

	if user.IsAccountConfirmed == false {
		credentials.Password = "anonymized"
		args := fmt.Sprintf("credentials %v", credentials)
		serviceError := apperror.AppError{
			StatusCode:      401,
			Message:         "Account not confirmed",
			StructAndMethod: "AuthService.Login()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      nil,
		}
		return nil, &serviceError
	}

	if user.IsAccountDeleted == true {
		credentials.Password = "anonymized"
		args := fmt.Sprintf("credentials %v", credentials)
		serviceError := apperror.AppError{
			StatusCode:      401,
			Message:         "Invalid username or password",
			StructAndMethod: "AuthService.Login()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      nil,
		}
		return nil, &serviceError
	}

	isPasswordCorrect := as.verifyPassword(credentials.Password, user.Password, user.Salt)
	if isPasswordCorrect == false {
		credentials.Password = "anonymized"
		args := fmt.Sprintf("credentials %v", credentials)
		serviceError := apperror.AppError{
			StatusCode:      401,
			Message:         "Invalid username or password",
			StructAndMethod: "AuthService.Login()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      nil,
		}
		return nil, &serviceError
	}

	userSession := as.convertUserModelToUserSession(*user)
	sessionId, sessionError := as.sms.CreateSession(*userSession)

	if sessionError != nil {
		credentials.Password = "anonymized"
		args := fmt.Sprintf("credentials %v", credentials)
		serviceError := apperror.AppError{
			StatusCode:      sessionError.StatusCode,
			Message:         sessionError.Message,
			StructAndMethod: "AuthService.Login()",
			Argument:        &args,
			ChildAppError:   sessionError,
			ChildError:      nil,
		}
		return nil, &serviceError
	}

	return sessionId, nil
}

func (as *authService) convertUserModelToUserSession(user model.User) *model.UserSession {
	userSession := model.UserSession{
		Id:       user.Id,
		Email:    user.Email,
		Username: user.Username,
	}

	return &userSession
}

func (as *authService) verifyPassword(provided_password string, user_password, salt []byte) bool {
	hashedPassword := as.HashPassword(provided_password, salt)
	return subtle.ConstantTimeCompare(*hashedPassword, user_password) == 1
}

func (as *authService) HashPassword(password string, salt []byte) *[]byte {
	timeCost := uint32(3)
	memoryCost := uint32(64 * 1024)
	threads := uint8(4)
	keyLength := uint32(32)
	hash := argon2.IDKey([]byte(password), salt, timeCost, memoryCost, threads, keyLength)
	return &hash
}

func (as *authService) GenerateSalt(length int) (*[]byte, *apperror.AppError) {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		args := fmt.Sprintf("length: %d", length)
		generationError := apperror.AppError{
			StatusCode:      500,
			Message:         "Error occured while generating salt",
			StructAndMethod: "AuthService.GenerateSalt()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      &err,
		}
		return nil, &generationError
	}
	return &salt, nil
}
