package config

import (
	"flag"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/chelentos/monitoring/internal/validation"
)

type serverConfigValues struct {
	Address string `env:"ADDRESS"`
}

func validateServerConfig(cfg serverConfigValues) error {
	addrErr := validation.CheckAddr(cfg.Address)
	if addrErr != nil {
		return addrErr
	}

	return nil
}

func parseServerCli(args []string) (serverConfigValues, error) {
	var cfg serverConfigValues

	fs := flag.NewFlagSet("server", flag.ContinueOnError)
	fs.StringVar(&cfg.Address, "a", "localhost:8080", "address and port to run server")

	parseErr := fs.Parse(args)
	if parseErr != nil {
		return serverConfigValues{}, parseErr
	}

	return cfg, nil
}

func parseServerEnv(cfg serverConfigValues, environ map[string]string) (serverConfigValues, error) {
	if environ == nil {
		environ = map[string]string{}
	}

	err := env.ParseWithOptions(&cfg, env.Options{
		Environment: environ,
	})
	if err != nil {
		return serverConfigValues{}, err
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

	err := validateServerConfig(envCfg)
	if err != nil {
		return ServerConfig{}, err
	}

	return ServerConfig(envCfg), nil
}
