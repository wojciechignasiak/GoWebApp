package service

import (
	apperror "app/internal/app_error"
	"app/internal/database"
	"app/internal/model"
	servicecomponent "app/internal/service_component"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type MessageService interface {
	CreateMessage(ctx context.Context, message model.Message) *apperror.AppError
	GetMessagesPaginated(ctx context.Context, chatId uuid.UUID, pageNumber int) (*[]model.Message, *apperror.AppError)
}

type messageService struct {
	uowFactory func() (database.UnitOfWork, error)
	ug         servicecomponent.UuidGenerator
}

func NewMessageService(uowFactory func() (database.UnitOfWork, error), ug servicecomponent.UuidGenerator) MessageService {
	return &messageService{
		uowFactory: uowFactory,
		ug:         ug,
	}
}

func (ms *messageService) CreateMessage(ctx context.Context, message model.Message) *apperror.AppError {
	uow, err := ms.uowFactory()
	if err != nil {
		args := fmt.Sprintf("message: %v", message)
		serviceError := apperror.AppError{
			StatusCode:      500,
			Message:         "error occured while creating unit of work in message service",
			StructAndMethod: "MessageService.CreateMessage()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      &err,
		}
		return &serviceError
	}
	uowError := uow.BeginTransaction()
	if uowError != nil {
		args := fmt.Sprintf("message: %v", message)
		serviceError := apperror.AppError{
			StatusCode:      uowError.StatusCode,
			Message:         uowError.Message,
			StructAndMethod: "UserService.CreateMessage()",
			Argument:        &args,
			ChildAppError:   uowError,
			ChildError:      uowError.ChildError,
		}
		return &serviceError
	}
	repoError := uow.MessageRepository().CreateMessage(ctx, message)
	if repoError != nil {
		args := fmt.Sprintf("message: %v", message)
		serviceError := apperror.AppError{
			StatusCode:      repoError.StatusCode,
			Message:         repoError.Message,
			StructAndMethod: "UserService.CreateMessage()",
			Argument:        &args,
			ChildAppError:   repoError,
			ChildError:      nil,
		}
		return &serviceError
	}
	uow.Commit()
	return nil
}

func (ms *messageService) GetMessagesPaginated(ctx context.Context, chatId uuid.UUID, pageNumber int) (*[]model.Message, *apperror.AppError) {
	uow, err := ms.uowFactory()
	if err != nil {
		args := fmt.Sprintf("chatId: %v, pageNumber: %d", chatId, pageNumber)
		serviceError := apperror.AppError{
			StatusCode:      500,
			Message:         "error occured while creating unit of work in message service",
			StructAndMethod: "MessageService.GetMessagesPaginated()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      &err,
		}
		return nil, &serviceError
	}
	messages, repoError := uow.MessageRepository().GetMessagesPaginated(ctx, chatId, pageNumber)
	if repoError != nil {
		args := fmt.Sprintf("chatId: %v, pageNumber: %d", chatId, pageNumber)
		serviceError := apperror.AppError{
			StatusCode:      repoError.StatusCode,
			Message:         repoError.Message,
			StructAndMethod: "UserService.GetMessagesPaginated()",
			Argument:        &args,
			ChildAppError:   repoError,
			ChildError:      nil,
		}
		return nil, &serviceError
	}
	return messages, nil
}
