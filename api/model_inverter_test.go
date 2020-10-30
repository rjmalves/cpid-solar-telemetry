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
)

func TestListInverters(t *testing.T) {
	ctx := context.Background()
	// Removes all data in the collection
	if err := s.RefreshInverterCollection(ctx); err != nil {
		log.Fatalf("Error refreshing the DB: %v", err)
	}
	// Seeds the collection for testing data
	if err := seed.LoadInverters(s.DB); err != nil {
		log.Fatalf("Error seeding the DB: %v", err)
	}
	// Verifies the inverters in DB
	invs, err := models.ListInverters(s.DB)
	if err != nil {
		t.Errorf("Error while listing inverters in DB: %v\n", err)
		return
	}
	assert.Equal(t, 3, len(invs))
}

func TestDetectInverterNotInDB(t *testing.T) {
	ctx := context.Background()
	// Removes all data in the collection
	if err := s.RefreshInverterCollection(ctx); err != nil {
		log.Fatalf("Error refreshing the DB: %v", err)
	}
	// Tries to read an inverter with the empty collection
	i := models.Inverter{
		Serial: "INVERTER1",
	}
	if i.AlreadyInDB(s.DB) {
		t.Errorf("Found an inverter that should not exist in DB\n")
		return
	}
}

func TestDetectInverterAlreadyInDB(t *testing.T) {
	ctx := context.Background()
	// Removes all data in the collection
	if err := s.RefreshInverterCollection(ctx); err != nil {
		log.Fatalf("Error refreshing the DB: %v", err)
	}
	// Seeds the collection for testing data
	if err := seed.LoadInverters(s.DB); err != nil {
		log.Fatalf("Error seeding the DB: %v", err)
	}
	// Tries to read an inverter by serial
	i := models.Inverter{
		Serial: "INVERTER1",
	}
	if !i.AlreadyInDB(s.DB) {
		t.Errorf("Failed to detect an inverter already in the DB\n")
		return
	}
}

func TestCreatingRepeatedInverter(t *testing.T) {
	ctx := context.Background()
	// Removes all data in the collection
	if err := s.RefreshInverterCollection(ctx); err != nil {
		log.Fatalf("Error refreshing the DB: %v", err)
	}
	// Seeds the collection for testing data
	if err := seed.LoadInverters(s.DB); err != nil {
		log.Fatalf("Error seeding the DB: %v", err)
	}
	// Tries to add an repeated inverter to DB
	i := models.Inverter{
		Serial: "INVERTER1",
	}
	if _, err := i.AddInverterToDB(s.DB); err == nil {
		t.Errorf("Should have failed while adding an repeated inverter\n")
		return
	}
}

func TestCreatingNewInverter(t *testing.T) {
	ctx := context.Background()
	// Removes all data in the collection
	if err := s.RefreshInverterCollection(ctx); err != nil {
		log.Fatalf("Error refreshing the DB: %v", err)
	}
	// Seeds the collection for testing data
	if err := seed.LoadInverters(s.DB); err != nil {
		log.Fatalf("Error seeding the DB: %v", err)
	}
	// Tries to create a new inverter
	i := models.Inverter{
		Serial: "INVERTER4",
	}
	if _, err := i.AddInverterToDB(s.DB); err != nil {
		t.Errorf("Failed while adding a new inverter to DB: %v\n", err)
		return
	}
	// List the existing inverters and checks the amount
	invs, _ := models.ListInverters(s.DB)
	assert.Equal(t, 4, len(invs))
}

func TestUpdateNonExistingInverter(t *testing.T) {
	ctx := context.Background()
	// Removes all data in the collection
	if err := s.RefreshInverterCollection(ctx); err != nil {
		log.Fatalf("Error refreshing the DB: %v", err)
	}
	// Tries to update a non existing inverter
	i := models.Inverter{
		Serial: "INVERTER1",
	}
	if err := i.UpdateInverterInDB(s.DB); err == nil {
		t.Errorf("Should have failed while updating inverter that didn't exist in DB\n")
		return
	}
}

func TestUpdateExistingInverter(t *testing.T) {
	ctx := context.Background()
	// Removes all data in the collection
	if err := s.RefreshInverterCollection(ctx); err != nil {
		log.Fatalf("Error refreshing the DB: %v", err)
	}
	// Seeds the collection for testing data
	if err := seed.LoadInverters(s.DB); err != nil {
		log.Fatalf("Error seeding the DB: %v", err)
	}
	// Reads the current inverter data
	i := models.Inverter{
		Serial: "INVERTER1",
	}
	i.ReadInverter(s.DB)
	i.Voltage = 380.0
	// Tries to update the inverter
	if err := i.UpdateInverterInDB(s.DB); err != nil {
		t.Errorf("Should have succeeded while updating inverter that existed in DB, but found: %v\n", err)
		return
	}
	// Reads again from DB and compares modified fields
	newi := models.Inverter{
		Serial: "INVERTER1",
	}
	newi.ReadInverter(s.DB)
	assert.Equal(t, i, newi)
}

func TestCreateInverterFromScrapper(t *testing.T) {
	ctx := context.Background()
	// Removes all data in the collection
	if err := s.RefreshInverterCollection(ctx); err != nil {
		log.Fatalf("Error refreshing the DB: %v", err)
	}
	// Visits the static server
	baseURL := fmt.Sprintf("http://%v:%v/", os.Getenv("APP_HOST"), os.Getenv("APP_PORT"))
	iURL := baseURL + s.InverterPaths[0] + "/"
	go s.InverterCollector.Visit(iURL)
	time.Sleep(100 * time.Millisecond)
	// Checks if the inverter is in DB
	i := models.Inverter{
		Serial: "7E1504FE-95",
	}
	if err := i.ReadInverter(s.DB); err != nil {
		t.Errorf("Couldn't create inverter from scrapper\n")
	}
}
