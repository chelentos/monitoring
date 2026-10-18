package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/chelentos/monitoring/internal/agent"
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 4 * time.Second
)

func main() {
	client := &http.Client{}
	serverURL := "http://localhost:8080"
	counter := 0
	pollTicker := time.NewTicker(pollInterval)
	reportTicker := time.NewTicker(reportInterval)

	defer pollTicker.Stop()
	defer reportTicker.Stop()

	gauges := agent.Collect()

	for {
		select {
		case <-pollTicker.C:
			gauges = agent.Collect()
			counter += 1
		case <-reportTicker.C:
			for name, value := range gauges {
				err := agent.SendGauge(client, serverURL, name, value)
				if err != nil {
					fmt.Println(err)
				}
			}

			err := agent.SendCounter(client, serverURL, "PollCount", int64(counter))
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
