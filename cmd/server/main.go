package main

import (
	"net/http"

	"github.com/chelentos/monitoring/internal/handler"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	r.Post("/update/{type}/{name}/{value}", handler.UpdateHandler)

	err := http.ListenAndServe(`:8080`, r)

	if err != nil {
		panic(err)
	}
}
