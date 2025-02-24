package model

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	Id      uuid.UUID `json:"id" db:"id"`
	UserId  uuid.UUID `json:"user_id" db:"user_id"`
	ChatId  uuid.UUID `json:"chat_id" db:"chat_id"`
	Content string    `json:"content" db:"content"`
	Created time.Time `json:"created" db:"created"`
}

type NewMessage struct {
	ChatId  uuid.UUID `json:"chat_id" db:"chat_id"`
	Content string    `json:"content" db:"content"`
}
