package adapter

import (
	"testing"
)

func TestMemoryStorageBlockKey(t *testing.T) {
	memoryStorage := NewMemoryStorage()
	key := "1"
	keyType := "token"

	memoryStorage.AddKey(key, keyType)

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

func TestMemoryStorageKey(t *testing.T) {
	memoryStorage := NewMemoryStorage()
	key := "1"
	keyType := "token"

	keyInfo, err := memoryStorage.GetKeyInfo(key, keyType)
	if err != KeyTypeNotFound {
		t.Errorf("Expected nil, got %v", err)
	}

	if keyInfo != nil {
		t.Errorf("Expected nil, got %v", keyInfo)
	}

	err = memoryStorage.AddKey(key, keyType)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	keyInfo, err = memoryStorage.GetKeyInfo(key, keyType)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	if keyInfo == nil {
		t.Errorf("Expected keyInfo, got nil")
	}

	err = memoryStorage.AddKey(key, keyType)
	if err != KeyAlreadyRegistered {
		t.Errorf("Expected KeyAlreadyExists, got %v", err)
	}
}

func TestMemoryStorageIncrement(t *testing.T) {
	memoryStorage := NewMemoryStorage()
	key := "1"
	keyType := "token"

	err := memoryStorage.Increment(key, keyType)
	if err != KeyTypeNotFound {
		t.Errorf("Expected nil, got %v", err)
	}

	err = memoryStorage.AddKey(key, keyType)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	err = memoryStorage.Increment(key, keyType)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
}
