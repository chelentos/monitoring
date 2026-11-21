package handler

import (
	"bytes"
	"encoding/json"
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

			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(strconv.FormatFloat(val, 'f', -1, 64)))
			return
		} else {
			val, err := reader.GetCounter(metricName)

			if err != nil {
				writeMetricReadError(w, err, "Counter not found")
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(strconv.FormatInt(val, 10)))
			return
		}
	}
}

func GetMetricJSONHandler(reader MetricsGetService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data models.Metrics
		var buf bytes.Buffer

		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err = json.Unmarshal(buf.Bytes(), &data); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if data.MType != models.Gauge && data.MType != models.Counter {
			http.Error(w, "Invalid metric type", http.StatusBadRequest)
			return
		}

		if data.MType == models.Gauge {
			val, err := reader.GetGauge(data.ID)

			if err != nil {
				writeMetricReadError(w, err, "Gauge not found")
				return
			}

			resp, err := json.Marshal(models.Metrics{
				ID:    data.ID,
				MType: data.MType,
				Value: &val,
			})

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(resp)
			return
		} else {
			val, err := reader.GetCounter(data.ID)

			if err != nil {
				writeMetricReadError(w, err, "Counter not found")
				return
			}

			resp, err := json.Marshal(models.Metrics{
				ID:    data.ID,
				MType: data.MType,
				Delta: &val,
			})

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(resp)
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
