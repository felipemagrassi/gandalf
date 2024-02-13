package ratelimiter

import (
	"testing"
	"time"

	"github.com/felipemagrassi/gandalf/ratelimiter/adapter"
)

func TestRateLimiterCanIncrease(t *testing.T) {
	rateLimiter := NewRateLimiter(
		adapter.NewMemoryStorage(),
	)
	key := "1"
	token := "token"
	rps := 10
	timeout := 1.5

	for i := 0; i < 10; i++ {
		err := rateLimiter.Increment(key, token, rps, timeout)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}
	}

	err := rateLimiter.Increment(key, token, rps, timeout)
	if err != ReachedMaxTries {
		t.Errorf("Expected %v, got %v", ReachedMaxTries, err)
	}

	time.Sleep(2 * time.Second)

	for i := 0; i < 10; i++ {
		err := rateLimiter.Increment(key, token, rps, timeout)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}
	}

	time.Sleep(1 * time.Second)

	for i := 0; i < 10; i++ {
		err := rateLimiter.Increment(key, token, rps, timeout)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}
	}

	rateLimiter.Increment(key, token, rps, timeout)
	err = rateLimiter.Increment(key, token, rps, timeout)
	if err != BlockedKey {
		t.Errorf("Expected %v, got %v", BlockedKey, err)
	}
}

func TestRateLimiterCannotIncrease(t *testing.T) {
	rateLimiter := NewRateLimiter(
		adapter.NewMemoryStorage(),
	)
	err := rateLimiter.Increment("1", "", 10, 1)
	if err != KeyTypeNotFound {
		t.Errorf("Expected %v, got %v", KeyTypeNotFound, err)
	}

	err = rateLimiter.Increment("", "token", 10, 1)
	if err != KeyNotFound {
		t.Errorf("Expected %v, got %v", KeyNotFound, err)
	}

	err = rateLimiter.Increment("1", "token", 10, 0)
	if err != InvalidTimeout {
		t.Errorf("Expected %v, got %v", InvalidTimeout, err)
	}

	err = rateLimiter.Increment("1", "token", 0, 1)
	if err != InvalidRps {
		t.Errorf("Expected %v, got %v", InvalidRps, err)
	}
}
