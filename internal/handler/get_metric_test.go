package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chelentos/monitoring/internal/repository"
)

func TestGetMetricHandler(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		metricType  string
		metricName  string
		contentType string
		wantStatus  int
		wantBody    string
	}{
		{
			name: "existing gauge", method: http.MethodGet,
			metricType: "gauge", metricName: "Alloc", contentType: "text/plain",
			wantStatus: http.StatusOK, wantBody: "12.5",
		},
		{
			name: "existing counter", method: http.MethodGet,
			metricType: "counter", metricName: "PollCount", contentType: "text/plain",
			wantStatus: http.StatusOK, wantBody: "3",
		},
		{
			name: "zero gauge", method: http.MethodGet,
			metricType: "gauge", metricName: "Zero", contentType: "text/plain",
			wantStatus: http.StatusOK, wantBody: "0",
		},
		{
			name: "zero counter", method: http.MethodGet,
			metricType: "counter", metricName: "Zero", contentType: "text/plain",
			wantStatus: http.StatusOK, wantBody: "0",
		},
		{
			name: "missing gauge", method: http.MethodGet,
			metricType: "gauge", metricName: "PollCount", contentType: "text/plain",
			wantStatus: http.StatusNotFound, wantBody: "Gauge not found\n",
		},
		{
			name: "missing counter", method: http.MethodGet,
			metricType: "counter", metricName: "Alloc", contentType: "text/plain",
			wantStatus: http.StatusNotFound, wantBody: "Counter not found\n",
		},
		{
			name: "invalid metric type", method: http.MethodGet,
			metricType: "unknown", metricName: "Alloc", contentType: "text/plain",
			wantStatus: http.StatusBadRequest, wantBody: "Invalid metric type\n",
		},
		{
			name: "missing metric name", method: http.MethodGet,
			metricType: "gauge", contentType: "text/plain",
			wantStatus: http.StatusNotFound, wantBody: "Metric name requested\n",
		},
		{
			name: "unsupported method", method: http.MethodPost,
			metricType: "gauge", metricName: "Alloc", contentType: "text/plain",
			wantStatus: http.StatusMethodNotAllowed,
		},
		{
			name: "unsupported content type", method: http.MethodGet,
			metricType: "gauge", metricName: "Alloc", contentType: "application/json",
			wantStatus: http.StatusUnsupportedMediaType, wantBody: "Invalid type, expected text/plain\n",
		},
		{
			name: "missing content type", method: http.MethodGet,
			metricType: "gauge", metricName: "Alloc",
			wantStatus: http.StatusUnsupportedMediaType, wantBody: "Invalid type, expected text/plain\n",
		},
		{
			name: "content type with charset", method: http.MethodGet,
			metricType: "gauge", metricName: "Alloc", contentType: "text/plain; charset=utf-8",
			wantStatus: http.StatusOK, wantBody: "12.5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := repository.InitMemStorage()
			storage.SetGauge("Alloc", 12.5)
			storage.SetGauge("Zero", 0)
			storage.AddCounter("PollCount", 3)
			storage.AddCounter("Zero", 0)

			req := httptest.NewRequest(tt.method, "/value/"+tt.metricType+"/"+tt.metricName, nil)
			req.SetPathValue("type", tt.metricType)
			req.SetPathValue("name", tt.metricName)
			if tt.contentType != "" {
				req.Header.Set("Content-Type", tt.contentType)
			}
			recorder := httptest.NewRecorder()

			GetMetricHandler(storage).ServeHTTP(recorder, req)

			if recorder.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", recorder.Code, tt.wantStatus)
			}
			if got := recorder.Body.String(); got != tt.wantBody {
				t.Errorf("body = %q, want %q", got, tt.wantBody)
			}
		})
	}
}
