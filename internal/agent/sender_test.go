package agent

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-resty/resty/v2"
)

func metricServer(t *testing.T, wantBody map[string]any, statusCode int) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/update" {
			t.Errorf("request = %s %s, want POST /update", r.Method, r.URL.Path)
		}
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", got)
		}

		if got := r.Header.Get("Content-Encoding"); got != "gzip" {
			t.Errorf("Content-Encoding = %q, want gzip", got)
		}
		reader, err := gzip.NewReader(r.Body)
		if err != nil {
			t.Errorf("open gzip request body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer reader.Close()

		body, err := io.ReadAll(reader)
		var got map[string]any
		if err != nil {
			t.Errorf("read request body: %v", err)
		} else if err := json.Unmarshal(body, &got); err != nil {
			t.Errorf("decode request body: %v", err)
		} else if !reflect.DeepEqual(got, wantBody) {
			t.Errorf("request body = %s, want %v", body, wantBody)
		}
		w.WriteHeader(statusCode)
	}))
	t.Cleanup(server.Close)
	return server
}

func TestSendGauge(t *testing.T) {
	tests := []struct {
		name       string
		value      float64
		statusCode int
		wantErr    bool
	}{
		{name: "successful request", value: 12.5, statusCode: http.StatusOK},
		{name: "zero value", value: 0, statusCode: http.StatusOK},
		{name: "server error", value: 12.5, statusCode: http.StatusInternalServerError, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantBody := map[string]any{"id": "Alloc", "type": "gauge", "value": tt.value}
			server := metricServer(t, wantBody, tt.statusCode)
			client := resty.NewWithClient(server.Client())
			err := SendGauge(client, server.URL, "Alloc", tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendGauge() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSendCounter(t *testing.T) {
	tests := []struct {
		name       string
		delta      int64
		statusCode int
		wantErr    bool
	}{
		{name: "successful request", delta: 3, statusCode: http.StatusOK},
		{name: "zero delta", delta: 0, statusCode: http.StatusOK},
		{name: "server error", delta: 3, statusCode: http.StatusInternalServerError, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantBody := map[string]any{"id": "PollCount", "type": "counter", "delta": float64(tt.delta)}
			server := metricServer(t, wantBody, tt.statusCode)
			client := resty.NewWithClient(server.Client())
			err := SendCounter(client, server.URL, "PollCount", tt.delta)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendCounter() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
