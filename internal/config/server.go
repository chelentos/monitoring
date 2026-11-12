package config

import (
	"flag"

	"github.com/chelentos/monitoring/internal/validation"
)

func ParseServer(args []string) (ServerConfig, error) {
	var cfg ServerConfig

	fs := flag.NewFlagSet("server", flag.ContinueOnError)
	fs.StringVar(&cfg.Address, "a", "localhost:8080", "address and port to run server")

	parseErr := fs.Parse(args)
	if parseErr != nil {
		return ServerConfig{}, parseErr
	}

	addrErr := validation.CheckAddr(cfg.Address)
	if addrErr != nil {
		return ServerConfig{}, addrErr
	}

	return cfg, nil
}
