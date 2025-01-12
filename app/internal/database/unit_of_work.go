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
			StatusCode:      500,
			Message:         "Transaction already started",
			StructAndMethod: "unitOfWork.BeginTransaction()",
			Argument:        nil,
			ChildAppError:   nil,
			ChildError:      nil,
		}
		return &uowError
	}
	var err error
	u.tx, err = u.db.Begin()
	if err != nil {
		uowError := apperror.AppError{
			StatusCode:      500,
			Message:         "Failed to begin transaction",
			StructAndMethod: "unitOfWork.BeginTransaction()",
			Argument:        nil,
			ChildAppError:   nil,
			ChildError:      &err,
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
			Message:       "Transaction already closed",
			ChildAppError: nil,
			ChildError:    nil,
		}
		return &uowError
	}
	err := u.tx.Commit()
	if err != nil {
		uowError := apperror.AppError{
			StatusCode:      500,
			Message:         "Error occured while commiting changes to database",
			StructAndMethod: "unitOfWork.Commit()",
			Argument:        nil,
			ChildAppError:   nil,
			ChildError:      &err,
		}
		return &uowError
	}
	u.tx = nil
	return nil
}

func (u *unitOfWork) Rollback() *apperror.AppError {
	if u.tx == nil {
		uowError := apperror.AppError{
			StatusCode:      500,
			Message:         "Transaction already closed",
			StructAndMethod: "unitOfWork.BeginRollback()",
			Argument:        nil,
			ChildAppError:   nil,
			ChildError:      nil,
		}
		return &uowError
	}
	err := u.tx.Rollback()
	if err != nil {
		uowError := apperror.AppError{
			StatusCode:      500,
			Message:         "Error occured while rolling back changes",
			StructAndMethod: "unitOfWork.BeginRollback()",
			Argument:        nil,
			ChildAppError:   nil,
			ChildError:      &err,
		}
		return &uowError
	}
	u.tx = nil
	return nil
}
