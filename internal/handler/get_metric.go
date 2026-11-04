package handler

import (
	"net/http"
	"strconv"

	models "github.com/chelentos/monitoring/internal/model"
)

func GetMetricHandler(storage MetricsStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
			val, ok := storage.GetGauge(metricName)

			if !ok {
				http.Error(w, "Gauge not found", http.StatusNotFound)
				return
			}

			w.Write([]byte(strconv.FormatFloat(val, 'f', -1, 64)))
			return
		} else {
			val, ok := storage.GetCounter(metricName)

			if !ok {
				http.Error(w, "Counter not found", http.StatusNotFound)
				return
			}

			w.Write([]byte(strconv.FormatInt(val, 10)))
			return
		}
	}
}
