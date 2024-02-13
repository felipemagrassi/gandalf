package adapter

import "time"

type Storage interface {
	GetKeyInfo(key string, keyType string) (*StorageConfig, error)
	GetBlockedKey(key string, keyType string) (*time.Time, error)
	BlockKey(key string, keyType string) error
	UnblockKey(key string, keyType string) error
	Increment(key string, keyType string) error
	ClearOldAccesses(key string, keyType string, between time.Duration) error
}

type StorageConfig struct {
	Key      string
	KeyType  string
	Accesses []time.Time
}
