package main

import (
	"fmt"
	"net/http"

	"github.com/chelentos/monitoring/internal/handler"
	"github.com/chelentos/monitoring/internal/repository"
	"github.com/go-chi/chi/v5"
)

func main() {
	flagsErr := parseServerFlags()

	if flagsErr != nil {
		fmt.Println(flagsErr)
		return
	}

	r := chi.NewRouter()

	storage := repository.InitMemStorage()

	r.Post("/update/{type}/{name}/{value}", handler.UpdateMetricHandler(storage))
	r.Get("/value/{type}/{name}", handler.GetMetricHandler(storage))
	r.Get("/", handler.GetAllMetricHandler(storage))

	fmt.Println("Running server on", flagRunAddr)

	err := http.ListenAndServe(flagRunAddr, r)

	if err != nil {
		panic(err)
	}
}
