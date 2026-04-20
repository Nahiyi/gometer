package loader

import (
	"testing"

	"gmeter/internal/config"
)

func TestNew(t *testing.T) {
	users := []config.UserConfig{
		{Headers: map[string]string{"token": "user-1"}},
		{Headers: map[string]string{"token": "user-2"}},
	}

	l := New(users)

	if l.UserCount() != 2 {
		t.Errorf("Expected UserCount 2, got %d", l.UserCount())
	}
}

func TestGetUserConfig(t *testing.T) {
	users := []config.UserConfig{
		{Headers: map[string]string{"token": "user-token-1", "x-user-id": "1001"}},
		{Headers: map[string]string{"token": "user-token-2", "x-user-id": "1002"}},
	}

	l := New(users)

	// Thread 0 should get users[0]
	user0 := l.GetUserConfig(0)
	if user0.Headers["token"] != "user-token-1" {
		t.Errorf("Thread 0: expected 'user-token-1', got '%s'", user0.Headers["token"])
	}

	// Thread 1 should get users[1]
	user1 := l.GetUserConfig(1)
	if user1.Headers["token"] != "user-token-2" {
		t.Errorf("Thread 1: expected 'user-token-2', got '%s'", user1.Headers["token"])
	}

	// Thread 2 should wrap around to users[0]
	user2 := l.GetUserConfig(2)
	if user2.Headers["token"] != "user-token-1" {
		t.Errorf("Thread 2: expected 'user-token-1', got '%s'", user2.Headers["token"])
	}

	// Thread 3 should wrap around to users[1]
	user3 := l.GetUserConfig(3)
	if user3.Headers["token"] != "user-token-2" {
		t.Errorf("Thread 3: expected 'user-token-2', got '%s'", user3.Headers["token"])
	}
}

func TestUserCountEmpty(t *testing.T) {
	l := New([]config.UserConfig{})

	if l.UserCount() != 0 {
		t.Errorf("Expected UserCount 0, got %d", l.UserCount())
	}
}
