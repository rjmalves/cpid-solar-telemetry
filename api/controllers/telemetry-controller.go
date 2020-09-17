package controllers

import (
	"fmt"
	"time"

	"github.com/gocolly/colly"
	"github.com/rjmalves/cpid-solar-telemetry/api/models"
)

// TelemetryDataCollectorConfig : configures the telemetry data scrapper
func (s *Server) TelemetryDataCollectorConfig() error {
	// When a div with id is found
	s.TelemetryCollector.OnHTML("div[id]", func(e *colly.HTMLElement) {
		// If not the root div, ignores
		if e.Attr("id") != "root" {
			return
		}
		t := models.TelemetryData{}
		t.FromScrapper(e)
		// Prints the built data
		fmt.Printf("Telemetry: %v\n", t)
	})

	// Before making a request print "Visiting ..."
	s.TelemetryCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})
	return nil
}

// TelemetryDataAcquisition : uses the scrapper for acquire telemetry data
func (s *Server) TelemetryDataAcquisition(baseURL string, tPeriod int64) {
	// Prepares the timers
	tTimers := map[string]int64{}
	for _, t := range s.TelemetryPaths {
		tTimers[t] = time.Now().Unix()
	}
	// Runs forever
	for {
		// For each telemetry
		for _, t := range s.TelemetryPaths {
			// Checks timeout
			cTime := time.Now().Unix()
			if cTime-tTimers[t] > tPeriod {
				tTimers[t] = cTime
				tURL := baseURL + t + "/"
				go s.TelemetryCollector.Visit(tURL)
			}
		}
		time.Sleep(1 * time.Second)
	}
}
