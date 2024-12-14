package repository

import (
	apperror "app/internal/app_error"
	"app/internal/model"
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user model.User) *apperror.AppError
	GetUserById(ctx context.Context, id uuid.UUID) (*model.User, *apperror.AppError)
	GetUserByEmail(ctx context.Context, email string) (*model.User, *apperror.AppError)
	GetUserByUsername(ctx context.Context, username string) (*model.User, *apperror.AppError)
	CreateAccountConfirmation(ctx context.Context, accountConfirmation model.AccountConfirmation) *apperror.AppError
	GetAccountConfirmationByConfirmationCode(ctx context.Context, confirmationCode uuid.UUID) (*model.AccountConfirmation, *apperror.AppError)
	ConfirmUserAccount(ctx context.Context, userId uuid.UUID) *apperror.AppError
}

type userRepository struct {
	tx *sql.Tx
	db *sql.DB
}

func NewUserRepository(tx *sql.Tx, db *sql.DB) *userRepository {
	if db == nil && tx == nil {
		panic("Missing connection and transaction in UserRepository.")
	}
	return &userRepository{
		tx: tx,
		db: db,
	}
}

func (ur *userRepository) CreateUser(ctx context.Context, user model.User) *apperror.AppError {
	query := `
		INSERT INTO user (id, username, email, password, salt)
		VALUES (?, ?, ?, ?, ?);
	`
	_, err := ur.tx.ExecContext(ctx, query, user.Id, user.Username, user.Email, user.Password, user.Salt)
	if err != nil {
		args := fmt.Sprintf("user: %v", user)
		repositoryError := apperror.AppError{
			StatusCode:      500,
			Message:         "Database error occurred while trying to create a new user",
			StructAndMethod: "userRepository.CreateUser()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      &err,
		}
		return &repositoryError
	}
	return nil
}

func (ur *userRepository) GetUserById(ctx context.Context, id uuid.UUID) (*model.User, *apperror.AppError) {
	query := `SELECT * FROM user WHERE id = ?;`
	row := ur.db.QueryRowContext(ctx, query, id)
	var user model.User
	err := row.Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.Salt, &user.RegistrationDate, &user.IsAccountConfirmed, &user.IsAccountDeleted)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			args := fmt.Sprintf("id: %s", id)
			repositoryError := apperror.AppError{
				StatusCode:      500,
				Message:         "Database error occurred while trying to get user by id",
				StructAndMethod: "userRepository.GetUserById()",
				Argument:        &args,
				ChildAppError:   nil,
				ChildError:      &err,
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
	err := row.Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.Salt, &user.RegistrationDate, &user.IsAccountConfirmed, &user.IsAccountDeleted)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			args := fmt.Sprintf("email: %s", email)
			repositoryError := apperror.AppError{
				StatusCode:      500,
				Message:         "Database error occurred while trying to get user by email",
				StructAndMethod: "userRepository.GetUserByEmail()",
				Argument:        &args,
				ChildAppError:   nil,
				ChildError:      &err,
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
	err := row.Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.Salt, &user.RegistrationDate, &user.IsAccountConfirmed, &user.IsAccountDeleted)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			args := fmt.Sprintf("username: %s", username)
			repositoryError := apperror.AppError{
				StatusCode:      500,
				Message:         "Database error occurred while trying to get user by username",
				StructAndMethod: "userRepository.GetUserByUsername()",
				Argument:        &args,
				ChildAppError:   nil,
				ChildError:      &err,
			}
			return nil, &repositoryError
		}
	}
	return &user, nil
}

func (ur *userRepository) CreateAccountConfirmation(ctx context.Context, accountConfirmation model.AccountConfirmation) *apperror.AppError {
	query := `
		INSERT INTO account_confirmation (user_id, confirmation_code, security_code) 
		VALUES (?, ?, ?)
	`
	_, err := ur.tx.ExecContext(ctx, query, &accountConfirmation.UserId, &accountConfirmation.ConfirmationCode, &accountConfirmation.SecurityCode)

	if err != nil {
		args := fmt.Sprintf("accountConfirmation: %v", accountConfirmation)
		repositoryError := apperror.AppError{
			StatusCode:      500,
			Message:         "Database error occurred while trying to create a account confirmation entry",
			StructAndMethod: "userRepository.CreateAccountConfirmation()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      &err,
		}
		return &repositoryError
	}
	return nil
}

func (ur *userRepository) GetAccountConfirmationByConfirmationCode(ctx context.Context, confirmationCode uuid.UUID) (*model.AccountConfirmation, *apperror.AppError) {
	query := `SELECT * FROM account_confirmation WHERE confirmation_code = ?;`
	row := ur.db.QueryRowContext(ctx, query, confirmationCode)
	var accountConfirmation model.AccountConfirmation
	err := row.Scan(&accountConfirmation.UserId, &accountConfirmation.ConfirmationCode, &accountConfirmation.SecurityCode)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			args := fmt.Sprintf("confirmationCode: %s", confirmationCode)
			repositoryError := apperror.AppError{
				StatusCode:      500,
				Message:         "Database error occurred while trying to get account confirmation by confirmation code",
				StructAndMethod: "userRepository.GetAccountConfirmationByConfirmationCode()",
				Argument:        &args,
				ChildAppError:   nil,
				ChildError:      &err,
			}
			return nil, &repositoryError
		}
	}
	return &accountConfirmation, nil
}

func (ur *userRepository) ConfirmUserAccount(ctx context.Context, userId uuid.UUID) *apperror.AppError {
	query := `
		UPDATE user
		SET is_account_confirmed = ?
		WHERE id = ?;
	`
	_, err := ur.tx.ExecContext(ctx, query, true, userId)
	if err != nil {
		args := fmt.Sprintf("userId: %v", userId)
		repositoryError := apperror.AppError{
			StatusCode:      500,
			Message:         "Database error occurred while trying to confirm user account",
			StructAndMethod: "userRepository.ConfirmUserAccount()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      &err,
		}
		return &repositoryError
	}
	return nil
}
