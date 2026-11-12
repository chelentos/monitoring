package main

import (
	"log"
	"os"
	"time"

	"github.com/chelentos/monitoring/internal/agent"
	"github.com/chelentos/monitoring/internal/config"
	"github.com/go-resty/resty/v2"
)

func main() {
	cfg, err := config.ParseAgent(os.Args[1:])
	if err != nil {
		log.Fatal(err)
		return
	}

	client := resty.New()
	serverURL := "http://" + cfg.Sender.Address
	counter := 0
	pollTicker := time.NewTicker(cfg.Collector.PollInterval)
	reportTicker := time.NewTicker(cfg.Sender.ReportInterval)

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
					log.Fatal(err)
				}
			}

			err := agent.SendCounter(client, serverURL, "PollCount", int64(counter))
			if err != nil {
				log.Fatal(err)
			} else {
				counter = 0
			}
		}
	}
}
