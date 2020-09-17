package controllers

import (
	"fmt"
	"time"

	"github.com/gocolly/colly"
	"github.com/rjmalves/cpid-solar-telemetry/api/models"
)

// InverterCollectorConfig : configures the inverter data scrapper
func (s *Server) InverterCollectorConfig() error {
	s.InverterCollector.AllowURLRevisit = true
	// Security parameters
	s.InverterCollector.Limit(&colly.LimitRule{
		Parallelism: 4,
		RandomDelay: 10 * time.Millisecond,
	})
	// When a div with id is found
	s.InverterCollector.OnHTML("div[id]", func(e *colly.HTMLElement) {
		// If not the root div, ignores
		if e.Attr("id") != "root" {
			return
		}
		// Processes the HTML
		i := models.Inverter{}
		i.FromScrapper(e)
		// Adds to DB or updates
		if !i.AlreadyInDB(s.DB) {
			if _, err := i.AddInverterToDB(s.DB); err != nil {
				fmt.Printf("Error while adding inverter: %v\n", err)
			}
		} else {
			if err := i.UpdateInverterInDB(s.DB); err != nil {
				fmt.Printf("Error while updating inverter: %v\n", err)
			}
		}
	})

	// Before making a request print "Visiting ..."
	s.InverterCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})
	return nil
}

// InverterAcquisition : uses the scrapper for acquire inverter data
func (s *Server) InverterAcquisition(baseURL string, iPeriod int64) {
	// Prepares the timers
	iTimers := map[string]int64{}
	for _, i := range s.InverterPaths {
		iTimers[i] = time.Now().Unix()
	}
	// Runs forever
	for {
		// For each inverter
		for _, i := range s.InverterPaths {
			// Checks timeout
			cTime := time.Now().Unix()
			if cTime-iTimers[i] >= iPeriod {
				iTimers[i] = cTime
				iURL := baseURL + i + "/"
				go s.InverterCollector.Visit(iURL)
			}
		}
		time.Sleep(1 * time.Second)
	}
}
