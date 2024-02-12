package adapter

import (
	"errors"
	"sync"
	"time"
)

var (
	KeyNotFound          = errors.New("Key not found")
	KeyTypeNotFound      = errors.New("KeyTypeNotFound")
	ReachedMaxTries      = errors.New("Reached max tries")
	KeyAlreadyRegistered = errors.New("Key already registered")
	KeyIsBlocked         = errors.New("Key is blocked")
)

type MemoryStorage struct {
	mutex       sync.Mutex
	accessMutex sync.Mutex
	blockMutex  sync.Mutex
	keyConfig   map[string]map[string]StorageConfig
	accesses    map[string]map[string][]time.Time
	blocked     map[string]map[string]time.Time
}

func NewMemoryStorage() Storage {
	return &MemoryStorage{
		keyConfig:   make(map[string]map[string]StorageConfig),
		blocked:     make(map[string]map[string]time.Time),
		accesses:    make(map[string]map[string][]time.Time),
		accessMutex: sync.Mutex{},
		blockMutex:  sync.Mutex{},
	}
}

func (ms *MemoryStorage) GetKeyInfo(key string, keyType string) (*StorageConfig, error) {
	ms.accessMutex.Lock()
	defer ms.accessMutex.Unlock()

	_, ok := ms.keyConfig[keyType]
	if !ok {
		return nil, KeyTypeNotFound
	}

	_, ok = ms.keyConfig[keyType][key]
	if !ok {
		return nil, KeyNotFound
	}

	return &StorageConfig{
		keyType:  keyType,
		key:      key,
		accesses: ms.accesses[keyType][key],
	}, nil
}

func (ms *MemoryStorage) GetBlockedKey(key string, keyType string) (*time.Time, error) {
	ms.blockMutex.Lock()
	defer ms.blockMutex.Unlock()

	_, ok := ms.blocked[keyType]
	if !ok {
		return nil, KeyTypeNotFound
	}
	_, ok = ms.blocked[keyType][key]
	if !ok {
		return nil, KeyNotFound
	}

	blockedAt := ms.blocked[keyType][key]
	return &blockedAt, nil
}

func (ms *MemoryStorage) Increment(key string, keyType string) error {
	ms.accessMutex.Lock()
	defer ms.accessMutex.Unlock()

	_, ok := ms.accesses[keyType]
	if !ok {
		return KeyTypeNotFound
	}
	_, ok = ms.accesses[keyType][key]
	if !ok {
		return KeyNotFound
	}

	ms.accesses[keyType][key] = append(ms.accesses[keyType][key], time.Now())

	return nil
}

func (ms *MemoryStorage) BlockKey(key string, keyType string) error {
	ms.blockMutex.Lock()
	defer ms.blockMutex.Unlock()

	_, ok := ms.keyConfig[keyType]
	if !ok {
		return KeyTypeNotFound
	}

	_, ok = ms.keyConfig[keyType][key]
	if !ok {
		return KeyNotFound
	}

	_, ok = ms.blocked[keyType]
	if !ok {
		ms.blocked[keyType] = make(map[string]time.Time)
	}

	ms.blocked[keyType][key] = time.Now()

	return nil
}

func (ms *MemoryStorage) UnblockKey(key string, keyType string) error {
	ms.blockMutex.Lock()
	defer ms.blockMutex.Unlock()
	_, ok := ms.blocked[keyType]
	if !ok {
		return KeyTypeNotFound
	}
	_, ok = ms.blocked[keyType][key]
	if !ok {
		return KeyNotFound
	}
	delete(ms.blocked[keyType], key)
	return nil
}

func (ms *MemoryStorage) AddKey(key string, keyType string) error {
	ms.accessMutex.Lock()
	defer ms.accessMutex.Unlock()

	_, ok := ms.keyConfig[keyType]

	if !ok {
		ms.keyConfig[keyType] = make(map[string]StorageConfig)
	}

	_, ok = ms.accesses[keyType]
	if !ok {
		ms.accesses[keyType] = make(map[string][]time.Time)
	}

	_, ok = ms.accesses[keyType][key]
	if !ok {
		ms.accesses[keyType][key] = make([]time.Time, 0)
	}

	_, ok = ms.keyConfig[keyType][key]

	if !ok {
		ms.keyConfig[keyType][key] = StorageConfig{
			key:      key,
			keyType:  keyType,
			accesses: ms.accesses[keyType][key],
		}
		return nil
	}

	return KeyAlreadyRegistered
}
