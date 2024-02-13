package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/felipemagrassi/gandalf/ratelimiter"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
)

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! %s %s", time.Now(), r.Header.Get("API_KEY"))
}

func main() {
	godotenv.Load(".env")
	ratelimiter := ratelimiter.NewRateLimiterMiddleware()
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(ratelimiter)
	r.Get("/", greet)

	fmt.Println("Server started at :8081")
	err := http.ListenAndServe(":8081", r)
	if err != nil {
		panic(err)
	}
}
