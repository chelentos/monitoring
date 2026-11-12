package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chelentos/monitoring/internal/service"
)

func TestUpdateMetricHandler(t *testing.T) {
	tests := []struct {
		name                             string
		path                             string
		contentType                      string
		wantStatus                       int
		wantBody                         string
		serviceErr                       error
		wantGaugeCalls, wantCounterCalls int
	}{
		{
			name: "valid gauge", path: "/update/gauge/Alloc/12.5",
			contentType: "text/plain", wantStatus: http.StatusOK,
			wantGaugeCalls: 1,
		},
		{
			name: "valid counter", path: "/update/counter/PollCount/3",
			contentType: "text/plain", wantStatus: http.StatusOK,
			wantCounterCalls: 1,
		},
		{
			name: "invalid counter value", path: "/update/counter/PollCount/1.5",
			contentType: "text/plain", wantStatus: http.StatusBadRequest,
			wantBody: "Invalid counter value\n",
		},
		{
			name: "content type is ignored", path: "/update/gauge/Alloc/12.5",
			contentType: "application/json", wantStatus: http.StatusOK,
			wantGaugeCalls: 1,
		},
		{
			name: "missing content type", path: "/update/gauge/Alloc/12.5",
			wantStatus:     http.StatusOK,
			wantGaugeCalls: 1,
		},
		{
			name: "invalid gauge value", path: "/update/gauge/Alloc/bad",
			wantStatus: http.StatusBadRequest, wantBody: "Invalid gauge value\n",
		},
		{
			name: "invalid metric type", path: "/update/unknown/Alloc/12.5",
			wantStatus: http.StatusBadRequest, wantBody: "Invalid metric type\n",
		},
		{
			name: "service rejects name", path: "/update/gauge/Alloc/12.5",
			serviceErr: service.ErrEmptyName, wantGaugeCalls: 1,
			wantStatus: http.StatusNotFound, wantBody: "Metric name requested\n",
		},
		{
			name: "service failure", path: "/update/counter/PollCount/3",
			serviceErr: errors.New("private failure details"), wantCounterCalls: 1,
			wantStatus: http.StatusInternalServerError, wantBody: "Internal server error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updater := &serviceStub{err: tt.serviceErr}
			mux := http.NewServeMux()
			mux.HandleFunc("POST /update/{type}/{name}/{value}", UpdateMetricHandler(updater))
			req := httptest.NewRequest(http.MethodPost, tt.path, nil)
			if tt.contentType != "" {
				req.Header.Set("Content-Type", tt.contentType)
			}
			recorder := httptest.NewRecorder()

			mux.ServeHTTP(recorder, req)

			if recorder.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", recorder.Code, tt.wantStatus)
			}
			if got := recorder.Body.String(); got != tt.wantBody {
				t.Errorf("body = %q, want %q", got, tt.wantBody)
			}
			if updater.gaugeCalls != tt.wantGaugeCalls || updater.counterCalls != tt.wantCounterCalls {
				t.Fatalf("calls: gauge=%d, counter=%d; want %d, %d", updater.gaugeCalls, updater.counterCalls, tt.wantGaugeCalls, tt.wantCounterCalls)
			}
			if tt.wantGaugeCalls > 0 && (updater.name != "Alloc" || updater.value != 12.5) {
				t.Errorf("gauge arguments = %q, %v", updater.name, updater.value)
			}
			if tt.wantCounterCalls > 0 && (updater.name != "PollCount" || updater.delta != 3) {
				t.Errorf("counter arguments = %q, %d", updater.name, updater.delta)
			}
		})
	}
}
