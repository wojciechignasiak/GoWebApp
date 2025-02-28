package service

import servicecomponent "app/internal/service_component"

type ChatService interface {
}

type chatService struct {
	us UserService
	ms MessageService
	ug servicecomponent.UuidGenerator
}

func NewChatService(us UserService, ms MessageService, ug servicecomponent.UuidGenerator) ChatService {
	return &chatService{
		us: us,
		ms: ms,
		ug: ug,
	}
}
