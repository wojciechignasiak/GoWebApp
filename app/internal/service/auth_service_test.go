package service

import (
	"app/internal/model"
	servicecomponent "app/internal/service_component"
	"testing"

	"github.com/google/uuid"
)

var TestCasesVerifyPassword = []struct {
	name             string
	correctPassword  string
	providedPassword string
	expectedResult   bool
}{
	{
		name:             "success",
		correctPassword:  "secure_password123",
		providedPassword: "secure_password123",
		expectedResult:   true,
	},
	{
		name:             "failure",
		correctPassword:  "secure_password123",
		providedPassword: "secure_password321",
		expectedResult:   false,
	},
}

func TestVerifyPassword(t *testing.T) {

	for _, tc := range TestCasesVerifyPassword {
		t.Run(tc.name, func(t *testing.T) {
			sgaph := servicecomponent.NewSaltGeneratorAndPasswordHasher()
			authService := &authService{
				sgaph: sgaph,
			}
			salt, _ := sgaph.GenerateSalt()

			hashedPassword := sgaph.HashPassword(tc.correctPassword, *salt)
			result := authService.verifyPassword(tc.providedPassword, *hashedPassword, *salt)
			if result != tc.expectedResult {
				t.Errorf("scenario: %s, expected: %v, got: %v", tc.name, tc.expectedResult, result)
			}
		})
	}
}

func TestConvertUserModelToUserSession(t *testing.T) {
	user := model.User{
		Id:                 uuid.MustParse("01945a93-b693-7713-aae0-b0ca59ef4ffd"),
		Username:           "john_doe",
		Email:              "john_doe@domain.com",
		Password:           []byte("secure_password123"),
		Salt:               []byte("G91ApOHlFuMCWuat"),
		RegistrationDate:   nil,
		IsAccountConfirmed: false,
		IsAccountDeleted:   false,
	}

	sgaph := servicecomponent.NewSaltGeneratorAndPasswordHasher()
	authService := &authService{
		sgaph: sgaph,
	}
	userSession := authService.convertUserModelToUserSessionModel(user)

	if userSession.Id != user.Id || userSession.Email != user.Email || userSession.Username != user.Username {
		t.Errorf("expected id: %v, got id: %v, expected email: %v, got email: %v, expected username: %v, got username: %v", user.Id, userSession.Id, user.Email, userSession.Email, user.Username, userSession.Username)
	}
}
