package handler

import (
	"errors"
	"net/http"
	"strconv"

	models "github.com/chelentos/monitoring/internal/model"
	"github.com/chelentos/monitoring/internal/service"
)

type MetricsUpdateService interface {
	UpdateGauge(name string, value float64) error
	UpdateCounter(name string, delta int64) error
}

func UpdateMetricHandler(updater MetricsUpdateService) http.HandlerFunc {
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

		var updateErr error
		if metricType == models.Gauge {
			value, err := strconv.ParseFloat(r.PathValue("value"), 64)

			if err != nil {
				http.Error(w, "Invalid gauge value", http.StatusBadRequest)
				return
			}

			updateErr = updater.UpdateGauge(metricName, value)
		} else {
			value, err := strconv.ParseInt(r.PathValue("value"), 10, 64)

			if err != nil {
				http.Error(w, "Invalid counter value", http.StatusBadRequest)
				return
			}

			updateErr = updater.UpdateCounter(metricName, value)
		}

		if errors.Is(updateErr, service.ErrEmptyName) {
			http.Error(w, "Metric name requested", http.StatusNotFound)
			return
		}
		if updateErr != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
