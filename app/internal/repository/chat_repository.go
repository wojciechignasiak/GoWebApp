package repository

import (
	apperror "app/internal/app_error"
	"app/internal/model"
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type ChatRepository interface {
}

type chatRepository struct {
	tx *sql.Tx
	db *sql.DB
}

func NewChatRepository(tx *sql.Tx, db *sql.DB) ChatRepository {
	if db == nil && tx == nil {
		panic("Missing connection and transaction in MessageRepository.")
	}
	return &chatRepository{
		tx: tx,
		db: db,
	}
}

func (cr *chatRepository) CreateChat(ctx context.Context, chat model.Chat) *apperror.AppError {
	query := `
		INSERT INTO chat (id, user_id, name)
		VALUES (?,?,?);
	`
	_, err := cr.tx.ExecContext(ctx, query, chat.Id, chat.UserId, chat.Name)
	if err != nil {
		args := fmt.Sprintf("chat: %v", chat)
		repoError := apperror.AppError{
			StatusCode:      500,
			Message:         "Database error occurred while trying to create a new chat",
			StructAndMethod: "chatRepository.CreateChat()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      &err,
		}
		return &repoError
	}
	return nil
}

func (cr *chatRepository) GetUserOwnedChats(ctx context.Context, userId uuid.UUID) (*[]model.Chat, *apperror.AppError) {
	query := `
		SELECT * FROM chat WHERE user_id = ?;
	`
	rows, err := cr.db.QueryContext(ctx, query, userId)
	if err != nil {
		args := fmt.Sprintf("userId: %v", userId)
		repoError := apperror.AppError{
			StatusCode:      500,
			Message:         "Database error occurred while trying to get chats owned by user",
			StructAndMethod: "chatRepository.GetUserOwnedChats()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      &err,
		}
		return nil, &repoError
	}
	var chats []model.Chat
	for rows.Next() {
		var chat model.Chat
		if err := rows.Scan(&chat.Id, &chat.UserId, &chat.Name); err != nil {
			args := fmt.Sprintf("userId: %v", userId)
			repoError := apperror.AppError{
				StatusCode:      500,
				Message:         "Database error occurred while trying to get chats owned by user",
				StructAndMethod: "chatRepository.GetUserOwnedChats()",
				Argument:        &args,
				ChildAppError:   nil,
				ChildError:      &err,
			}
			return nil, &repoError
		}
		chats = append(chats, chat)
	}
	if len(chats) == 0 {
		return nil, nil
	} else {
		return &chats, nil
	}
}

func (cr *chatRepository) GetChatsInWhichUserParticipate(ctx context.Context, userId uuid.UUID) (*[]model.Chat, *apperror.AppError) {
	query := `
		SELECT chat.*
		FROM chat
		JOIN chat_participant 
			ON chat.id = chat_participant.chat_id
		WHERE chat_participant.user_id = ?;
	`
	rows, err := cr.db.QueryContext(ctx, query, userId)
	if err != nil {
		args := fmt.Sprintf("userId: %v", userId)
		repoError := apperror.AppError{
			StatusCode:      500,
			Message:         "Database error occurred while trying to get chats in which user participate",
			StructAndMethod: "chatRepository.GetChatsInWhichUserParticipate()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      &err,
		}
		return nil, &repoError
	}
	var chats []model.Chat
	for rows.Next() {
		var chat model.Chat
		if err := rows.Scan(&chat.Id, &chat.UserId, &chat.Name); err != nil {
			args := fmt.Sprintf("userId: %v", userId)
			repoError := apperror.AppError{
				StatusCode:      500,
				Message:         "Database error occurred while trying to get chats in which user participate",
				StructAndMethod: "chatRepository.GetUserOwnedChats()",
				Argument:        &args,
				ChildAppError:   nil,
				ChildError:      &err,
			}
			return nil, &repoError
		}
		chats = append(chats, chat)
	}
	if len(chats) == 0 {
		return nil, nil
	} else {
		return &chats, nil
	}
}

func (cr *chatRepository) AddChatParticipant(ctx context.Context, chatParticipant model.ChatParticipant) *apperror.AppError {
	query := `
		INSERT INTO chat_participant (chat_id, user_id)
		VALUES (?,?);
	`
	_, err := cr.tx.ExecContext(ctx, query, chatParticipant.ChatId, chatParticipant.UserId)
	if err != nil {
		args := fmt.Sprintf("chatParticipant: %v", chatParticipant)
		repoError := apperror.AppError{
			StatusCode:      500,
			Message:         "Database error occurred while trying to add participant to chat",
			StructAndMethod: "chatRepository.AddChatParticipant()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      &err,
		}
		return &repoError
	}
	return nil
}

func (cr *chatRepository) RemoveChatParticipant(ctx context.Context, chatParticipant model.ChatParticipant) *apperror.AppError {
	query := `
		DELETE FROM chat_participant WHERE chat_id = ? AND user_id = ?;
	`

	_, err := cr.tx.ExecContext(ctx, query, chatParticipant.ChatId, chatParticipant.UserId)
	if err != nil {
		args := fmt.Sprintf("chatParticipant: %v", chatParticipant)
		repoError := apperror.AppError{
			StatusCode:      500,
			Message:         "Database error occurred while trying to remove participant from chat",
			StructAndMethod: "chatRepository.RemoveChatParticipant()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      &err,
		}
		return &repoError
	}
	return nil
}

// func GetChatParticipants(ctx context.Context, chatId uuid.UUID) {
// 	query := `
// 		SELECT * FROM ch;
// 	`
// 	rows, err := cr.db.QueryContext(ctx, query, userId)
// 	if err != nil {
// 		args := fmt.Sprintf("userId: %v", userId)
// 		repoError := apperror.AppError{
// 			StatusCode:      500,
// 			Message:         "Database error occurred while trying to get chats in which user participate",
// 			StructAndMethod: "chatRepository.GetChatsInWhichUserParticipate()",
// 			Argument:        &args,
// 			ChildAppError:   nil,
// 			ChildError:      &err,
// 		}
// 		return nil, &repoError
// 	}
// 	var chats []model.Chat
// 	for rows.Next() {
// 		var chat model.Chat
// 		if err := rows.Scan(&chat.Id, &chat.UserId, &chat.Name); err != nil {
// 			args := fmt.Sprintf("userId: %v", userId)
// 			repoError := apperror.AppError{
// 				StatusCode:      500,
// 				Message:         "Database error occurred while trying to get chats in which user participate",
// 				StructAndMethod: "chatRepository.GetUserOwnedChats()",
// 				Argument:        &args,
// 				ChildAppError:   nil,
// 				ChildError:      &err,
// 			}
// 			return nil, &repoError
// 		}
// 		chats = append(chats, chat)
// 	}
// 	if len(chats) == 0 {
// 		return nil, nil
// 	} else {
// 		return &chats, nil
// 	}
// }
