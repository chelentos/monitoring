package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"

	models "github.com/chelentos/monitoring/internal/model"
	"github.com/go-resty/resty/v2"
)

func sendMetric(client *resty.Client, serverURL string, metric models.Metrics) error {
	data, jsonErr := json.Marshal(metric)
	if jsonErr != nil {
		return jsonErr
	}

	var compressedData bytes.Buffer

	w := gzip.NewWriter(&compressedData)

	_, err := w.Write(data)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	resp, err := client.R().
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Content-Type", "application/json").
		SetBody(&compressedData).
		Post(serverURL + "/update")

	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("server returned %s", resp.Status())
	}

	return nil
}

func SendGauge(client *resty.Client, serverURL string, name string, value float64) error {
	return sendMetric(client, serverURL, models.Metrics{
		ID:    name,
		MType: "gauge",
		Value: &value,
	})
}

func SendCounter(client *resty.Client, serverURL string, name string, value int64) error {
	return sendMetric(client, serverURL, models.Metrics{
		ID:    name,
		MType: "counter",
		Delta: &value,
	})
}
