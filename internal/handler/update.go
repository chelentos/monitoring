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

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

func UpdateMetricJSONHandler(updater MetricsUpdateService) http.HandlerFunc {
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

		if (data.MType == models.Counter && data.Delta == nil) ||
			(data.MType == models.Gauge && data.Value == nil) {
			http.Error(w, "No value or delta", http.StatusBadRequest)
			return
		}

		var updateErr error
		var resp []byte
		var jsonErr error
		switch data.MType {
		case models.Gauge:
			updateErr = updater.UpdateGauge(data.ID, *data.Value)

			resp, jsonErr = json.Marshal(models.Metrics{
				ID:    data.ID,
				MType: data.MType,
				Value: data.Value,
			})
		case models.Counter:
			updateErr = updater.UpdateCounter(data.ID, *data.Delta)

			resp, jsonErr = json.Marshal(models.Metrics{
				ID:    data.ID,
				MType: data.MType,
				Delta: data.Delta,
			})
		default:
			http.Error(w, "Unsupported metric type", http.StatusBadRequest)
			return
		}

		if errors.Is(updateErr, service.ErrEmptyName) {
			http.Error(w, "Metric name requested", http.StatusBadRequest)
			return
		}
		if updateErr != nil {
			http.Error(w, updateErr.Error(), http.StatusInternalServerError)
			return
		}
		if jsonErr != nil {
			http.Error(w, jsonErr.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
}
