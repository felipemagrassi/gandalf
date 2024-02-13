package adapter

import (
	"testing"
	"time"
)

func TestRedisStorageBlockKey(t *testing.T) {
	redisStorage := NewRedisStorage("localhost", "6379", "", 0)
	key := "1"
	keyType := "token"

	_, err := redisStorage.GetBlockedKey(key, keyType)

	if err != KeyNotFound {
		t.Errorf("Expected nil, got %v", err)
	}

	err = redisStorage.BlockKey(key, keyType)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	blockedAt, err := redisStorage.GetBlockedKey(key, keyType)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	if blockedAt == nil {
		t.Errorf("Expected time, got %v", blockedAt)
	}

	err = redisStorage.UnblockKey(key, keyType)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	_, err = redisStorage.GetBlockedKey(key, keyType)
	if err != KeyNotFound {
		t.Errorf("Expected nil, got %v", err)
	}
}

func TestRedisStorageIncrement(t *testing.T) {
	redisStorage := NewRedisStorage("localhost", "6379", "", 0)
	key := "1"
	keyType := "token"
	redisStorage.ClearOldAccesses(key, keyType, 0*time.Microsecond)

	err := redisStorage.Increment(key, keyType)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	err = redisStorage.Increment(key, keyType)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	config, err := redisStorage.GetKeyInfo(key, keyType)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	if config.Accesses != 2 {
		t.Errorf("Expected 2, got %v", config.Accesses)
	}

	time.Sleep(1 * time.Microsecond)

	err = redisStorage.ClearOldAccesses(key, keyType, 1*time.Microsecond)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	config, err = redisStorage.GetKeyInfo(key, keyType)

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	if config.Accesses != 0 {
		t.Errorf("Expected 0, got %v", config.Accesses)
	}

}
