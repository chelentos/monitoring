package config

import (
	"time"
)

type ServerConfig struct {
	Address string `env:"ADDRESS"`
}

type CollectorConfig struct {
	PollInterval time.Duration `env:"POLL_INTERVAL"`
}

type SenderConfig struct {
	Address        string        `env:"ADDRESS"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
}

type AgentConfig struct {
	Collector CollectorConfig
	Sender    SenderConfig
}
