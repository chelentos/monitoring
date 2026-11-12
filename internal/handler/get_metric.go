package handler

import (
	"errors"
	"net/http"
	"strconv"

	models "github.com/chelentos/monitoring/internal/model"
	"github.com/chelentos/monitoring/internal/service"
)

type MetricsGetService interface {
	GetGauge(name string) (float64, error)
	GetCounter(name string) (int64, error)
}

func GetMetricHandler(reader MetricsGetService) http.HandlerFunc {
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
			val, err := reader.GetGauge(metricName)

			if err != nil {
				writeMetricReadError(w, err, "Gauge not found")
				return
			}

			w.Write([]byte(strconv.FormatFloat(val, 'f', -1, 64)))
			return
		} else {
			val, err := reader.GetCounter(metricName)

			if err != nil {
				writeMetricReadError(w, err, "Counter not found")
				return
			}

			w.Write([]byte(strconv.FormatInt(val, 10)))
			return
		}
	}
}

func writeMetricReadError(w http.ResponseWriter, err error, notFoundMessage string) {
	switch {
	case errors.Is(err, service.ErrEmptyName):
		http.Error(w, "Metric name requested", http.StatusNotFound)
	case errors.Is(err, service.ErrMetricNotFound):
		http.Error(w, notFoundMessage, http.StatusNotFound)
	default:
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
