package service

import (
	"sync"

	"github.com/google/uuid"
)

type OnlineStatusService interface {
	SetUserStatusToOffline(userId uuid.UUID)
	SetUserStatusToOnline(userId uuid.UUID)
	IsUserOnline(userId uuid.UUID) bool
}

type onlineStatusService struct {
	usersOnline map[uuid.UUID]struct{}
	mu          sync.RWMutex
}

func NewOnlineStatusService() OnlineStatusService {
	return &onlineStatusService{}
}

func (oss *onlineStatusService) SetUserStatusToOffline(userId uuid.UUID) {
	oss.mu.Lock()
	defer oss.mu.Unlock()
	delete(oss.usersOnline, userId)
}

func (oss *onlineStatusService) SetUserStatusToOnline(userId uuid.UUID) {
	oss.mu.RLock()
	defer oss.mu.RUnlock()
	oss.usersOnline[userId] = struct{}{}
}

func (oss *onlineStatusService) IsUserOnline(userId uuid.UUID) bool {
	oss.mu.RLock()
	defer oss.mu.RUnlock()
	_, exists := oss.usersOnline[userId]
	return exists
}
