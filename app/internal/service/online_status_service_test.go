package service

import (
	"testing"

	"github.com/google/uuid"
)

func TestSetUserStatusToOffline(t *testing.T) {
	onlineStatusService := onlineStatusService{
		usersOnline: make(map[uuid.UUID]struct{}),
	}
	userId, _ := uuid.NewRandom()

	onlineStatusService.usersOnline[userId] = struct{}{}
	onlineStatusService.SetUserStatusToOffline(userId)
	_, exists := onlineStatusService.usersOnline[userId]

	if exists != false {
		t.Errorf("expected: false, got: %v", exists)
	}
}

func TestSetUserStatusToOnline(t *testing.T) {
	onlineStatusService := onlineStatusService{
		usersOnline: make(map[uuid.UUID]struct{}),
	}

	userId, _ := uuid.NewRandom()

	onlineStatusService.SetUserStatusToOnline(userId)
	_, exists := onlineStatusService.usersOnline[userId]

	if exists != true {
		t.Errorf("expected: true, got: %v", exists)
	}
}

func TestIsUserOnline(t *testing.T) {
	onlineStatusService := onlineStatusService{
		usersOnline: make(map[uuid.UUID]struct{}),
	}
	userId, _ := uuid.NewRandom()
	onlineStatusService.usersOnline[userId] = struct{}{}
	exists := onlineStatusService.IsUserOnline(userId)

	if exists != true {
		t.Errorf("scenario: user online, expected: true, got: %v", exists)
	}

	delete(onlineStatusService.usersOnline, userId)
	exists = onlineStatusService.IsUserOnline(userId)

	if exists != false {
		t.Errorf("scenario: user offline, expected: false, got: %v", exists)
	}
}
