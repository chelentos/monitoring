package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/chelentos/monitoring/internal/validation"
)

type agentConfigValues struct {
	Address       string `env:"ADDRESS"`
	PollSeconds   int64  `env:"POLL_INTERVAL"`
	ReportSeconds int64  `env:"REPORT_INTERVAL"`
}

func validateAgentConfig(cfg agentConfigValues) error {
	addrErr := validation.CheckAddr(cfg.Address)
	if addrErr != nil {
		return addrErr
	}

	if !(cfg.PollSeconds > 0 && cfg.ReportSeconds > 0) {
		return fmt.Errorf(
			"poll(-p) and report(-r) intervals mast be >0",
		)
	}

	return nil
}

func parseAgentCli(args []string) (agentConfigValues, error) {
	var cfg agentConfigValues

	fs := flag.NewFlagSet("agent", flag.ContinueOnError)
	fs.StringVar(&cfg.Address, "a", "localhost:8080", "address and port of server")
	fs.Int64Var(&cfg.ReportSeconds, "r", 10, "report interval")
	fs.Int64Var(&cfg.PollSeconds, "p", 2, "polling interval")

	parseErr := fs.Parse(args)
	if parseErr != nil {
		return agentConfigValues{}, parseErr
	}

	return cfg, nil
}

func parseAgentEnv(cfg agentConfigValues, environ map[string]string) (agentConfigValues, error) {
	if environ == nil {
		environ = map[string]string{}
	}

	err := env.ParseWithOptions(&cfg, env.Options{
		Environment: environ,
	})
	if err != nil {
		return agentConfigValues{}, err
	}

	validationErr := validateAgentConfig(cfg)
	if validationErr != nil {
		return agentConfigValues{}, validationErr
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

	err := validateAgentConfig(envCfg)
	if err != nil {
		return AgentConfig{}, err
	}

	return AgentConfig{
		Collector: CollectorConfig{
			PollInterval: time.Duration(envCfg.PollSeconds) * time.Second,
		},
		Sender: SenderConfig{
			Address:        envCfg.Address,
			ReportInterval: time.Duration(envCfg.ReportSeconds) * time.Second,
		},
	}, nil
}
