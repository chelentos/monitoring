package config

import "time"

type ServerConfig struct {
	Address string
}

type CollectorConfig struct {
	PollInterval time.Duration
}

type SenderConfig struct {
	Address        string
	ReportInterval time.Duration
}

type AgentConfig struct {
	Collector CollectorConfig
	Sender    SenderConfig
}
