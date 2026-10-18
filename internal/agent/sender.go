package agent

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-resty/resty/v2"
)

func SendGauge(client *resty.Client, serverURL string, name string, value float64) error {
	resp, err := client.R().SetPathParams(map[string]string{
		"name":  name,
		"value": strconv.FormatFloat(value, 'f', -1, 64),
	}).SetHeader("Content-Type", "text/plain").Post(serverURL + "/update/gauge/{name}/{value}")

	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("server returned %s", resp.Status())
	}

	return nil
}

func SendCounter(client *resty.Client, serverURL string, name string, value int64) error {
	resp, err := client.R().SetPathParams(map[string]string{
		"name":  name,
		"value": strconv.FormatInt(value, 10),
	}).SetHeader("Content-Type", "text/plain").Post(serverURL + "/update/counter/{name}/{value}")

	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("server returned %s", resp.Status())
	}

	return nil
}
