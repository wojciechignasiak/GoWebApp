package service

import (
	"app/internal/model"
	servicecomponent "app/internal/service_component"
	"reflect"
	"testing"
)

func TestConvertNewUserModelToUserModel(t *testing.T) {
	registrationService := registrationService{}
	saltGeneratorAndPasswordHasher := servicecomponent.NewSaltGeneratorAndPasswordHasher()
	uuidGenerator := servicecomponent.NewUuidGenerator()

	newUser := model.NewUser{
		Username:        "john_doe",
		Email:           "john.doe@example1.com",
		ConfirmEmail:    "john.doe@example1.com",
		Password:        "secure_password123",
		ConfirmPassword: "secure_password123",
	}
	id, _ := uuidGenerator.GenerateUuid()
	salt, _ := saltGeneratorAndPasswordHasher.GenerateSalt()
	hashedPassword := saltGeneratorAndPasswordHasher.HashPassword(newUser.Password, *salt)

	user := registrationService.convertNewUserModelToUserModel(newUser, *id, *salt, *hashedPassword)

	if reflect.TypeOf(user).Name() != "User" {
		t.Errorf("returned model is not type User, got type: %v", reflect.TypeOf(user).Name())
	}
}
