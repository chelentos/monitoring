package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateHandler(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		contentType string
		wantStatus  int
		wantBody    string
	}{
		{
			name: "valid gauge", path: "/update/gauge/Alloc/12.5",
			contentType: "text/plain", wantStatus: http.StatusOK,
		},
		{
			name: "valid counter", path: "/update/counter/PollCount/3",
			contentType: "text/plain", wantStatus: http.StatusOK,
		},
		{
			name: "invalid counter value", path: "/update/counter/PollCount/1.5",
			contentType: "text/plain", wantStatus: http.StatusBadRequest,
			wantBody: "Invalid counter value\n",
		},
		{
			name: "unsupported content type", path: "/update/gauge/Alloc/12.5",
			contentType: "application/json", wantStatus: http.StatusUnsupportedMediaType,
			wantBody: "Invalid type, expected text/plain\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("POST /update/{type}/{name}/{value}", UpdateHandler)
			req := httptest.NewRequest(http.MethodPost, tt.path, nil)
			req.Header.Set("Content-Type", tt.contentType)
			recorder := httptest.NewRecorder()

			mux.ServeHTTP(recorder, req)

			if recorder.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", recorder.Code, tt.wantStatus)
			}
			if got := recorder.Body.String(); got != tt.wantBody {
				t.Errorf("body = %q, want %q", got, tt.wantBody)
			}
		})
	}
}
