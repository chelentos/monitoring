package main

import (
	"fmt"
	"time"

	"github.com/chelentos/monitoring/internal/agent"
	"github.com/go-resty/resty/v2"
)

func main() {
	parseAgentFlags()

	client := resty.New()
	serverURL := "http://" + flagRunAddr
	counter := 0
	pollTicker := time.NewTicker(time.Duration(flagPollInterval))
	reportTicker := time.NewTicker(time.Duration(flagReportInterval))

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
