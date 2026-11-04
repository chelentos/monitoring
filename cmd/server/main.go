package main

import (
	"log"
	"net/http"
	"os"

	"github.com/chelentos/monitoring/internal/config"
	"github.com/chelentos/monitoring/internal/handler"
	"github.com/chelentos/monitoring/internal/repository"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg, parseErr := config.ParseServer(os.Args[1:])
	if parseErr != nil {
		log.Println(parseErr)
		return
	}

	r := chi.NewRouter()

	storage := repository.InitMemStorage()

	r.Post("/update/{type}/{name}/{value}", handler.UpdateMetricHandler(storage))
	r.Get("/value/{type}/{name}", handler.GetMetricHandler(storage))
	r.Get("/", handler.GetAllMetricHandler(storage))

	log.Println("Running server on", cfg.Address)
	err := http.ListenAndServe(cfg.Address, r)

	if err != nil {
		log.Fatal(err)
	}
}
