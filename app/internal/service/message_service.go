package service

import (
	apperror "app/internal/app_error"
	"app/internal/database"
	"app/internal/model"
	servicecomponent "app/internal/service_component"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type MessageService interface {
	CreateMessage(ctx context.Context, message model.NewMessage) *apperror.AppError
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

func (ms *messageService) CreateMessage(ctx context.Context, newMessage model.NewMessage) *apperror.AppError {
	isMessageContentNotTooLong := ms.validateNewMessageLength(newMessage.Content)
	if isMessageContentNotTooLong != false {
		args := fmt.Sprintf("newMessage: %v", newMessage)
		serviceError := apperror.AppError{
			StatusCode:      400,
			Message:         "message is too long",
			StructAndMethod: "MessageService.CreateMessage()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      nil,
		}
		return &serviceError
	}

	message, convertionError := ms.convertNewMessageToMessage(newMessage)
	if convertionError != nil {
		args := fmt.Sprintf("newMessage: %v", newMessage)
		serviceError := apperror.AppError{
			StatusCode:      500,
			Message:         "error occurred durning message conversion",
			StructAndMethod: "MessageService.CreateMessage()",
			Argument:        &args,
			ChildAppError:   convertionError,
			ChildError:      nil,
		}
		return &serviceError
	}
	uow, err := ms.uowFactory()
	if err != nil {
		args := fmt.Sprintf("newMessage: %v", newMessage)
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
		args := fmt.Sprintf("newMessage: %v", newMessage)
		serviceError := apperror.AppError{
			StatusCode:      uowError.StatusCode,
			Message:         uowError.Message,
			StructAndMethod: "UserService.CreateMessage()",
			Argument:        &args,
			ChildAppError:   uowError,
			ChildError:      nil,
		}
		return &serviceError
	}
	repoError := uow.MessageRepository().CreateMessage(ctx, *message)
	if repoError != nil {
		args := fmt.Sprintf("newMessage: %v", newMessage)
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
	uowError = uow.Commit()
	if uowError != nil {
		args := fmt.Sprintf("newMessage: %v", newMessage)
		serviceError := apperror.AppError{
			StatusCode:      uowError.StatusCode,
			Message:         uowError.Message,
			StructAndMethod: "UserService.CreateMessage()",
			Argument:        &args,
			ChildAppError:   uowError,
			ChildError:      nil,
		}
		return &serviceError
	}
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

func (ms *messageService) validateNewMessageLength(messageContent string) bool {
	if len(messageContent) > 500 {
		return false
	}
	return true
}

func (ms *messageService) convertNewMessageToMessage(newMessage model.NewMessage) (*model.Message, *apperror.AppError) {
	uuid, generationError := ms.ug.GenerateUuid()

	if generationError != nil {
		args := fmt.Sprintf("newMessage: %v", newMessage)
		serviceError := apperror.AppError{
			StatusCode:      generationError.StatusCode,
			Message:         generationError.Message,
			StructAndMethod: "UserService.convertNewMessageToMessage()",
			Argument:        &args,
			ChildAppError:   generationError,
			ChildError:      nil,
		}
		return nil, &serviceError
	}

	currentTime := time.Now()
	message := model.Message{
		Id:      *uuid,
		UserId:  newMessage.UserId,
		ChatId:  newMessage.ChatId,
		Content: newMessage.Content,
		Created: currentTime,
	}
	return &message, nil
}
