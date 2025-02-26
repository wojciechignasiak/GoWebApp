package model

import "github.com/google/uuid"

type NewChat struct {
	name string
}

type Chat struct {
	Id     uuid.UUID
	UserId uuid.UUID
	Name   string
}

type ChatParticipant struct {
	ChatId uuid.UUID
	UserId uuid.UUID
}
