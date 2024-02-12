package ratelimiter

import (
	"errors"
	"time"

	"github.com/felipemagrassi/gandalf/ratelimiter/adapter"
)

type keyConfig struct {
	rps     int
	timeout time.Duration
}

type RateLimiter struct {
	StorageAdapter adapter.Storage
}

func NewRateLimiter(storageAdapter adapter.Storage) *RateLimiter {
	return &RateLimiter{
		StorageAdapter: storageAdapter,
	}
}

func (rl *RateLimiter) AddKey(key string, keyType string, rps, timeout int) error {
	if key == "" {
		return errors.New("Key is required")
	}

	if keyType == "" {
		return errors.New("KeyType is required")
	}

	if rps <= 0 || timeout <= 0 {
		return errors.New("RPS and Timeout must be greater than 0")
	}

	err := rl.StorageAdapter.AddKey(key, keyType, rps, timeout)
	return err
}

func (rl *RateLimiter) Increment(key string, keyType string) (*time.Duration, error) {
	if key == "" {
		return nil, errors.New("Key is required")
	}

	if keyType == "" {
		return nil, errors.New("KeyType is required")
	}

	timeLeft, err := rl.StorageAdapter.Increment(key, keyType)
	return timeLeft, err
}
