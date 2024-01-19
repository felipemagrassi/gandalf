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

	r.Get("/", greet)

	r.Use(middleware.Recoverer)

	err := http.ListenAndServe(":8081", r)
	if err != nil {
		panic(err)
	}
}
