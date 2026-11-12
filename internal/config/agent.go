package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/chelentos/monitoring/internal/validation"
)

func validateAgentConfig(cfg AgentConfig) error {
	addrErr := validation.CheckAddr(cfg.Sender.Address)
	if addrErr != nil {
		return addrErr
	}

	if !(cfg.Collector.PollInterval > 0 && cfg.Sender.ReportInterval > 0) {
		return fmt.Errorf(
			"poll(-p) and report(-r) intervals mast be >0",
		)
	}

	return nil
}

func parseAgentCli(args []string) (AgentConfig, error) {
	var cfg AgentConfig
	var pollSeconds, reportSeconds int64

	fs := flag.NewFlagSet("agent", flag.ContinueOnError)
	fs.StringVar(&cfg.Sender.Address, "a", "localhost:8080", "address and port of server")
	fs.Int64Var(&reportSeconds, "r", 10, "report interval")
	fs.Int64Var(&pollSeconds, "p", 2, "polling interval")

	parseErr := fs.Parse(args)
	if parseErr != nil {
		return AgentConfig{}, parseErr
	}

	cfg.Collector.PollInterval = time.Duration(pollSeconds) * time.Second
	cfg.Sender.ReportInterval = time.Duration(reportSeconds) * time.Second

	return cfg, nil
}

func parseAgentEnv(cfg AgentConfig, environ map[string]string) (AgentConfig, error) {
	if environ == nil {
		environ = map[string]string{}
	}

	err := env.ParseWithOptions(&cfg, env.Options{
		Environment: environ,
	})
	if err != nil {
		return AgentConfig{}, err
	}

	validationErr := validateAgentConfig(cfg)
	if validationErr != nil {
		return AgentConfig{}, validationErr
	}

	return cfg, nil
}

func ParseAgentConfig(args []string) (AgentConfig, error) {
	cliCfg, errCli := parseAgentCli(args)
	if errCli != nil {
		return AgentConfig{}, errCli
	}

	envCfg, errEnv := parseAgentEnv(cliCfg, env.ToMap(os.Environ()))
	if errEnv != nil {
		return AgentConfig{}, errEnv
	}

	return envCfg, nil
}
