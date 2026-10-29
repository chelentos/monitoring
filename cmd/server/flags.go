package main

import (
	"flag"

	"github.com/chelentos/monitoring/internal/validation"
)

var flagRunAddr string

func parseServerFlags() error {
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")

	flag.Parse()

	err := validation.CheckAddr(flagRunAddr)

	return err
}
