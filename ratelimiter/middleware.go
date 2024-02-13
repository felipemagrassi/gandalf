package ratelimiter

import (
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/felipemagrassi/gandalf/ratelimiter/adapter"
)

const (
	EnvTokenTimeout  = "TOKEN_TIMEOUT"
	EnvTokenRps      = "TOKEN_RPS"
	EnvIpTimeout     = "IP_TIMEOUT"
	EnvIpRps         = "IP_RPS"
	EnvRedisHost     = "REDIS_HOST"
	EnvRedisPort     = "REDIS_PORT"
	EnvRedisDatabase = "REDIS_DATABASE"
	EnvRedisPassword = "REDIS_PASSWORD"
)

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	Database int
}

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
	tokenTimeout := envFloat64Lookup(EnvTokenTimeout, 5.0)
	tokenRps := envIntLookup(EnvTokenRps, 10)
	ipTimeout := envFloat64Lookup(EnvIpTimeout, 5.0)
	ipRps := envIntLookup(EnvIpRps, 10)

	return &RateLimiterConfig{
		tokenTimeout: tokenTimeout,
		tokenRps:     tokenRps,
		ipTimeout:    ipTimeout,
		ipRps:        ipRps,
		rateLimiter:  NewRateLimiter(storage),
	}
}

func newStorage() adapter.Storage {
	redisHost := envLookup(EnvRedisHost, "")
	redisPort := envLookup(EnvRedisPort, "")
	redisPassword := envLookup(EnvRedisPassword, "")
	redisDatabase := envIntLookup(EnvRedisDatabase, 0)

	redis := RedisConfig{
		Host:     redisHost,
		Port:     redisPort,
		Password: redisPassword,
		Database: redisDatabase,
	}

	if redis.Host != "" {
		return adapter.NewMemoryStorage()
	}

	return adapter.NewMemoryStorage()
}

func RateLimitError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusTooManyRequests)
	w.Write([]byte("you have reached the maximum number of requests or actions allowed within a certain time frame"))
}
func envLookup(key string, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	return value
}

func envFloat64Lookup(key string, defaultValue float64) float64 {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	parsedValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return defaultValue
	}

	return parsedValue
}

func envIntLookup(key string, defaultValue int) int {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return parsedValue
}
