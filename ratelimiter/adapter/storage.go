package adapter

import "time"

type Storage interface {
	AddKey(key string, keyType string, rps, timeout int) error
	Increment(key string, keyType string) (*time.Duration, error)
	DeleteKey(key string, keyType string) error
}

type StorageConfig struct {
	rps     int
	timeout time.Duration
}
