package integration_test

import (
	"app/internal/controller"
	controllercomponent "app/internal/controller_component"
	"app/internal/database"
	integration "app/internal/integration_test"
	"app/internal/logs"
	"app/internal/model"
	"app/internal/server"
	"app/internal/service"
	servicecomponent "app/internal/service_component"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
)

var currentTime time.Time = time.Now()

func insertUserAndAccountConfirmation(db *sql.DB, user model.User, accountConfirmation model.AccountConfirmation) {
	ctx := context.Background()
	unitOfWork := database.NewUnitOfWork(db)
	unitOfWork.BeginTransaction()
	unitOfWork.UserRepository().CreateUser(ctx, user)
	unitOfWork.UserRepository().CreateAccountConfirmation(ctx, accountConfirmation)
	unitOfWork.Commit()
}

func IsAccountConfirmed(t *testing.T, db *sql.DB, id uuid.UUID) bool {
	ctx := context.Background()
	unitOfWork := database.NewUnitOfWork(db)
	user, dbErr := unitOfWork.UserRepository().GetUserById(ctx, id)

	if dbErr != nil {
		t.Fatalf("failed to fetch user from database")
	}
	if user.IsAccountConfirmed == true {
		return true
	}

	return false
}

var confirmAccountTestCases = []struct {
	name                string
	user                model.User
	accountConfirmation model.AccountConfirmation
	confirmationCode    string
	securityCode        string
	expectedResponse    map[string]interface{}
	expectedStatusCode  int
	isUserConfirmedInDB bool
}{
	{
		name: "success",
		user: model.User{
			Id:                 uuid.MustParse("01945a93-b693-7713-aae0-b0ca59ef4ffd"),
			Username:           "john_doe",
			Email:              "john_doe@domain.com",
			Password:           []byte("secure_password123"),
			Salt:               []byte("G91ApOHlFuMCWuat"),
			RegistrationDate:   &currentTime,
			IsAccountConfirmed: false,
			IsAccountDeleted:   false,
		},
		accountConfirmation: model.AccountConfirmation{
			UserId:           uuid.MustParse("01945a93-b693-7713-aae0-b0ca59ef4ffd"),
			ConfirmationCode: uuid.MustParse("01945af2-19d5-7e5f-8546-00a94aeaea22"),
			SecurityCode:     "12avc3",
		},
		confirmationCode: "01945af2-19d5-7e5f-8546-00a94aeaea22",
		securityCode:     "12avc3",
		expectedResponse: map[string]interface{}{
			"message": "account confirmed successfully",
		},
		expectedStatusCode:  200,
		isUserConfirmedInDB: true,
	},
	{
		name: "wrong confirmation code",
		user: model.User{
			Id:                 uuid.MustParse("01945a93-b693-7713-aae0-b0ca59ef4ffd"),
			Username:           "john_doe",
			Email:              "john_doe@domain.com",
			Password:           []byte("secure_password123"),
			Salt:               []byte("G91ApOHlFuMCWuat"),
			RegistrationDate:   &currentTime,
			IsAccountConfirmed: false,
			IsAccountDeleted:   false,
		},
		accountConfirmation: model.AccountConfirmation{
			UserId:           uuid.MustParse("01945a93-b693-7713-aae0-b0ca59ef4ffd"),
			ConfirmationCode: uuid.MustParse("01945af2-19d5-7e5f-8546-00a94aeaea22"),
			SecurityCode:     "12avc3",
		},
		confirmationCode: "01945b26-99a4-7761-93b7-4667321d889a",
		securityCode:     "12avc3",
		expectedResponse: map[string]interface{}{
			"message": "content not found",
		},
		expectedStatusCode:  404,
		isUserConfirmedInDB: false,
	},
	{
		name: "wrong confirmation security code",
		user: model.User{
			Id:                 uuid.MustParse("01945a93-b693-7713-aae0-b0ca59ef4ffd"),
			Username:           "john_doe",
			Email:              "john_doe@domain.com",
			Password:           []byte("secure_password123"),
			Salt:               []byte("G91ApOHlFuMCWuat"),
			RegistrationDate:   &currentTime,
			IsAccountConfirmed: false,
			IsAccountDeleted:   false,
		},
		accountConfirmation: model.AccountConfirmation{
			UserId:           uuid.MustParse("01945a93-b693-7713-aae0-b0ca59ef4ffd"),
			ConfirmationCode: uuid.MustParse("01945af2-19d5-7e5f-8546-00a94aeaea22"),
			SecurityCode:     "12avc3",
		},
		confirmationCode: "01945af2-19d5-7e5f-8546-00a94aeaea22",
		securityCode:     "111111",
		expectedResponse: map[string]interface{}{
			"message": "content not found",
		},
		expectedStatusCode:  404,
		isUserConfirmedInDB: false,
	},
	{
		name: "account deleted",
		user: model.User{
			Id:                 uuid.MustParse("01945a93-b693-7713-aae0-b0ca59ef4ffd"),
			Username:           "john_doe",
			Email:              "john_doe@domain.com",
			Password:           []byte("secure_password123"),
			Salt:               []byte("G91ApOHlFuMCWuat"),
			RegistrationDate:   &currentTime,
			IsAccountConfirmed: true,
			IsAccountDeleted:   true,
		},
		accountConfirmation: model.AccountConfirmation{
			UserId:           uuid.MustParse("01945a93-b693-7713-aae0-b0ca59ef4ffd"),
			ConfirmationCode: uuid.MustParse("01945af2-19d5-7e5f-8546-00a94aeaea22"),
			SecurityCode:     "12avc3",
		},
		confirmationCode: "01945af2-19d5-7e5f-8546-00a94aeaea22",
		securityCode:     "12avc3",
		expectedResponse: map[string]interface{}{
			"message": "content not found",
		},
		expectedStatusCode:  404,
		isUserConfirmedInDB: true,
	},
	{
		name: "account arleady confirmed",
		user: model.User{
			Id:                 uuid.MustParse("01945a93-b693-7713-aae0-b0ca59ef4ffd"),
			Username:           "john_doe",
			Email:              "john_doe@domain.com",
			Password:           []byte("secure_password123"),
			Salt:               []byte("G91ApOHlFuMCWuat"),
			RegistrationDate:   &currentTime,
			IsAccountConfirmed: true,
			IsAccountDeleted:   false,
		},
		accountConfirmation: model.AccountConfirmation{
			UserId:           uuid.MustParse("01945a93-b693-7713-aae0-b0ca59ef4ffd"),
			ConfirmationCode: uuid.MustParse("01945af2-19d5-7e5f-8546-00a94aeaea22"),
			SecurityCode:     "12avc3",
		},
		confirmationCode: "01945af2-19d5-7e5f-8546-00a94aeaea22",
		securityCode:     "12avc3",
		expectedResponse: map[string]interface{}{
			"message": "account already confirmed",
		},
		expectedStatusCode:  200,
		isUserConfirmedInDB: true,
	},
}

