package main

import (
	"net/http"

	"github.com/chelentos/monitoring/internal/handler"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /update/{type}/{name}/{value}", handler.UpdateHandler)

	err := http.ListenAndServe(`:8080`, mux)

	if err != nil {
		panic(err)
	}
}
