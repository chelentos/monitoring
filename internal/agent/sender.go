package agent

import (
	"fmt"
	"net/http"
)

func SendGauge(client *http.Client, serverURL string, name string, value float64) error {
	requestURL := fmt.Sprintf(
		"%s/update/gauge/%s/%g",
		serverURL,
		name,
		value,
	)

	resp, err := client.Post(requestURL, "text/plain", nil)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned %s", resp.Status)
	}

	fmt.Println("gauge", resp.Status)

	return nil
}

func SendCounter(client *http.Client, serverURL string, name string, value int64) error {
	requestURL := fmt.Sprintf(
		"%s/update/counter/%s/%d",
		serverURL,
		name,
		value,
	)

	resp, err := client.Post(requestURL, "text/plain", nil)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned %s", resp.Status)
	}

	fmt.Println("counter", resp.Status)

	return nil
}
