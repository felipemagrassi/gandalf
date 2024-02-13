package ratelimiter

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
)

func TestCanCreateRateLimiterMiddleware(t *testing.T) {
	t.Setenv(EnvTokenTimeout, "1")
	t.Setenv(EnvTokenRps, "10")
	t.Setenv(EnvUseRedis, "0")

	r := chi.NewRouter()
	path := "/"

	r.Use(NewRateLimiterMiddleware())

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	req, err := http.NewRequest("GET", ts.URL+path, nil)
	req.Header.Set("API_KEY", "123")
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusTooManyRequests {
		t.Errorf("Expected status code %d, got %d", http.StatusTooManyRequests, resp.StatusCode)
	}
}

func TestCanCreateRedisRateLimiterMiddleware(t *testing.T) {
	t.Setenv(EnvTokenTimeout, "1")
	t.Setenv(EnvTokenRps, "10")
	t.Setenv(EnvRedisHost, "localhost")
	t.Setenv(EnvRedisPort, "6379")
	t.Setenv(EnvRedisPassword, "")
	t.Setenv(EnvUseRedis, "1")

	r := chi.NewRouter()
	path := "/"

	r.Use(NewRateLimiterMiddleware())

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	req, err := http.NewRequest("GET", ts.URL+path, nil)
	req.Header.Set("API_KEY", "123")
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusTooManyRequests {
		t.Errorf("Expected status code %d, got %d", http.StatusTooManyRequests, resp.StatusCode)
	}

}
