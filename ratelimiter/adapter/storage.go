package adapter

import "time"

type Storage interface {
	AddKey(key string, keyType string) error
	GetKeyInfo(key string, keyType string) (*StorageConfig, error)
	GetBlockedKey(key string, keyType string) (*time.Time, error)
	BlockKey(key string, keyType string) error
	UnblockKey(key string, keyType string) error
	Increment(key string, keyType string) error
}

type StorageConfig struct {
	key      string
	keyType  string
	accesses []time.Time
}
