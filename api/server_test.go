package api

import (
	"log"
	"os"
	"testing"

	"github.com/rjmalves/cpid-solar-telemetry/api/tests"
)

var ts = tests.StaticServer{}

func TestMain(m *testing.M) {

	// Starts the static test server
	ts.Initialize()
	ts.Run(os.Getenv("APP_PORT"))

	// Initializes the scrapper
	if err := s.Initialize(os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_DATABASE"),
		os.Getenv("INVERTER_PATHS"),
		os.Getenv("TELEMETRY_PATHS")); err != nil {
		log.Fatalf("Error initializing the service: %v", err)
	}

	os.Exit(m.Run())
}