func TestIntegration_ConfirmAccount(t *testing.T) {
	db, err := integration.SetupTestDB()
	if err != nil {
		t.Fatalf("failed setup database connection: %v", err)
	}

	defer db.Close()

	userService := service.NewUserService(func() (database.UnitOfWork, error) {
		return database.NewUnitOfWork(db), nil
	}, servicecomponent.NewUuidGenerator())

	logger := logs.NewLogger()
	responseHandler := controllercomponent.NewResponseHandler()
	userController := controller.NewUserController(userService, responseHandler, logger)

	testServer := httptest.NewServer(server.NewServer("", 8080, userController, nil).Handler)
	defer testServer.Close()

	for _, tc := range confirmAccountTestCases {
		t.Run(tc.name, func(t *testing.T) {

			insertUserAndAccountConfirmation(db, tc.user, tc.accountConfirmation)
			url := fmt.Sprintf("/user/confirm-account/%s/%s", tc.confirmationCode, tc.securityCode)

			req, err := http.NewRequest(http.MethodPut, testServer.URL+url, nil)
			if err != nil {
				fmt.Println("Error creating request:", err)
				return
			}

			client := &http.Client{}
			resp, err := client.Do(req)

			if err != nil {
				fmt.Println("Error creating request:", err)
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

			if resp.StatusCode != tc.expectedStatusCode || !reflect.DeepEqual(responseBodyMap, tc.expectedResponse) || tc.isUserConfirmedInDB != IsAccountConfirmed(t, db, tc.user.Id) {
				t.Errorf("\nscenario: %s, expected response body: %v, got response body: %v, expected status code: %v, got status code: %v\n",
					tc.name, tc.expectedResponse, responseBodyMap, tc.expectedStatusCode, resp.StatusCode)
			}

			defer resp.Body.Close()
			defer integration.TruncateAccountConfirmationTable(db)
			defer integration.TruncateUserTable(db)

		})
	}

}
