package api

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/rjmalves/cpid-solar-telemetry/api/models"
	"github.com/rjmalves/cpid-solar-telemetry/api/seed"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestListTelemetryData(t *testing.T) {
	ctx := context.Background()
	// Removes all data in the collection
	if err := s.RefreshTelemetryDataCollection(ctx); err != nil {
		log.Fatalf("Error refreshing the DB: %v", err)
	}
	// Seeds the collection for testing data
	if err := seed.LoadTelemetryData(s.DB); err != nil {
		log.Fatalf("Error seeding the DB: %v", err)
	}
	// Verifies the data in DB
	invs, err := models.ListTelemetryData(s.DB, bson.M{})
	if err != nil {
		t.Errorf("Error while listing telemetry data in DB: %v\n", err)
		return
	}
	assert.Equal(t, 300, len(invs))
}

func TestDetectTelemetryDataNotInDB(t *testing.T) {
	ctx := context.Background()
	// Removes all data in the collection
	if err := s.RefreshTelemetryDataCollection(ctx); err != nil {
		log.Fatalf("Error refreshing the DB: %v", err)
	}
	// Tries to read data with the empty collection
	d := models.TelemetryData{
		Serial:            "INVERTER1",
		LastTelemetryTime: 0,
	}
	if d.AlreadyAcquired(s.DB) {
		t.Errorf("Found data that should not exist in DB\n")
		return
	}
}

func TestDetectTelemetryDataAlreadyInDB(t *testing.T) {
	ctx := context.Background()
	// Removes all data in the collection
	if err := s.RefreshTelemetryDataCollection(ctx); err != nil {
		log.Fatalf("Error refreshing the DB: %v", err)
	}
	// Seeds the collection for testing data
	if err := seed.LoadTelemetryData(s.DB); err != nil {
		log.Fatalf("Error seeding the DB: %v", err)
	}
	// Tries to read data by serial
	d := models.TelemetryData{
		Serial:            "INVERTER1",
		LastTelemetryTime: 0,
	}
	if !d.AlreadyAcquired(s.DB) {
		t.Errorf("Failed to detect data already in the DB\n")
		return
	}
}

func TestCreatingRepeatedTelemetryData(t *testing.T) {
	ctx := context.Background()
	// Removes all data in the collection
	if err := s.RefreshTelemetryDataCollection(ctx); err != nil {
		log.Fatalf("Error refreshing the DB: %v", err)
	}
	// Seeds the collection for testing data
	if err := seed.LoadTelemetryData(s.DB); err != nil {
		log.Fatalf("Error seeding the DB: %v", err)
	}
	// Tries to add repeated data to DB
	d := models.TelemetryData{
		Serial:            "INVERTER1",
		LastTelemetryTime: 0,
	}
	if _, err := d.AddDataToDB(s.DB); err == nil {
		t.Errorf("Should have failed while adding repeated data\n")
		return
	}
}

func TestCreatingNewTelemetryData(t *testing.T) {
	ctx := context.Background()
	// Removes all data in the collection
	if err := s.RefreshTelemetryDataCollection(ctx); err != nil {
		log.Fatalf("Error refreshing the DB: %v", err)
	}
	// Seeds the collection for testing data
	if err := seed.LoadTelemetryData(s.DB); err != nil {
		log.Fatalf("Error seeding the DB: %v", err)
	}
	// Tries to create new data
	d := models.TelemetryData{
		Serial:            "INVERTER4",
		LastTelemetryTime: 0,
	}
	if _, err := d.AddDataToDB(s.DB); err != nil {
		t.Errorf("Failed while adding new data to DB: %v\n", err)
		return
	}
	// List the existing data and checks the amount
	invs, _ := models.ListTelemetryData(s.DB, bson.M{})
	assert.Equal(t, 301, len(invs))
}

func TestCreateTelemetryDataFromScrapper(t *testing.T) {
	ctx := context.Background()
	// Removes all data in the collection
	if err := s.RefreshTelemetryDataCollection(ctx); err != nil {
		log.Fatalf("Error refreshing the DB: %v", err)
	}
	// Visits the static server
	baseURL := fmt.Sprintf("http://%v:%v/", os.Getenv("APP_HOST"), os.Getenv("APP_PORT"))
	tURL := baseURL + s.TelemetryPaths[0] + "/"
	go s.TelemetryCollector.Visit(tURL)
	time.Sleep(100 * time.Millisecond)
	// Checks if the data is in DB
	td, _ := models.ListTelemetryData(s.DB, bson.M{})
	assert.Equal(t, 1, len(td))
}
