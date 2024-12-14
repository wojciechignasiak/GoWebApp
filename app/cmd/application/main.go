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
)

func main() {
	dbConfig := database.DatabaseConfig{
		Host:     "mysql",
		Port:     3306,
		User:     "melkey",
		Password: "password1234",
		Database: "blueprint",
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
