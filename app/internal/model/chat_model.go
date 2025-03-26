package model

import "github.com/google/uuid"

type NewChat struct {
	UserId uuid.UUID `json:"user_id" db:"user_id"`
	Name   string    `json:"name" db:"name"`
}

type Chat struct {
	Id     uuid.UUID `json:"id" db:"id"`
	UserId uuid.UUID `json:"user_id" db:"user_id"`
	Name   string    `json:"name" db:"name"`
}

type NewChatParticipant struct {
	ChatId           uuid.UUID `json:"chat_id"`
	ChatOwnerId      uuid.UUID `json:"chat_owner_id"`
	NewParticipantId uuid.UUID `json:"new_participant_id" db:"new_participant_id"`
}

type ChatParticipant struct {
	ChatId uuid.UUID `json:"chat_id" db:"chat_id"`
	UserId uuid.UUID `json:"user_id" db:"user_id"`
}
