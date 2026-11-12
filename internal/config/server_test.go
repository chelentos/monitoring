package config

import (
	"os"
	"testing"
)

func TestParseServerCli(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    ServerConfig
		wantErr bool
	}{
		{name: "defaults", want: ServerConfig{Address: "localhost:8080"}},
		{name: "address flag", args: []string{"-a", "localhost:9090"}, want: ServerConfig{Address: "localhost:9090"}},
		{name: "unknown flag", args: []string{"-unknown"}, wantErr: true},
		{name: "missing value", args: []string{"-a"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseServerCli(tt.args)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseServerCli() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && got != tt.want {
				t.Errorf("parseServerCli() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestParseServerEnv(t *testing.T) {
	base := ServerConfig{Address: "localhost:9090"}
	tests := []struct {
		name    string
		environ map[string]string
		want    ServerConfig
		wantErr bool
	}{
		{name: "nil environment", want: base},
		{name: "empty environment", environ: map[string]string{}, want: base},
		{name: "unrelated variable", environ: map[string]string{"OTHER": "value"}, want: base},
		{name: "address override", environ: map[string]string{"ADDRESS": "localhost:7070"}, want: ServerConfig{Address: "localhost:7070"}},
		{name: "invalid address", environ: map[string]string{"ADDRESS": "invalid"}, wantErr: true},
		{name: "invalid port", environ: map[string]string{"ADDRESS": "localhost:abc"}, wantErr: true},
		{name: "port out of range", environ: map[string]string{"ADDRESS": "localhost:65536"}, wantErr: true},
	}
	// Explicit input must be independent of the process environment, even for nil.
	t.Setenv("ADDRESS", "invalid")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := base
			got, err := parseServerEnv(input, tt.environ)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseServerEnv() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && got != tt.want {
				t.Errorf("parseServerEnv() = %+v, want %+v", got, tt.want)
			}
			if input != base {
				t.Errorf("input changed: %+v, want %+v", input, base)
			}
		})
	}
}

func TestValidateServerConfig(t *testing.T) {
	tests := []struct {
		name    string
		address string
		wantErr bool
	}{
		{name: "valid", address: "localhost:8080"},
		{name: "IPv4", address: "127.0.0.1:8080"},
		{name: "IPv6", address: "[::1]:8080"},
		{name: "minimum port", address: "localhost:1"},
		{name: "maximum port", address: "localhost:65535"},
		{name: "empty address", wantErr: true},
		{name: "missing port", address: "localhost", wantErr: true},
		{name: "empty port", address: "localhost:", wantErr: true},
		{name: "invalid port", address: "localhost:abc", wantErr: true},
		{name: "zero port", address: "localhost:0", wantErr: true},
		{name: "negative port", address: "localhost:-1", wantErr: true},
		{name: "port out of range", address: "localhost:65536", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateServerConfig(ServerConfig{Address: tt.address}); (err != nil) != tt.wantErr {
				t.Errorf("validateServerConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseServerConfig(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		environ map[string]string
		want    ServerConfig
		wantErr bool
	}{
		{name: "defaults", want: ServerConfig{Address: "localhost:8080"}},
		{name: "flag", args: []string{"-a", "localhost:9090"}, want: ServerConfig{Address: "localhost:9090"}},
		{name: "environment", environ: map[string]string{"ADDRESS": "localhost:7070"}, want: ServerConfig{Address: "localhost:7070"}},
		{name: "environment overrides flag", args: []string{"-a", "localhost:9090"}, environ: map[string]string{"ADDRESS": "localhost:7070"}, want: ServerConfig{Address: "localhost:7070"}},
		{name: "environment corrects invalid flag value", args: []string{"-a", "invalid"}, environ: map[string]string{"ADDRESS": "localhost:7070"}, want: ServerConfig{Address: "localhost:7070"}},
		{name: "invalid env does not fall back to flag", args: []string{"-a", "localhost:9090"}, environ: map[string]string{"ADDRESS": "invalid"}, wantErr: true},
		{name: "unknown flag", args: []string{"-unknown"}, wantErr: true},
		{name: "missing flag value", args: []string{"-a"}, environ: map[string]string{"ADDRESS": "localhost:7070"}, wantErr: true},
		{name: "invalid final address", args: []string{"-a", "invalid"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Register restoration before removing any inherited ADDRESS.
			t.Setenv("ADDRESS", "")
			if err := os.Unsetenv("ADDRESS"); err != nil {
				t.Fatal(err)
			}
			for key, value := range tt.environ {
				t.Setenv(key, value)
			}
			got, err := ParseServerConfig(tt.args)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseServerConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && got != tt.want {
				t.Errorf("ParseServerConfig() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
