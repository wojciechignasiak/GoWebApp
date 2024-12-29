package main

import (
	"app/internal/controller"
	controllercomponent "app/internal/controller_component"
	"app/internal/database"
	"app/internal/logs"
	"app/internal/server"
	"app/internal/service"
	servicecomponent "app/internal/service_component"
	"fmt"
	"os"
	"strconv"
)

func main() {
	db_port, err := strconv.Atoi(os.Getenv("DATABASE_PORT"))
	if err != nil {
		fmt.Println("Wrong database port: ", err)
	}
	dbConfig := database.DatabaseConfig{
		Host:     os.Getenv("DATABASE_HOST"),
		Port:     db_port,
		User:     os.Getenv("DATABASE_USERNAME"),
		Password: os.Getenv("DATABASE_PASSWORD"),
		Database: os.Getenv("DATABASE_NAME"),
	}
	db, err := database.InitializeDatabase(dbConfig)

	if err != nil {
		fmt.Println("Failed to initialize database: ", err)
	}
	commonTools := servicecomponent.NewCommonTools()
	uowFactory := func() (database.UnitOfWork, error) {
		return database.NewUnitOfWork(db), err
	}
	userService := service.NewUserService(uowFactory, commonTools)

	logger := logs.NewLogger()
	responseHandler := controllercomponent.NewResponseHandler()
	userController := controller.NewUserController(userService, responseHandler, logger)
	server := server.NewServer(
		"",
		80,
		userController,
	)

	server.ListenAndServe()

}
