package ratelimiter

import (
	"testing"

	"github.com/felipemagrassi/gandalf/ratelimiter/adapter"
)

func TestRateLimiterCanIncrease(t *testing.T) {
	rateLimiter := NewRateLimiter(adapter.NewMemoryStorage())

	err := rateLimiter.AddKey("1", "token", 10, 1)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	_, err = rateLimiter.Increment("1", "token")
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
}

func TestRateLimiterCanChangeConfig(t *testing.T) {
	rateLimiter := NewRateLimiter(adapter.NewMemoryStorage())

	err := rateLimiter.AddKey("1", "token", 10, 1)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	err = rateLimiter.AddKey("1", "token", 20, 2)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
}

func TestRateLimiterCantIncreaseWithoutAdd(t *testing.T) {
	rateLimiter := NewRateLimiter(adapter.NewMemoryStorage())

	_, err := rateLimiter.Increment("1", "token")
	if err == nil {
		t.Errorf("Expected Key not added, got %v", err)
	}
}

func TestRateLimiterCantAddInvalidRps(t *testing.T) {
	rateLimiter := NewRateLimiter(adapter.NewMemoryStorage())

	err := rateLimiter.AddKey("1", "token", 0, 1)
	if err == nil {
		t.Errorf("Expected RPS and Timeout must be greater than 0, got %v", err)
	}

	err = rateLimiter.AddKey("1", "token", 1, 0)
	if err == nil {
		t.Errorf("Expected RPS and Timeout must be greater than 0, got %v", err)
	}
}

func TestRateLimiterTimeout(t *testing.T) {
	rateLimiter := NewRateLimiter(adapter.NewMemoryStorage())
	err := rateLimiter.AddKey("1", "token", 10, 1)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	for i := 0; i < 10; i++ {
		_, err = rateLimiter.Increment("1", "token")
		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}
	}

	timeLeft, err := rateLimiter.Increment("1", "token")
	if err == nil {
		t.Errorf("Expected Key is blocked, got %v", err)
	}

	if timeLeft == nil {
		t.Errorf("Expected timeLeft, got %v", timeLeft)
	}

}
