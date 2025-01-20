package service

import (
	apperror "app/internal/app_error"
	servicecomponent "app/internal/service_component"
	"context"
	"crypto/subtle"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
)

type AuthService interface {
}

type authService struct {
	us UserService
	ct servicecomponent.CommonTools
}

func NewAuthService(us UserService, ct servicecomponent.CommonTools) *authService {
	return &authService{
		us: us,
		ct: ct,
	}
}

func (as *authService) LogIn(ctx context.Context, email, password string) (*uuid.UUID, *apperror.AppError) {
	credentials, userServiceError := as.us.GetUserCredentialsByEmail(ctx, email)
	if userServiceError != nil {
		args := fmt.Sprintf("email: %s, password: anonimized", email)
		serviceError := apperror.AppError{
			StatusCode:      userServiceError.StatusCode,
			Message:         userServiceError.Message,
			StructAndMethod: "AuthService.LogIn()",
			Argument:        &args,
			ChildAppError:   userServiceError,
			ChildError:      nil,
		}
		return nil, &serviceError
	}

	is_password_correct := as.verify_password(password, credentials.Password, credentials.Salt)
	if is_password_correct == false {
		args := fmt.Sprintf("email: %s, password: anonimized", email)
		serviceError := apperror.AppError{
			StatusCode:      401,
			Message:         "Invalid email or password",
			StructAndMethod: "AuthService.LogIn()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      nil,
		}
		return nil, &serviceError
	}

	session_id, ctError := as.ct.GenerateUUID()
	if ctError != nil {
		args := fmt.Sprintf("email: %s, password: anonimized", email)
		serviceError := apperror.AppError{
			StatusCode:      ctError.StatusCode,
			Message:         ctError.Message,
			StructAndMethod: "AuthService.LogIn()",
			Argument:        &args,
			ChildAppError:   ctError,
			ChildError:      nil,
		}
		return nil, &serviceError
	}

	return session_id, nil
}

func (as *authService) verify_password(provided_password string, user_password []byte, salt []byte) bool {
	hashedPassword := as.HashPassword(provided_password, salt)
	return subtle.ConstantTimeCompare(hashedPassword, user_password) == 1
}

func (as *authService) HashPassword(password string, salt []byte) []byte {
	timeCost := uint32(3)
	memoryCost := uint32(64 * 1024)
	threads := uint8(4)
	keyLength := uint32(32)
	hash := argon2.IDKey([]byte(password), salt, timeCost, memoryCost, threads, keyLength)
	return hash
}
