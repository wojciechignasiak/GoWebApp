package integration_test

import (
	"app/internal/controller"
	controllercomponent "app/internal/controller_component"
	"app/internal/database"
	integration "app/internal/integration_test"
	"app/internal/logs"
	"app/internal/server"
	"app/internal/service"
	servicecomponent "app/internal/service_component"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var registerUserTestCases = []struct {
	name               string
	requestBody        map[string]interface{}
	expectedResponse   map[string]interface{}
	expectedStatusCode int
}{
	{
		"success",
		map[string]interface{}{
			"username":         "john_doe",
			"email":            "john.doe@example1.com",
			"confirm_email":    "john.doe@example1.com",
			"password":         "secure_password123",
			"confirm_password": "secure_password123",
		},
		map[string]interface{}{
			"message": "registered successfully",
		},
		201,
	},
	{
		"wrong email format",
		map[string]interface{}{
			"username":         "john_doe",
			"email":            "john.doeexample1.com",
			"confirm_email":    "john.doeexample1.com",
			"password":         "secure_password123",
			"confirm_password": "secure_password123",
		},
		map[string]interface{}{
			"message": "invalid email format",
		},
		400,
	},
	{
		"username too short",
		map[string]interface{}{
			"username":         "j",
			"email":            "john.doe@example1.com",
			"confirm_email":    "john.doe@example1.com",
			"password":         "secure_password123",
			"confirm_password": "secure_password123",
		},
		map[string]interface{}{
			"message": "username must contain between 5 and 20 characters",
		},
		400,
	},
	{
		"username too long",
		map[string]interface{}{
			"username":         "jaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			"email":            "john.doe@example1.com",
			"confirm_email":    "john.doe@example1.com",
			"password":         "secure_password123",
			"confirm_password": "secure_password123",
		},
		map[string]interface{}{
			"message": "username must contain between 5 and 20 characters",
		},
		400,
	},
	{
		"password too short",
		map[string]interface{}{
			"username":         "john_doe",
			"email":            "john.doe@example1.com",
			"confirm_email":    "john.doe@example1.com",
			"password":         "s_p1",
			"confirm_password": "s_p1",
		},
		map[string]interface{}{
			"message": "password must contain at least 8 characters",
		},
		400,
	},
	{
		"password missing special characters",
		map[string]interface{}{
			"username":         "john_doe",
			"email":            "john.doe@example1.com",
			"confirm_email":    "john.doe@example1.com",
			"password":         "mypassworddd",
			"confirm_password": "mypassworddd",
		},
		map[string]interface{}{
			"message": "password must contain at least one digit and one special character",
		},
		403,
	},
	{
		"email not the same",
		map[string]interface{}{
			"username":         "john_doe",
			"email":            "john.doe@example1.com",
			"confirm_email":    "john.doe@example2.com",
			"password":         "secure_password123",
			"confirm_password": "secure_password123",
		},
		map[string]interface{}{
			"message": "provided emails do not match",
		},
		400,
	},
	{
		"password not the same",
		map[string]interface{}{
			"username":         "john_doe",
			"email":            "john.doe@example1.com",
			"confirm_email":    "john.doe@example1.com",
			"password":         "secure_password123",
			"confirm_password": "secure_password1234",
		},
		map[string]interface{}{
			"message": "provided passwords are not the same",
		},
		400,
	},
	{
		"username already in use",
		map[string]interface{}{
			"username":         "john_doe",
			"email":            "john.doe@example1.com",
			"confirm_email":    "john.doe@example1.com",
			"password":         "secure_password123",
			"confirm_password": "secure_password123",
		},
		map[string]interface{}{
			"message": "username already in use",
		},
		409,
	},
	{
		"email already in use",
		map[string]interface{}{
			"username":         "john_doe",
			"email":            "john.doe@example1.com",
			"confirm_email":    "john.doe@example1.com",
			"password":         "secure_password123",
			"confirm_password": "secure_password123",
		},
		map[string]interface{}{
			"message": "email already in use",
		},
		409,
	},
}

func TestIntegration_RegisterUser(t *testing.T) {
	db, err := integration.SetupTestDB()
	if err != nil {
		t.Fatalf("failed setup database connection: %v", err)
	}

	defer db.Close()

	userService := service.NewUserService(func() (database.UnitOfWork, error) {
		return database.NewUnitOfWork(db), nil
	}, servicecomponent.NewUuidGenerator())

	sessionManagementService := service.NewSessionManagementService(servicecomponent.NewUuidGenerator())

	authService := service.NewAuthService(userService, sessionManagementService)

	registrationService := service.NewRegistrationService(authService, userService, servicecomponent.NewCredentialsValidator(), servicecomponent.NewUuidGenerator())

	logger := logs.NewLogger()
	responseHandler := controllercomponent.NewResponseHandler()
	authController := controller.NewAuthController(authService, registrationService, responseHandler, logger)

	testServer := httptest.NewServer(server.NewServer("", 8080, nil, authController).Handler)
	defer testServer.Close()

	for _, tc := range registerUserTestCases {
		t.Run(tc.name, func(t *testing.T) {

			postBody, _ := json.Marshal(tc.requestBody)

			if tc.name == "username already in use" || tc.name == "email already in use" {
				newRequestBody := make(map[string]interface{})
				for k, v := range tc.requestBody {
					newRequestBody[k] = v
				}

				if tc.name == "username already in use" {
					newRequestBody["email"] = "john.doe@example2.com"
					newRequestBody["confirm_email"] = "john.doe@example2.com"
				}

				if tc.name == "email already in use" {
					newRequestBody["username"] = "jon_doe1234"
				}

				postBody, _ = json.Marshal(newRequestBody)

				_, err := http.Post(testServer.URL+"/auth/register", "application/json", bytes.NewReader(postBody))
				if err != nil {
					t.Fatalf("failed to send request to register user for duplicate test case: %v", err)
				}

				postBody, _ = json.Marshal(tc.requestBody)
			}

			resp, err := http.Post(testServer.URL+"/auth/register", "application/json", bytes.NewReader(postBody))
			if err != nil {
				t.Fatalf("failed to send request: %v", err)
			}
			defer resp.Body.Close()

			responseBody, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("failed to read response body: %v", err)
			}

			var responseBodyMap map[string]interface{}
			err = json.Unmarshal(responseBody, &responseBodyMap)
			if err != nil {
				t.Fatalf("failed to unmarshal response body: %v", err)
			}

			if !reflect.DeepEqual(responseBodyMap, tc.expectedResponse) || resp.StatusCode != tc.expectedStatusCode {
				t.Errorf("\nscenario: %s, expected response body: %v, got response body: %v, expected status code: %v, got status code: %v\n",
					tc.name, tc.expectedResponse, responseBodyMap, tc.expectedStatusCode, resp.StatusCode)
			}

			defer integration.TruncateUserTable(db)
		})
	}
}
