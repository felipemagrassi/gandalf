package adapter

import (
	"strings"
	"testing"
	"time"
)

func TestMemoryStorageIncrementWithinOneSecond(t *testing.T) {
	memoryStorage := NewMemoryStorage()
	memoryStorage.AddKey("1", "token", 2, 1)

	memoryStorage.Increment("1", "token")
	memoryStorage.Increment("1", "token")
	timeLeft, rateLimitedError := memoryStorage.Increment("1", "token")

	if rateLimitedError != RateLimitExceeded {
		t.Errorf("Expected RateLimitExceeded, got %v", rateLimitedError)
	}

	if timeLeft == nil {
		t.Errorf("Expected timeLeft, got %v", timeLeft)
	}

	timeLeft, blockedError := memoryStorage.Increment("1", "token")
	if !strings.Contains(blockedError.Error(), "Key is blocked") {
		t.Errorf("Expected Key is blocked for, got %v", blockedError)
	}

	if timeLeft == nil {
		t.Errorf("Expected timeLeft, got %v", timeLeft)
	}
}

func TestMemoryStorageIncrementAfter1Second(t *testing.T) {
	memoryStorage := NewMemoryStorage()
	memoryStorage.AddKey("1", "token", 2, 1)

	memoryStorage.Increment("1", "token")
	memoryStorage.Increment("1", "token")
	time.Sleep(1 * time.Second)
	timeLeft, error := memoryStorage.Increment("1", "token")

	if error != nil {
		t.Errorf("Expected nil, got %v", error)
	}

	if timeLeft != nil {
		t.Errorf("Expected nil, got %v", timeLeft)
	}
}

func TestMemoryStorageIncrementWithoutAdd(t *testing.T) {
	memoryStorage := NewMemoryStorage()
	_, error := memoryStorage.Increment("1", "token")

	if error != KeyTypeNotFound {
		t.Errorf("Expected KeyTypeNotFound, got %v", error)
	}
}

func TestMemoryStorageAddDifferentRps(t *testing.T) {
	memoryStorage := NewMemoryStorage()

	memoryStorage.AddKey("1", "token", 5, 1)
	error := memoryStorage.AddKey("1", "token", 10, 1)

	if error != KeyWithDifferentConfig {
		t.Errorf("Expected KeyAlreadyExists, got %v", error)
	}

}

func TestMemoryStorageDeleteKey(t *testing.T) {
	memoryStorage := NewMemoryStorage()

	memoryStorage.AddKey("1", "token", 5, 1)
	error := memoryStorage.DeleteKey("1", "token")
	if error != nil {
		t.Errorf("Expected nil, got %v", error)
	}

	alreadyDeletedError := memoryStorage.DeleteKey("1", "token")
	if !strings.Contains(alreadyDeletedError.Error(), "Key doesn't exist") {
		t.Errorf("Expected Key doesn't exist, got %v", alreadyDeletedError)
	}
}
