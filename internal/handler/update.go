package handler

import (
	"mime"
	"net/http"
	"strconv"

	models "github.com/chelentos/monitoring/internal/model"
)

func UpdateMetricHandler(storage MetricsStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil || mediaType != "text/plain" {
			http.Error(w, "Invalid type, expected text/plain", http.StatusUnsupportedMediaType)
			return
		}

		metricType := r.PathValue("type")

		if metricType != models.Gauge && metricType != models.Counter {
			http.Error(w, "Invalid metric type", http.StatusBadRequest)
			return
		}

		metricName := r.PathValue("name")

		if metricName == "" {
			http.Error(w, "Metric name requested", http.StatusNotFound)
			return
		}

		if metricType == models.Gauge {
			value, err := strconv.ParseFloat(r.PathValue("value"), 64)

			if err != nil {
				http.Error(w, "Invalid gauge value", http.StatusBadRequest)
				return
			}

			storage.SetGauge(metricName, value)
			w.WriteHeader(http.StatusOK)
			return
		} else {
			value, err := strconv.ParseInt(r.PathValue("value"), 10, 64)

			if err != nil {
				http.Error(w, "Invalid counter value", http.StatusBadRequest)
				return
			}

			storage.AddCounter(metricName, value)
			w.WriteHeader(http.StatusOK)
			return
		}
	}
}
