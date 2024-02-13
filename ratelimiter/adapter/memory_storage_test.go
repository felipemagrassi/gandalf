package adapter

import (
	"testing"
	"time"
)

func TestMemoryStorageBlockKey(t *testing.T) {
	memoryStorage := NewMemoryStorage()
	key := "1"
	keyType := "token"

	_, err := memoryStorage.GetBlockedKey(key, keyType)

	if err != KeyTypeNotFound {
		t.Errorf("Expected nil, got %v", err)
	}

	err = memoryStorage.BlockKey(key, keyType)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	blockedAt, err := memoryStorage.GetBlockedKey(key, keyType)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	if blockedAt == nil {
		t.Errorf("Expected time, got nil")
	}

	err = memoryStorage.UnblockKey(key, keyType)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	_, err = memoryStorage.GetBlockedKey(key, keyType)
	if err != KeyNotFound {
		t.Errorf("Expected nil, got %v", err)
	}
}

func TestMemoryStorageIncrement(t *testing.T) {
	memoryStorage := NewMemoryStorage()
	key := "1"
	keyType := "token"

	err := memoryStorage.Increment(key, keyType)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	err = memoryStorage.Increment(key, keyType)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	config, err := memoryStorage.GetKeyInfo(key, keyType)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	if len(config.Accesses) != 2 {
		t.Errorf("Expected 2, got %v", len(config.Accesses))
	}

	time.Sleep(1 * time.Microsecond)

	err = memoryStorage.ClearOldAccesses(key, keyType, 1*time.Microsecond)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	config, err = memoryStorage.GetKeyInfo(key, keyType)

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	if len(config.Accesses) != 0 {
		t.Errorf("Expected 0, got %v", len(config.Accesses))
	}

}
