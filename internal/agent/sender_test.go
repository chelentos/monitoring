package agent

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
)

func TestSendGauge(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		wantErr    bool
	}{
		{name: "successful request", statusCode: http.StatusOK},
		{name: "server error", statusCode: http.StatusInternalServerError, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("method = %q, want %q", r.Method, http.MethodPost)
				}
				if r.URL.Path != "/update/gauge/Alloc/12.5" {
					t.Errorf("path = %q, want %q", r.URL.Path, "/update/gauge/Alloc/12.5")
				}
				if got := r.Header.Get("Content-Type"); got != "text/plain" {
					t.Errorf("Content-Type = %q, want text/plain", got)
				}
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			client := resty.NewWithClient(server.Client())
			err := SendGauge(client, server.URL, "Alloc", 12.5)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendGauge() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSendCounter(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		wantErr    bool
	}{
		{name: "successful request", statusCode: http.StatusOK},
		{name: "server error", statusCode: http.StatusInternalServerError, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("method = %q, want %q", r.Method, http.MethodPost)
				}
				if r.URL.Path != "/update/counter/PollCount/3" {
					t.Errorf("path = %q, want %q", r.URL.Path, "/update/counter/PollCount/3")
				}
				if got := r.Header.Get("Content-Type"); got != "text/plain" {
					t.Errorf("Content-Type = %q, want text/plain", got)
				}
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			client := resty.NewWithClient(server.Client())
			err := SendCounter(client, server.URL, "PollCount", 3)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendCounter() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
