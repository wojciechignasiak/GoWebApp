package repository

import (
	apperror "app/internal/app_error"
	"app/internal/model"
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user model.User) *apperror.AppError
	GetUserById(ctx context.Context, id uuid.UUID) (*model.User, *apperror.AppError)
	GetUserByEmail(ctx context.Context, email string) (*model.User, *apperror.AppError)
	GetUserByUsername(ctx context.Context, username string) (*model.User, *apperror.AppError)
}

type userRepository struct {
	tx *sql.Tx
	db *sql.DB
}

func NewUserRepository(tx *sql.Tx, db *sql.DB) *userRepository {
	if db == nil && tx == nil {
		panic("missing connection and transaction in UserRepository.")
	}
	return &userRepository{
		tx: tx,
		db: db,
	}
}

func (ur *userRepository) CreateUser(ctx context.Context, user model.User) *apperror.AppError {
	query := `
		INSERT INTO user (id, username, email, password, salt, phone_number)
		VALUES (?, ?, ?, ?, ?, ?);
	`
	_, err := ur.tx.ExecContext(ctx, query, user.Id, user.Username, user.Email, user.Password, user.Salt, user.PhoneNumber)
	if err != nil {
		repositoryError := apperror.AppError{
			StatusCode:    500,
			Message:       "database error occurred while trying to create a new user",
			ChildAppError: nil,
			ChildError:    &err,
			Logging:       true,
		}
		return &repositoryError
	}

	return nil
}

func (ur *userRepository) GetUserById(ctx context.Context, id uuid.UUID) (*model.User, *apperror.AppError) {
	query := `SELECT * FROM user WHERE id = ?;`
	row := ur.db.QueryRowContext(ctx, query, id)
	var user model.User
	err := row.Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.Salt, &user.PhoneNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			repositoryError := apperror.AppError{
				StatusCode:    500,
				Message:       "database error occurred while trying to get user by id",
				ChildAppError: nil,
				ChildError:    &err,
				Logging:       true,
			}
			return nil, &repositoryError
		}
	}

	return &user, nil
}

func (ur *userRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, *apperror.AppError) {
	query := `SELECT * FROM user WHERE email = ?;`
	row := ur.db.QueryRowContext(ctx, query, email)
	var user model.User
	err := row.Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.Salt, &user.PhoneNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			repositoryError := apperror.AppError{
				StatusCode:    500,
				Message:       "database error occurred while trying to get user by email",
				ChildAppError: nil,
				ChildError:    &err,
				Logging:       true,
			}
			return nil, &repositoryError
		}

	}

	return &user, nil
}

func (ur *userRepository) GetUserByUsername(ctx context.Context, username string) (*model.User, *apperror.AppError) {
	query := `SELECT * FROM user WHERE username = ?;`
	row := ur.db.QueryRowContext(ctx, query, username)
	var user model.User
	err := row.Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.Salt, &user.PhoneNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			repositoryError := apperror.AppError{
				StatusCode:    500,
				Message:       "database error occurred while trying to get user by username",
				ChildAppError: nil,
				ChildError:    &err,
				Logging:       true,
			}
			return nil, &repositoryError
		}
	}

	return &user, nil
}
