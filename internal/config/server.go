package config

import (
	"flag"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/chelentos/monitoring/internal/validation"
)

func validateServerConfig(cfg ServerConfig) error {
	addrErr := validation.CheckAddr(cfg.Address)
	if addrErr != nil {
		return addrErr
	}

	return nil
}

func parseServerCli(args []string) (ServerConfig, error) {
	var cfg ServerConfig

	fs := flag.NewFlagSet("server", flag.ContinueOnError)
	fs.StringVar(&cfg.Address, "a", "localhost:8080", "address and port to run server")

	parseErr := fs.Parse(args)
	if parseErr != nil {
		return ServerConfig{}, parseErr
	}

	return cfg, nil
}

func parseServerEnv(cfg ServerConfig, environ map[string]string) (ServerConfig, error) {
	if environ == nil {
		environ = map[string]string{}
	}

	err := env.ParseWithOptions(&cfg, env.Options{
		Environment: environ,
	})
	if err != nil {
		return ServerConfig{}, err
	}

	validationErr := validateServerConfig(cfg)
	if validationErr != nil {
		return ServerConfig{}, validationErr
	}

	return cfg, nil
}

func ParseServerConfig(args []string) (ServerConfig, error) {
	cliCfg, errCli := parseServerCli(args)
	if errCli != nil {
		return ServerConfig{}, errCli
	}

	envCfg, errEnv := parseServerEnv(cliCfg, env.ToMap(os.Environ()))
	if errEnv != nil {
		return ServerConfig{}, errEnv
	}

	return envCfg, nil
}
