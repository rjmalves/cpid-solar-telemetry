package api

import (
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

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
