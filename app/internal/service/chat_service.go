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

type ChatService interface {
	CreateChat(ctx context.Context, newChat *model.NewChat) *apperror.AppError
	AddChatParticipant(ctx context.Context, newChatParticipant *model.NewChatParticipant) *apperror.AppError
}

type chatService struct {
	uowFactory func() (database.UnitOfWork, error)
	us         UserService
	ms         MessageService
	ug         servicecomponent.UuidGenerator
}

func NewChatService(uowFactory func() (database.UnitOfWork, error), us UserService, ms MessageService, ug servicecomponent.UuidGenerator) ChatService {
	return &chatService{
		uowFactory: uowFactory,
		us:         us,
		ms:         ms,
		ug:         ug,
	}
}

func (cs *chatService) CreateChat(ctx context.Context, newChat *model.NewChat) *apperror.AppError {
	rollbackNeeded := false

	isChatNameTooLong := cs.validateNewChatNameLength(newChat.Name)
	if isChatNameTooLong {
		args := fmt.Sprintf("newChat: %v", *newChat)
		serviceError := apperror.AppError{
			StatusCode:      400,
			Message:         "chat name is too long",
			StructAndMethod: "ChatService.CreateChat()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      nil,
		}
		return &serviceError
	}

	uuid, generationError := cs.ug.GenerateUuid()
	if generationError != nil {
		args := fmt.Sprintf("newChat: %v", *newChat)
		serviceError := apperror.AppError{
			StatusCode:      generationError.StatusCode,
			Message:         "error occured while generating uuid for new chat",
			StructAndMethod: "ChatService.CreateChat()",
			Argument:        &args,
			ChildAppError:   generationError,
			ChildError:      nil,
		}
		return &serviceError
	}

	chat := cs.convertNewChatToChat(*uuid, newChat)

	uow, err := cs.uowFactory()
	if err != nil {
		args := fmt.Sprintf("newChat: %v", *newChat)
		serviceError := apperror.AppError{
			StatusCode:      500,
			Message:         "error occured while creating unit of work in chat service",
			StructAndMethod: "ChatService.CreateChat()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      &err,
		}
		return &serviceError
	}

	transactionError := uow.BeginTransaction()
	if transactionError != nil {
		args := fmt.Sprintf("newChat: %v", *newChat)
		serviceError := apperror.AppError{
			StatusCode:      transactionError.StatusCode,
			Message:         transactionError.Message,
			StructAndMethod: "ChatService.CreateChat()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      &err,
		}
		return &serviceError
	}

	repositoryError := uow.ChatRepository().CreateChat(ctx, chat)
	if repositoryError != nil {
		rollbackNeeded = true
		args := fmt.Sprintf("newChat: %v", *newChat)
		serviceError := apperror.AppError{
			StatusCode:      500,
			Message:         repositoryError.Message,
			StructAndMethod: "ChatService.CreateChat()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      &err,
		}
		return &serviceError
	}

	defer func() {
		if rollbackNeeded {
			uow.Rollback()
		}
	}()

	uow.Commit()
	return nil
}

func (cs *chatService) AddChatParticipant(ctx context.Context, newChatParticipant *model.NewChatParticipant) *apperror.AppError {
	// rollbackNeeded := false

	user, userServiceError := cs.us.GetUserById(ctx, newChatParticipant.NewParticipantId)
	if userServiceError != nil {
		args := fmt.Sprintf("newChatParticipant: %v", *newChatParticipant)
		serviceError := apperror.AppError{
			StatusCode:      500,
			Message:         userServiceError.Message,
			StructAndMethod: "ChatService.CreateAddChatParticipant()",
			Argument:        &args,
			ChildAppError:   userServiceError,
			ChildError:      nil,
		}
		return &serviceError
	}

	if user == nil {
		args := fmt.Sprintf("newChatParticipant: %v", *newChatParticipant)
		serviceError := apperror.AppError{
			StatusCode:      404,
			Message:         "user with provided id not found",
			StructAndMethod: "ChatService.AddChatParticipant()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      nil,
		}
		return &serviceError
	}

	if user.IsAccountDeleted {
		args := fmt.Sprintf("newChatParticipant: %v", *newChatParticipant)
		serviceError := apperror.AppError{
			StatusCode:      404,
			Message:         "user with provided id not found",
			StructAndMethod: "ChatService.AddChatParticipant()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      nil,
		}
		return &serviceError
	}

	return nil
}

func (cs *chatService) GetChatById(ctx context.Context, chatId uuid.UUID) (*model.Chat, *apperror.AppError) {
	uow, err := cs.uowFactory()
	if err != nil {
		args := fmt.Sprintf("chatId: %s", chatId)
		serviceError := apperror.AppError{
			StatusCode:      500,
			Message:         "Error occured while creating unit of work in chat service",
			StructAndMethod: "UserService.GetChatById()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      &err,
		}
		return nil, &serviceError
	}

	chat, repositoryError := uow.ChatRepository().GetChatById(ctx, chatId)
	if repositoryError != nil {
		args := fmt.Sprintf("chatId: %s", chatId)
		serviceError := apperror.AppError{
			StatusCode:      500,
			Message:         repositoryError.Message,
			StructAndMethod: "UserService.GetChatById()",
			Argument:        &args,
			ChildAppError:   repositoryError,
			ChildError:      nil,
		}
		return nil, &serviceError
	}

	return chat, nil
}

func (cs *chatService) validateNewChatNameLength(name string) bool {
	return len(name) < 100
}

func (cs *chatService) convertNewChatToChat(uuid uuid.UUID, newChat *model.NewChat) *model.Chat {
	return &model.Chat{
		Id:     uuid,
		UserId: newChat.UserId,
		Name:   newChat.Name,
	}
}
