package integration_test

import (
	integration_test_tools "app/internal/integration_test"
	"app/internal/model"
	servicecomponent "app/internal/service_component"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

var LoginTestCases = []struct {
	name               string
	user               map[string]interface{}
	requestBody        map[string]string
	expectedResponse   map[string]interface{}
	expectedStatusCode int
}{
	{
		"success",
		map[string]interface{}{
			"id":                   uuid.MustParse("8f50842d-ce0f-4bd1-8411-68a220d797a6"),
			"username":             "john_doe",
			"email":                "john.doe@example1.com",
			"password":             "secure_password123",
			"is_account_confirmed": true,
			"is_account_deleted":   false,
		},
		map[string]string{
			"username": "john_doe",
			"password": "secure_password123",
		},
		map[string]interface{}{
			"message": "Login successful",
		},
		200,
	},
	{
		"invalid password",
		map[string]interface{}{
			"id":                   uuid.MustParse("8f50842d-ce0f-4bd1-8411-68a220d797a6"),
			"username":             "john_doe",
			"email":                "john.doe@example1.com",
			"password":             "secure_password123",
			"is_account_confirmed": true,
			"is_account_deleted":   false,
		},
		map[string]string{
			"username": "john_doe",
			"password": "secure_password123!@#",
		},
		map[string]interface{}{
			"message": "Invalid username or password",
		},
		401,
	},
	{
		"user do not exists or wrong username",
		map[string]interface{}{
			"id":                   uuid.MustParse("8f50842d-ce0f-4bd1-8411-68a220d797a6"),
			"username":             "john_doe",
			"email":                "john.doe@example1.com",
			"password":             "secure_password123",
			"is_account_confirmed": true,
			"is_account_deleted":   false,
		},
		map[string]string{
			"username": "john_doe1",
			"password": "secure_password123!@#",
		},
		map[string]interface{}{
			"message": "Invalid username or password",
		},
		401,
	},
	{
		"user not confirmed",
		map[string]interface{}{
			"id":                   uuid.MustParse("8f50842d-ce0f-4bd1-8411-68a220d797a6"),
			"username":             "john_doe",
			"email":                "john.doe@example1.com",
			"password":             "secure_password123",
			"is_account_confirmed": false,
			"is_account_deleted":   false,
		},
		map[string]string{
			"username": "john_doe",
			"password": "secure_password123",
		},
		map[string]interface{}{
			"message": "Account not confirmed",
		},
		401,
	},
	{
		"user deleted",
		map[string]interface{}{
			"id":                   uuid.MustParse("8f50842d-ce0f-4bd1-8411-68a220d797a6"),
			"username":             "john_doe",
			"email":                "john.doe@example1.com",
			"password":             "secure_password123",
			"is_account_confirmed": true,
			"is_account_deleted":   true,
		},
		map[string]string{
			"username": "john_doe",
			"password": "secure_password123",
		},
		map[string]interface{}{
			"message": "Invalid username or password",
		},
		401,
	},
}

func TestIntegration_Login(t *testing.T) {

	testServerTools := integration_test_tools.NewSetupIntegrationTestServerTools()
	testDatabaseTools := integration_test_tools.NewSetupIntegrationTestDatabaseTools()
	saltGeneratorAndPasswordHasher := servicecomponent.NewSaltGeneratorAndPasswordHasher()
	db, err := testDatabaseTools.SetupTestDB()
	if err != nil {
		t.Fatalf("failed setup database connection: %v", err)
	}
	defer db.Close()

	testServer := testServerTools.SetupTestServer(db)
	defer testServer.Close()

	for _, tc := range LoginTestCases {
		t.Run(tc.name, func(t *testing.T) {
			salt, sgaphError := saltGeneratorAndPasswordHasher.GenerateSalt()
			if sgaphError != nil {
				t.Fatalf("failed to generate salt: %v", sgaphError.ChildError)
			}
			hashedPassword := saltGeneratorAndPasswordHasher.HashPassword(tc.user["password"].(string), *salt)

			user := model.User{
				Id:                 tc.user["id"].(uuid.UUID),
				Username:           tc.user["username"].(string),
				Email:              tc.user["email"].(string),
				Password:           *hashedPassword,
				Salt:               *salt,
				RegistrationDate:   testServerTools.GetCurrentTime(),
				IsAccountConfirmed: tc.user["is_account_confirmed"].(bool),
				IsAccountDeleted:   tc.user["is_account_deleted"].(bool),
			}

			testDatabaseTools.InsertUser(db, user)

			url := fmt.Sprint("/auth/login")
			postBody, _ := json.Marshal(tc.requestBody)
			req, err := http.NewRequest(http.MethodPost, testServer.URL+url, bytes.NewBuffer(postBody))
			if err != nil {
				fmt.Printf("Error creating request: %v", err)
				return
			}

			client := &http.Client{}
			resp, err := client.Do(req)

			if err != nil {
				fmt.Printf("Error creating request: %v", err)
				return
			}

			responseBody, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("failed to read response body: %v", err)
			}

			var responseBodyMap map[string]interface{}
			err = json.Unmarshal(responseBody, &responseBodyMap)
			if err != nil {
				t.Fatalf("failed to unmarshal response body: %v", err)
			}

			if resp.StatusCode != tc.expectedStatusCode || !reflect.DeepEqual(responseBodyMap, tc.expectedResponse) {
				t.Errorf("\nscenario: %s, expected response body: %v, got response body: %v, expected status code: %v, got status code: %v\n",
					tc.name, tc.expectedResponse, responseBodyMap, tc.expectedStatusCode, resp.StatusCode)

				if resp.StatusCode != tc.expectedStatusCode {
					t.Error("Status code is not the same.")
				}

				if !reflect.DeepEqual(responseBodyMap, tc.expectedResponse) {
					t.Error("Response body is not the same.")
				}
			}

			if tc.name == "success" {
				cookies := resp.Cookies()
				var sessionCookie *http.Cookie

				for _, cookie := range cookies {
					if cookie.Name == "sessionId" {
						sessionCookie = cookie
						break
					}
				}

				if sessionCookie == nil {
					t.Errorf("Expected sessionId cookie to be set, but it was not found")
				}
			}

			defer resp.Body.Close()
			defer testDatabaseTools.TruncateUserTable(db)
		})
	}

}
