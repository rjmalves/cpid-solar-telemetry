package api

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/rjmalves/cpid-solar-telemetry/api/models"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestTelemetryDataAcquisition(t *testing.T) {
	ctx := context.Background()
	// Removes all data in the collection
	if err := s.RefreshTelemetryDataCollection(ctx); err != nil {
		log.Fatalf("Error refreshing the DB: %v", err)
	}
	// Lets the inverter acquisition run for 5s
	baseURL := fmt.Sprintf("http://%v:%v/", os.Getenv("APP_HOST"), os.Getenv("APP_PORT"))
	c := make(chan bool)
	go s.TelemetryDataAcquisition(baseURL, 1, c)
	time.Sleep(5 * time.Second)
	c <- true
	// Verifies the telemetry data in DB
	filter := bson.M{
		"serial": "7E1504FE-95",
	}
	data, err := models.ListTelemetryData(s.DB, filter)
	if err != nil {
		t.Errorf("Error while reading data in DB: %v\n", err)
		return
	}
	// Checks if the DB has repeated data
	assert.Equal(t, 1, len(data))
}
