package main

import (
	"fmt"
	"net/http"

	"github.com/felipemagrassi/gandalf/ratelimiter"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
)

func greet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!\n"))
}

func main() {
	godotenv.Load(".env")
	ratelimiter := ratelimiter.NewRateLimiterMiddleware()
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(ratelimiter)
	r.Get("/", greet)

	fmt.Println("Server started at :8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
