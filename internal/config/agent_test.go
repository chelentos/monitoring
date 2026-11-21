package config

import (
	"os"
	"testing"
	"time"
)

func agentConfig(address string, poll, report time.Duration) AgentConfig {
	return AgentConfig{
		Collector: CollectorConfig{PollInterval: poll},
		Sender:    SenderConfig{Address: address, ReportInterval: report},
	}
}

func testAgentValues(address string, poll, report int64) agentConfigValues {
	return agentConfigValues{
		Address:       address,
		PollSeconds:   poll,
		ReportSeconds: report,
	}
}

func TestParseAgentCli(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    agentConfigValues
		wantErr bool
	}{
		{name: "defaults", want: testAgentValues("localhost:8080", 2, 10)},
		{name: "all flags", args: []string{"-a", "localhost:9090", "-p", "3", "-r", "15"}, want: testAgentValues("localhost:9090", 3, 15)},
		{name: "partial flags", args: []string{"-p", "5"}, want: testAgentValues("localhost:8080", 5, 10)},
		{name: "unknown flag", args: []string{"-unknown"}, wantErr: true},
		{name: "missing value", args: []string{"-a"}, wantErr: true},
		{name: "invalid poll", args: []string{"-p", "abc"}, wantErr: true},
		{name: "invalid report", args: []string{"-r", "1.5"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseAgentCli(tt.args)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseAgentCli() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && got != tt.want {
				t.Errorf("parseAgentCli() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestParseAgentEnv(t *testing.T) {
	base := testAgentValues("localhost:9090", 3, 15)
	tests := []struct {
		name    string
		environ map[string]string
		want    agentConfigValues
		wantErr bool
	}{
		{name: "nil environment", want: base},
		{name: "empty environment", environ: map[string]string{}, want: base},
		{name: "unrelated variable", environ: map[string]string{"OTHER": "value"}, want: base},
		{name: "address override", environ: map[string]string{"ADDRESS": "localhost:7070"}, want: testAgentValues("localhost:7070", 3, 15)},
		{name: "poll override", environ: map[string]string{"POLL_INTERVAL": "1"}, want: testAgentValues("localhost:9090", 1, 15)},
		{name: "report override", environ: map[string]string{"REPORT_INTERVAL": "60"}, want: testAgentValues("localhost:9090", 3, 60)},
		{name: "all overrides", environ: map[string]string{"ADDRESS": "localhost:7070", "POLL_INTERVAL": "4", "REPORT_INTERVAL": "20"}, want: testAgentValues("localhost:7070", 4, 20)},
		{name: "invalid address", environ: map[string]string{"ADDRESS": "invalid"}, wantErr: true},
		{name: "invalid poll", environ: map[string]string{"POLL_INTERVAL": "abc"}, wantErr: true},
		{name: "duration with unit", environ: map[string]string{"POLL_INTERVAL": "3s"}, wantErr: true},
		{name: "fractional seconds", environ: map[string]string{"POLL_INTERVAL": "1.5"}, wantErr: true},
		{name: "integer overflow", environ: map[string]string{"REPORT_INTERVAL": "9223372036854775808"}, wantErr: true},
		{name: "invalid report", environ: map[string]string{"REPORT_INTERVAL": "abc"}, wantErr: true},
		{name: "zero poll", environ: map[string]string{"POLL_INTERVAL": "0"}, wantErr: true},
		{name: "negative report", environ: map[string]string{"REPORT_INTERVAL": "-1"}, wantErr: true},
	}
	// Explicit input must take precedence over the process environment, even for nil.
	t.Setenv("ADDRESS", "invalid")
	t.Setenv("POLL_INTERVAL", "invalid")
	t.Setenv("REPORT_INTERVAL", "invalid")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := base
			got, err := parseAgentEnv(input, tt.environ)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseAgentEnv() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && got != tt.want {
				t.Errorf("parseAgentEnv() = %+v, want %+v", got, tt.want)
			}
			if input != base {
				t.Errorf("input changed: %+v, want %+v", input, base)
			}
		})
	}
}

func TestValidateAgentConfig(t *testing.T) {
	tests := []struct {
		name    string
		cfg     agentConfigValues
		wantErr bool
	}{
		{name: "valid", cfg: testAgentValues("localhost:8080", 1, 1)},
		{name: "empty address", cfg: testAgentValues("", 1, 1), wantErr: true},
		{name: "missing port", cfg: testAgentValues("localhost", 1, 1), wantErr: true},
		{name: "invalid port", cfg: testAgentValues("localhost:abc", 1, 1), wantErr: true},
		{name: "zero port", cfg: testAgentValues("localhost:0", 1, 1), wantErr: true},
		{name: "port out of range", cfg: testAgentValues("localhost:65536", 1, 1), wantErr: true},
		{name: "zero poll", cfg: testAgentValues("localhost:8080", 0, 1), wantErr: true},
		{name: "negative poll", cfg: testAgentValues("localhost:8080", -1, 1), wantErr: true},
		{name: "zero report", cfg: testAgentValues("localhost:8080", 1, 0), wantErr: true},
		{name: "negative report", cfg: testAgentValues("localhost:8080", 1, -1), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateAgentConfig(tt.cfg); (err != nil) != tt.wantErr {
				t.Errorf("validateAgentConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseAgentConfig(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		environ map[string]string
		want    AgentConfig
		wantErr bool
	}{
		{name: "defaults", want: agentConfig("localhost:8080", 2*time.Second, 10*time.Second)},
		{name: "flags", args: []string{"-a", "localhost:9090", "-p", "3", "-r", "15"}, want: agentConfig("localhost:9090", 3*time.Second, 15*time.Second)},
		{name: "integer seconds from environment", environ: map[string]string{"POLL_INTERVAL": "1", "REPORT_INTERVAL": "5"}, want: agentConfig("localhost:8080", time.Second, 5*time.Second)},
		{name: "environment", environ: map[string]string{"ADDRESS": "localhost:7070", "POLL_INTERVAL": "4", "REPORT_INTERVAL": "20"}, want: agentConfig("localhost:7070", 4*time.Second, 20*time.Second)},
		{name: "environment overrides flags", args: []string{"-a", "localhost:9090", "-p", "3", "-r", "15"}, environ: map[string]string{"ADDRESS": "localhost:7070", "POLL_INTERVAL": "4", "REPORT_INTERVAL": "20"}, want: agentConfig("localhost:7070", 4*time.Second, 20*time.Second)},
		{name: "all three sources", args: []string{"-p", "3"}, environ: map[string]string{"ADDRESS": "localhost:7070"}, want: agentConfig("localhost:7070", 3*time.Second, 10*time.Second)},
		{name: "environment corrects invalid flag values", args: []string{"-a", "invalid", "-p", "0", "-r", "-1"}, environ: map[string]string{"ADDRESS": "localhost:7070", "POLL_INTERVAL": "4", "REPORT_INTERVAL": "20"}, want: agentConfig("localhost:7070", 4*time.Second, 20*time.Second)},
		{name: "invalid env does not fall back to flag", args: []string{"-p", "3"}, environ: map[string]string{"POLL_INTERVAL": "abc"}, wantErr: true},
		{name: "invalid address env", environ: map[string]string{"ADDRESS": "invalid"}, wantErr: true},
		{name: "invalid flag syntax", args: []string{"-p", "abc"}, environ: map[string]string{"POLL_INTERVAL": "4"}, wantErr: true},
		{name: "unknown flag", args: []string{"-unknown"}, wantErr: true},
		{name: "invalid final address", args: []string{"-a", "invalid"}, wantErr: true},
		{name: "invalid final interval", args: []string{"-r", "0"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, key := range []string{"ADDRESS", "POLL_INTERVAL", "REPORT_INTERVAL"} {
				// Register restoration with testing before removing inherited values.
				t.Setenv(key, "")
				if err := os.Unsetenv(key); err != nil {
					t.Fatal(err)
				}
			}
			for key, value := range tt.environ {
				t.Setenv(key, value)
			}
			got, err := ParseAgentConfig(tt.args)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseAgentConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && got != tt.want {
				t.Errorf("ParseAgentConfig() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
