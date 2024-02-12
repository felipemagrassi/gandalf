package adapter

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	RateLimitExceeded      = errors.New("Rate limit exceeded")
	KeyWithDifferentConfig = errors.New("Key already exists with different rps")
	KeyTypeNotFound        = errors.New("Key type not found, please use AddKey method")
	KeyNotFound            = errors.New("Key not found, please use AddKey method")
)

type BlockedKey struct {
	duration time.Duration
}

func (b *BlockedKey) Error() string {
	return fmt.Sprintf("Key is blocked for %v, please try again later", b.duration)
}

type MemoryStorage struct {
	mutex     sync.Mutex
	keyConfig map[string]map[string]StorageConfig
	accesses  map[string]map[string][]time.Time
	blocked   map[string]map[string]time.Time
}

func NewMemoryStorage() Storage {
	return &MemoryStorage{
		keyConfig: make(map[string]map[string]StorageConfig),
		blocked:   make(map[string]map[string]time.Time),
		accesses:  make(map[string]map[string][]time.Time),
		mutex:     sync.Mutex{},
	}
}

func (ms *MemoryStorage) Increment(key string, keyType string) (*time.Duration, error) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	config, err := ms.getConfig(key, keyType)
	if err != nil {
		return nil, err
	}

	_, ok := ms.blocked[keyType][key]

	if ok {
		blockedAt := ms.blocked[keyType][key]
		if time.Since(blockedAt) > config.timeout {
			delete(ms.blocked[keyType], key)
		} else {
			timeLeft := config.timeout - time.Since(blockedAt)
			return &timeLeft, &BlockedKey{duration: timeLeft}
		}

	}

	_, ok = ms.accesses[keyType]
	if !ok {
		ms.accesses[keyType] = make(map[string][]time.Time)
	}

	_, ok = ms.accesses[keyType][key]
	if !ok {
		ms.accesses[keyType][key] = make([]time.Time, 0)
	}

	for _, access := range ms.accesses[keyType][key] {
		if time.Since(access) > time.Second {
			ms.accesses[keyType][key] = ms.accesses[keyType][key][1:]
		}
	}

	if len(ms.accesses[keyType][key]) >= config.rps {
		_, ok = ms.blocked[keyType]
		if !ok {
			ms.blocked[keyType] = make(map[string]time.Time)
		}

		ms.blocked[keyType][key] = time.Now()

		return &config.timeout, RateLimitExceeded
	}

	ms.accesses[keyType][key] = append(ms.accesses[keyType][key], time.Now())

	return nil, nil
}

func (ms *MemoryStorage) AddKey(key string, keyType string, rps, timeout int) error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	foundType, ok := ms.keyConfig[keyType]

	if !ok {
		ms.keyConfig[keyType] = make(map[string]StorageConfig)
	}

	foundKey, ok := foundType[key]

	if !ok {
		ms.keyConfig[keyType][key] = StorageConfig{
			rps:     rps,
			timeout: time.Duration(timeout) * time.Second,
		}
		return nil
	}

	if foundKey.rps != rps {
		return KeyWithDifferentConfig
	}

	if foundKey.timeout != time.Duration(timeout)*time.Second {
		return KeyWithDifferentConfig
	}

	return nil
}

func (ms *MemoryStorage) DeleteKey(key string, keyType string) error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	_, ok := ms.keyConfig[keyType][key]
	if !ok {
		return errors.New("Key doesn't exist")
	}

	delete(ms.keyConfig[keyType], key)

	return nil
}

func (ms *MemoryStorage) getConfig(key string, keyType string) (*StorageConfig, error) {
	_, ok := ms.keyConfig[keyType]
	if !ok {
		return nil, KeyTypeNotFound
	}

	_, ok = ms.keyConfig[keyType][key]
	if !ok {
		return nil, KeyNotFound
	}

	return &StorageConfig{rps: ms.keyConfig[keyType][key].rps, timeout: ms.keyConfig[keyType][key].timeout}, nil
}
