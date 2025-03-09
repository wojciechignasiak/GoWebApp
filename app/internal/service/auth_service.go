package service

import (
	apperror "app/internal/app_error"
	"app/internal/model"
	servicecomponent "app/internal/service_component"
	"context"
	"crypto/subtle"
	"fmt"

	"github.com/google/uuid"
)

type AuthService interface {
	Login(ctx context.Context, credentials model.Credentials) (*uuid.UUID, *apperror.AppError)
	Logout(sessionId string)
}

type authService struct {
	us    UserService
	sms   SessionManagementService
	oss   OnlineStatusService
	sgaph servicecomponent.SaltGeneratorAndPasswordHasher
}

func NewAuthService(us UserService, sms SessionManagementService, oss OnlineStatusService, sgaph servicecomponent.SaltGeneratorAndPasswordHasher) AuthService {
	return &authService{
		us:    us,
		sms:   sms,
		oss:   oss,
		sgaph: sgaph,
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

	if !user.IsAccountConfirmed {
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

	if user.IsAccountDeleted {
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
	if !isPasswordCorrect {
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
	userSession := as.convertUserModelToUserSessionModel(*user)
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

	as.oss.SetUserStatusToOnline(user.Id)

	return sessionId, nil
}

func (as *authService) convertUserModelToUserSessionModel(user model.User) *model.UserSession {
	userSession := model.UserSession{
		Id:       user.Id,
		Email:    user.Email,
		Username: user.Username,
	}

	return &userSession
}

func (as *authService) verifyPassword(provided_password string, user_password, salt []byte) bool {
	hashedPassword := as.sgaph.HashPassword(provided_password, salt)
	return subtle.ConstantTimeCompare(*hashedPassword, user_password) == 1
}

func (as *authService) Logout(sessionId string) {
	sessionIdUuid := uuid.MustParse(sessionId)
	userSession := as.sms.GetUserSession(sessionIdUuid)
	if userSession == nil {
		return
	}
	as.sms.DeleteSession(sessionIdUuid)
	as.oss.SetUserStatusToOffline(userSession.Id)
}
