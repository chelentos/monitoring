package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chelentos/monitoring/internal/repository"
	"github.com/go-chi/chi/v5"
)

func TestMetricRoundTrip(t *testing.T) {
	tests := []struct {
		name       string
		metricType string
		values     []string
		want       []string
	}{
		{name: "fractional gauge", metricType: "gauge", values: []string{"12.3456789012345"}, want: []string{"12.3456789012345"}},
		{name: "zero gauge", metricType: "gauge", values: []string{"0"}, want: []string{"0"}},
		{name: "negative gauge", metricType: "gauge", values: []string{"-12.5"}, want: []string{"-12.5"}},
		{name: "gauge numeric value", metricType: "gauge", values: []string{"1.25e2"}, want: []string{"125"}},
		{name: "gauge replacement", metricType: "gauge", values: []string{"12.5", "3.25", "0"}, want: []string{"12.5", "3.25", "0"}},
		{name: "zero counter", metricType: "counter", values: []string{"0"}, want: []string{"0"}},
		{name: "counter accumulation", metricType: "counter", values: []string{"3", "2", "-7"}, want: []string{"3", "5", "-2"}},
		{name: "large counter without precision loss", metricType: "counter", values: []string{"9007199254740993"}, want: []string{"9007199254740993"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := repository.InitMemStorage()
			router := chi.NewRouter()
			router.Post("/update/{type}/{name}/{value}", UpdateMetricHandler(storage))
			router.Get("/value/{type}/{name}", GetMetricHandler(storage))

			for i, value := range tt.values {
				update := httptest.NewRequest(http.MethodPost, "/update/"+tt.metricType+"/Metric/"+value, nil)
				written := httptest.NewRecorder()
				router.ServeHTTP(written, update)
				if written.Code != http.StatusOK {
					t.Fatalf("update %d: status = %d, body = %q", i, written.Code, written.Body.String())
				}

				request := httptest.NewRequest(http.MethodGet, "/value/"+tt.metricType+"/Metric", nil)
				response := httptest.NewRecorder()
				router.ServeHTTP(response, request)
				if response.Code != http.StatusOK {
					t.Fatalf("read after update %d: status = %d, body = %q", i, response.Code, response.Body.String())
				}
				if got := response.Body.String(); got != tt.want[i] {
					t.Errorf("read after update %d: metric = %q, want %q", i, got, tt.want[i])
				}
			}
		})
	}
}
