package ratelimiter

import (
	"errors"
	"time"

	"github.com/felipemagrassi/gandalf/ratelimiter/adapter"
)

var (
	KeyNotFound     = errors.New("Key not found")
	KeyTypeNotFound = errors.New("KeyTypeNotFound")
	InvalidRps      = errors.New("RpsNotFound")
	InvalidTimeout  = errors.New("TimeoutNotFound")
	ReachedMaxTries = errors.New("Reached max tries")
	BlockedKey      = errors.New("Blocked key")
)

type RateLimiter struct {
	StorageAdapter adapter.Storage
}

func NewRateLimiter(storageAdapter adapter.Storage) *RateLimiter {
	return &RateLimiter{
		StorageAdapter: storageAdapter,
	}
}

func (rl *RateLimiter) Increment(key string, keyType string, rps int, timeout float64) error {
	if key == "" {
		return KeyNotFound
	}

	if keyType == "" {
		return KeyTypeNotFound
	}

	if rps <= 0 {
		return InvalidRps
	}

	if timeout <= 0 {
		return InvalidTimeout
	}

	blockedAt, _ := rl.StorageAdapter.GetBlockedKey(key, keyType)
	if blockedAt != nil {
		if time.Since(*blockedAt) < time.Duration(timeout)*time.Second {
			return BlockedKey
		}

		rl.StorageAdapter.UnblockKey(key, keyType)
	}

	rl.StorageAdapter.ClearOldAccesses(key, keyType, 1*time.Second)
	keyInfo, _ := rl.StorageAdapter.GetKeyInfo(key, keyType)
	if keyInfo != nil {
		if len(keyInfo.Accesses) >= rps {
			rl.StorageAdapter.BlockKey(key, keyType)
			return ReachedMaxTries
		}
	}

	err := rl.StorageAdapter.Increment(key, keyType)
	if err != nil {
		return err
	}

	return nil
}
