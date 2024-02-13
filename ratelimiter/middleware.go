package ratelimiter

import (
	"net"
	"net/http"

	"github.com/felipemagrassi/gandalf/ratelimiter/adapter"
)

type RateLimiterConfig struct {
	tokenTimeout float64
	tokenRps     int
	ipTimeout    float64
	ipRps        int
	rateLimiter  *RateLimiter
}

func NewRateLimiterMiddleware() func(http.Handler) http.Handler {
	storage := newStorage()
	config := newRateLimitConfig(storage)

	return func(next http.Handler) http.Handler {
		return rateLimiterMiddleware(config, next)
	}
}

func rateLimiterMiddleware(config *RateLimiterConfig, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("API_KEY")
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)

		if token != "" {
			err := config.rateLimiter.Increment(token, "token", config.tokenRps, config.tokenTimeout)
			if err != nil {
				RateLimitError(w, r)
				next.ServeHTTP(w, r)
			}
		} else {
			err := config.rateLimiter.Increment(ip, "ip", config.ipRps, config.ipTimeout)
			if err != nil {
				RateLimitError(w, r)
				next.ServeHTTP(w, r)
			}
		}

		next.ServeHTTP(w, r)
	})
}

func newRateLimitConfig(storage adapter.Storage) *RateLimiterConfig {
	return &RateLimiterConfig{
		tokenTimeout: 1.5,
		tokenRps:     5,
		ipTimeout:    1.5,
		ipRps:        5,
		rateLimiter:  NewRateLimiter(storage),
	}
}

func newStorage() adapter.Storage {
	return adapter.NewMemoryStorage()
}

func RateLimitError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusTooManyRequests)
	w.Write([]byte("you have reached the maximum number of requests or actions allowed within a certain time frame"))
}
