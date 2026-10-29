package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/chelentos/monitoring/internal/repository"
)

func TestGetAllMetricHandler(t *testing.T) {
	tests := []struct {
		name     string
		populate func(*repository.MemStorage)
		want     []string
	}{
		{
			name:     "empty storage",
			populate: func(*repository.MemStorage) {},
			want:     []string{"Метрик пока нет", "Счётчиков пока нет"},
		},
		{
			name: "stored metric names and values",
			populate: func(storage *repository.MemStorage) {
				storage.SetGauge("Alloc", 12.5)
				storage.SetGauge("Alloc", -3.25)
				storage.SetGauge("ZeroGauge", 0)
				storage.SetGauge("Shared", 1.5)
				storage.SetGauge("Load<&>", 2.5)
				storage.AddCounter("PollCount", 3)
				storage.AddCounter("PollCount", 2)
				storage.AddCounter("ZeroCounter", 0)
				storage.AddCounter("Shared", -7)
				storage.AddCounter("LargeCounter", 9007199254740993)
			},
			want: []string{
				"Alloc: -3.25", "ZeroGauge: 0", "Shared: 1.5", "Load&lt;&amp;&gt;: 2.5",
				"PollCount: 5", "ZeroCounter: 0", "Shared: -7", "LargeCounter: 9007199254740993",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := repository.InitMemStorage()
			tt.populate(storage)
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			response := httptest.NewRecorder()

			GetAllMetricHandler(storage).ServeHTTP(response, request)

			if response.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
			}
			if got := response.Header().Get("Content-Type"); got != "text/html; charset=utf-8" {
				t.Errorf("Content-Type = %q, want text/html; charset=utf-8", got)
			}
			body := response.Body.String()
			for _, metric := range tt.want {
				if !strings.Contains(body, "<li>"+metric+"</li>") {
					t.Errorf("body does not contain list item %q: %s", metric, body)
				}
			}
			if got := strings.Count(body, "<li>"); got != len(tt.want) {
				t.Errorf("list item count = %d, want %d", got, len(tt.want))
			}
		})
	}
}
