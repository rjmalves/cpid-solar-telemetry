package api

import (
	"log"
	"os"

	"github.com/rjmalves/cpid-solar-telemetry/api/controllers"
)

var s = controllers.Server{}

// Run : launches the service
func Run() {

	if err := s.Initialize(os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_DATABASE"),
		os.Getenv("INVERTER_PATHS"),
		os.Getenv("TELEMETRY_PATHS")); err != nil {
		log.Fatalf("Error initializing the service: %v", err)
	}

	s.Run(os.Getenv("APP_HOST"),
		os.Getenv("APP_PORT"),
		os.Getenv("INVERTER_ACQ_PERIOD"),
		os.Getenv("TELEMETRY_ACQ_PERIOD"))
}
