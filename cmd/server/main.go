package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! %s", time.Now())
}

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Get("/", greet)

	fmt.Println("Server started at :8081")
	err := http.ListenAndServe(":8081", r)
	if err != nil {
		panic(err)
	}
}
