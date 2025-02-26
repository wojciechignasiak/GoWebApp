package repository

import (
	apperror "app/internal/app_error"
	"app/internal/model"
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type MessageRepository interface {
	CreateMessage(ctx context.Context, message model.Message) *apperror.AppError
	GetMessagesPaginated(ctx context.Context, chatId uuid.UUID, pageNumber int) (*[]model.Message, *apperror.AppError)
}

type messageRepository struct {
	tx *sql.Tx
	db *sql.DB
}

func NewMessageRepository(tx *sql.Tx, db *sql.DB) MessageRepository {
	if db == nil && tx == nil {
		panic("Missing connection and transaction in MessageRepository.")
	}
	return &messageRepository{
		tx: tx,
		db: db,
	}
}

func (mr *messageRepository) CreateMessage(ctx context.Context, message model.Message) *apperror.AppError {
	query := `
		INSERT INTO message (id, chat_id, user_id, content, created)
		VALUES (?, ?, ?, ?, ?);
	`
	_, err := mr.tx.ExecContext(ctx, query, message.Id, message.ChatId, message.UserId, message.Content, message.Created)
	if err != nil {
		args := fmt.Sprintf("message: %v", message)
		repositoryError := apperror.AppError{
			StatusCode:      500,
			Message:         "Database error occurred while trying to create a new message",
			StructAndMethod: "messageRepository.CreateMessage()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      &err,
		}
		return &repositoryError
	}
	return nil
}

func (mr *messageRepository) GetMessagesPaginated(ctx context.Context, chatId uuid.UUID, pageNumber int) (*[]model.Message, *apperror.AppError) {
	query := `
		SELECT * FROM message WHERE chat_id = ? ORDER BY created DESC LIMIT 10 OFFSET ?;
	`
	offset := (int(pageNumber) - 1) * 10
	if pageNumber == 0 {
		offset = 0
	}
	rows, err := mr.db.QueryContext(ctx, query, chatId, offset)
	if err != nil {
		args := fmt.Sprintf("chatId: %v, pageNumber: %d", chatId, pageNumber)
		repositoryError := apperror.AppError{
			StatusCode:      500,
			Message:         "Database error occurred while trying to get messages",
			StructAndMethod: "messageRepository.GetMessagesPaginated()",
			Argument:        &args,
			ChildAppError:   nil,
			ChildError:      &err,
		}
		return nil, &repositoryError
	}
	var messages []model.Message
	for rows.Next() {
		var message model.Message
		if err := rows.Scan(&message.Id, &message.UserId, &message.ChatId, &message.Content, &message.Created); err != nil {
			args := fmt.Sprintf("chatId: %v, pageNumber: %d", chatId, pageNumber)
			repositoryError := apperror.AppError{
				StatusCode:      500,
				Message:         "Database error occurred while trying to get messages",
				StructAndMethod: "messageRepository.GetMessagesPaginated()",
				Argument:        &args,
				ChildAppError:   nil,
				ChildError:      &err,
			}
			return nil, &repositoryError
		}
		messages = append(messages, message)
	}
	if len(messages) == 0 {
		return nil, nil
	} else {
		return &messages, nil
	}
}
