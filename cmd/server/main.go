package main

import (
	"net/http"

	"github.com/chelentos/monitoring/internal/handler"
	"github.com/chelentos/monitoring/internal/repository"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	storage := repository.InitMemStorage()

	r.Post("/update/{type}/{name}/{value}", handler.UpdateMetricHandler(storage))
	r.Get("/value/{type}/{name}", handler.GetMetricHandler(storage))
	r.Get("/", handler.GetAllMetricHandler(storage))

	err := http.ListenAndServe(`:8080`, r)

	if err != nil {
		panic(err)
	}
}
