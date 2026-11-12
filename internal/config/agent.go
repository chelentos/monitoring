package config

import (
	"flag"
	"fmt"
	"time"

	"github.com/chelentos/monitoring/internal/validation"
)

func ParseAgent(args []string) (AgentConfig, error) {
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

	addrErr := validation.CheckAddr(cfg.Sender.Address)
	if addrErr != nil {
		return AgentConfig{}, addrErr
	}

	if pollSeconds <= 0 || reportSeconds <= 0 {
		return AgentConfig{}, fmt.Errorf(
			"intervals -p и -r must be >0",
		)
	}

	cfg.Collector.PollInterval = time.Duration(pollSeconds) * time.Second
	cfg.Sender.ReportInterval = time.Duration(reportSeconds) * time.Second

	return cfg, nil
}
