package database

import (
	apperror "app/internal/app_error"
	"app/internal/repository"
	"database/sql"
	"sync"
)

type UnitOfWork interface {
	UserRepository() repository.UserRepository
	Commit() *apperror.AppError
	Rollback() *apperror.AppError
	BeginTransaction() *apperror.AppError
}

type unitOfWork struct {
	db             *sql.DB
	tx             *sql.Tx
	userRepository repository.UserRepository
	repoInitOnce   sync.Once
}

func NewUnitOfWork(db *sql.DB) UnitOfWork {
	return &unitOfWork{
		db: db,
	}
}

func (u *unitOfWork) BeginTransaction() *apperror.AppError {
	if u.tx != nil {
		uowError := apperror.AppError{
			StatusCode:    500,
			Message:       "transaction already started",
			ChildAppError: nil,
			ChildError:    nil,
			Logging:       false,
		}
		return &uowError
	}
	var err error
	u.tx, err = u.db.Begin()
	if err != nil {
		uowError := apperror.AppError{
			StatusCode:    500,
			Message:       "failed to begin transaction",
			ChildAppError: nil,
			ChildError:    &err,
			Logging:       true,
		}
		return &uowError
	}
	return nil
}

func (u *unitOfWork) UserRepository() repository.UserRepository {
	u.repoInitOnce.Do(func() {
		u.userRepository = repository.NewUserRepository(u.tx, u.db)
	})
	return u.userRepository
}

func (u *unitOfWork) Commit() *apperror.AppError {
	if u.tx == nil {
		uowError := apperror.AppError{
			StatusCode:    500,
			Message:       "transaction already closed",
			ChildAppError: nil,
			ChildError:    nil,
			Logging:       false,
		}
		return &uowError
	}
	err := u.tx.Commit()
	if err != nil {
		uowError := apperror.AppError{
			StatusCode:    500,
			Message:       "error occured while commiting changes to database",
			ChildAppError: nil,
			ChildError:    &err,
			Logging:       true,
		}
		return &uowError
	}
	u.tx = nil
	return nil
}

func (u *unitOfWork) Rollback() *apperror.AppError {
	if u.tx == nil {
		uowError := apperror.AppError{
			StatusCode:    500,
			Message:       "transaction already closed",
			ChildAppError: nil,
			ChildError:    nil,
			Logging:       false,
		}
		return &uowError
	}
	err := u.tx.Rollback()
	if err != nil {
		uowError := apperror.AppError{
			StatusCode:    500,
			Message:       "error occured while rolling back changes",
			ChildAppError: nil,
			ChildError:    &err,
			Logging:       true,
		}
		return &uowError
	}
	u.tx = nil
	return nil
}
