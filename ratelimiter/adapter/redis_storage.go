package adapter

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var ConnectionError = errors.New("connection error")

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(host, port, password string, database int) (Storage, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       database,
	})
	_, err := client.Ping(context.Background()).Result()

	if err != nil {
		return nil, ConnectionError
	}

	return &RedisStorage{
		client: client,
	}, nil
}

func (rs *RedisStorage) GetKeyInfo(key string, keyType string) (*StorageConfig, error) {
	redisKey := rs.redisKey("access", key, keyType)
	cn := rs.client.Conn()
	defer cn.Close()

	accesses, err := cn.ZCard(context.Background(), redisKey).Result()
	if err != nil {
		return nil, err
	}

	return &StorageConfig{
		KeyType:  keyType,
		Key:      key,
		Accesses: int(accesses),
	}, nil
}

func (rs *RedisStorage) ClearOldAccesses(key string, keyType string, between time.Duration) error {
	cn := rs.client.Conn()
	defer cn.Close()

	redisKey := rs.redisKey("access", key, keyType)
	until := time.Now().Add(-between).UnixNano()

	_, err := cn.ZRemRangeByScore(context.Background(), redisKey, "-inf", fmt.Sprintf("%d", until)).Result()
	if err != nil {
		return err
	}

	return nil
}

func (rs *RedisStorage) Increment(key string, keyType string) error {
	cn := rs.client.Conn()
	defer cn.Close()

	redisKey := rs.redisKey("access", key, keyType)
	now := time.Now()

	_, err := cn.ZAdd(context.Background(), redisKey, redis.Z{Score: float64(now.UnixNano()), Member: now.Format(time.RFC3339Nano)}).Result()
	if err != nil {
		return err
	}

	return nil
}

func (rs *RedisStorage) GetBlockedKey(key string, keyType string) (*time.Time, error) {
	cn := rs.client.Conn()
	defer cn.Close()

	redisKey := rs.redisKey("block", key, keyType)

	blockedAt, err := cn.Get(context.Background(), redisKey).Result()
	if err == redis.Nil {
		return nil, KeyNotFound
	}

	if err != nil {
		return nil, err
	}

	parsedBlockedAt, err := time.Parse(time.RFC3339Nano, blockedAt)
	if err != nil {
		return nil, err
	}

	return &parsedBlockedAt, nil
}

func (rs *RedisStorage) BlockKey(key string, keyType string) error {
	cn := rs.client.Conn()
	defer cn.Close()

	redisKey := rs.redisKey("block", key, keyType)

	_, err := cn.Set(context.Background(), redisKey, time.Now().Format(time.RFC3339Nano), 0).Result()
	if err != nil {
		return err
	}

	return nil
}

func (rs *RedisStorage) UnblockKey(key string, keyType string) error {
	cn := rs.client.Conn()
	defer cn.Close()

	redisKey := rs.redisKey("block", key, keyType)

	_, err := cn.Del(context.Background(), redisKey).Result()

	if err != nil {
		return err
	}

	return nil
}

func (rs *RedisStorage) redisKey(redisSet string, key string, keyType string) string {
	return fmt.Sprintf("%s:%s:%s", redisSet, keyType, key)
}
