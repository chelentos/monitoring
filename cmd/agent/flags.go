package main

import (
	"flag"

	"github.com/chelentos/monitoring/internal/validation"
)

var flagRunAddr string
var flagReportInterval int
var flagPollInterval int

func parseAgentFlags() error {
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&flagReportInterval, "r", 10, "report interval")
	flag.IntVar(&flagPollInterval, "p", 2, "polling interval")

	flag.Parse()

	err := validation.CheckAddr(flagRunAddr)

	if err != nil {
		return err
	}

	return err
}
