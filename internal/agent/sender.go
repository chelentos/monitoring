package agent

import (
	"fmt"
	"net/http"

	models "github.com/chelentos/monitoring/internal/model"
	"github.com/go-resty/resty/v2"
)

func SendGauge(client *resty.Client, serverURL string, name string, value float64) error {
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(models.Metrics{
			ID:    name,
			MType: "gauge",
			Value: &value,
		}).
		Post(serverURL + "/update")

	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("server returned %s", resp.Status())
	}

	return nil
}

func SendCounter(client *resty.Client, serverURL string, name string, value int64) error {
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(models.Metrics{
			ID:    name,
			MType: "counter",
			Delta: &value,
		}).
		Post(serverURL + "/update")

	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("server returned %s", resp.Status())
	}

	return nil
}
