package integration

import (
	"app/internal/database"
	"database/sql"
	"fmt"
	"os"
)

func SetupTestDB() (*sql.DB, error) {
	// db_port, err := strconv.Atoi(os.Getenv("SIT_DATABASE_PORT"))
	// if err != nil {
	// 	return nil, err
	// }
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

func InitializeDBInstance(db *sql.DB) error {
	db_name := os.Getenv("SIT_DATABASE_NAME")

	// Sprawdzenie, czy baza danych już istnieje.
	var dbExists bool
	err := db.QueryRow("SELECT 1 FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = ?", db_name).Scan(&dbExists)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error durning checking does exists: %w", err)
	}

	if dbExists {
		fmt.Printf("Database '%s' already exists. Skipping.\n", db_name)
		return nil
	}

	query := fmt.Sprintf("CREATE DATABASE `%s`;", db_name)
	_, err = db.Exec(query)
	if err != nil {
		return fmt.Errorf("error occured durning creating database: %w", err) // Dodajemy kontekst do błędu
	}
	return nil
}

// func InitializeDBSchema(db *sql.DB) error { // Dodajemy argument *testing.T
// 	var projectRoot string

// 	_, filename, _, ok := runtime.Caller(0)
// 	if !ok {
// 		return fmt.Errorf("nie można uzyskać informacji o pliku")
// 	}
// 	dirname := filepath.Dir(filename)
// 	projectRoot = filepath.Join(dirname, "..", "..") // Cofamy się o dwa katalogi (do korzenia projektu)

// 	path := filepath.Join(projectRoot, "internal", "database", "schema.sql")
// 	fmt.Printf("Path: %v\n", path) // Dodajemy newline dla czytelności

// 	c, err := os.ReadFile(path)
// 	if err != nil {
// 		return fmt.Errorf("błąd odczytu pliku schema.sql: %w", err)
// 	}
// 	queries := string(c)

// 	parts := strings.Split(queries, ";")
// 	for _, query := range parts {
// 		_, err = db.Exec(query)
// 		if err != nil {
// 			return fmt.Errorf("błąd wykonania zapytania SQL: %w", err)
// 		}
// 	}
// 	return nil
// }

func TruncateUserTable(db *sql.DB) error {

	query := `
		DELETE FROM user;
	`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
