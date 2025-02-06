package integration

import (
	"app/internal/controller"
	controllercomponent "app/internal/controller_component"
	"app/internal/database"
	"app/internal/logs"
	"app/internal/model"
	"app/internal/server"
	"app/internal/service"
	servicecomponent "app/internal/service_component"
	"context"
	"database/sql"
	"net/http/httptest"
	"os"
	"time"
)

type SetupIntegrationTestServerTools struct{}

func NewSetupIntegrationTestServerTools() SetupIntegrationTestServerTools {
	return SetupIntegrationTestServerTools{}
}

func (sitst *SetupIntegrationTestServerTools) SetupTestServer(db *sql.DB) *httptest.Server {

	credentialsValidator := servicecomponent.NewCredentialsValidator()
	uuidGenerator := servicecomponent.NewUuidGenerator()
	saltGeneratorAndPasswordHasher := servicecomponent.NewSaltGeneratorAndPasswordHasher()

	userService := service.NewUserService(func() (database.UnitOfWork, error) {
		return database.NewUnitOfWork(db), nil
	}, uuidGenerator)
	sessionManagementService := service.NewSessionManagementService(uuidGenerator)
	authService := service.NewAuthService(userService, sessionManagementService, saltGeneratorAndPasswordHasher)

	registrationService := service.NewRegistrationService(userService, saltGeneratorAndPasswordHasher, credentialsValidator, uuidGenerator)

	logger := logs.NewLogger()

	responseHandler := controllercomponent.NewResponseHandler()

	authController := controller.NewAuthController(authService, registrationService, responseHandler, logger)
	userController := controller.NewUserController(userService, responseHandler, logger)

	testServer := httptest.NewServer(server.NewServer("", 8080, userController, authController).Handler)

	return testServer
}

func (sitst *SetupIntegrationTestServerTools) GetCurrentTime() *time.Time {
	currentTime := time.Now()
	return &currentTime
}

type SetupIntegrationTestDatabaseTools struct{}

func NewSetupIntegrationTestDatabaseTools() SetupIntegrationTestDatabaseTools {
	return SetupIntegrationTestDatabaseTools{}
}

func (sitdt *SetupIntegrationTestDatabaseTools) SetupTestDB() (*sql.DB, error) {

	dbConfig := database.DatabaseConfig{
		Host:     os.Getenv("SIT_DATABASE_HOST"),
		Port:     3306,
		User:     os.Getenv("SIT_DATABASE_USERNAME"),
		Password: os.Getenv("SIT_DATABASE_PASSWORD"),
		Database: os.Getenv("SIT_DATABASE_NAME"),
	}

	db, err := database.InitializeDatabase(dbConfig)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func (sitdt *SetupIntegrationTestDatabaseTools) InsertUser(db *sql.DB, user model.User) {
	ctx := context.Background()
	unitOfWork := database.NewUnitOfWork(db)
	unitOfWork.BeginTransaction()
	unitOfWork.UserRepository().CreateUser(ctx, user)
	unitOfWork.Commit()
}

func (sitdt *SetupIntegrationTestDatabaseTools) TruncateUserTable(db *sql.DB) error {

	query := `
		DELETE FROM user;
	`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (sitdt *SetupIntegrationTestDatabaseTools) InsertAccountConfirmation(db *sql.DB, accountConfirmation model.AccountConfirmation) {
	ctx := context.Background()
	unitOfWork := database.NewUnitOfWork(db)
	unitOfWork.BeginTransaction()
	unitOfWork.UserRepository().CreateAccountConfirmation(ctx, accountConfirmation)
	unitOfWork.Commit()
}

func (sitdt *SetupIntegrationTestDatabaseTools) TruncateAccountConfirmationTable(db *sql.DB) error {

	query := `
		DELETE FROM account_confirmation;
	`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
