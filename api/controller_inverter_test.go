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
)

func TestInverterAcquisition(t *testing.T) {
	ctx := context.Background()
	// Removes all data in the collection
	if err := s.RefreshInverterCollection(ctx); err != nil {
		log.Fatalf("Error refreshing the DB: %v", err)
	}
	// Lets the inverter acquisition run for 5s
	baseURL := fmt.Sprintf("http://%v:%v/", os.Getenv("APP_HOST"), os.Getenv("APP_PORT"))
	c := make(chan bool)
	go s.InverterAcquisition(baseURL, 1, c)
	time.Sleep(5 * time.Second)
	c <- true
	// Verifies the inverter in DB
	i := models.Inverter{
		Serial: "7E1504FE-95",
	}
	err := i.ReadInverter(s.DB)
	if err != nil {
		t.Errorf("Error while reading inverter in DB: %v\n", err)
		return
	}
	// Checks if the DB has repeated inverters
	invs, _ := models.ListInverters(s.DB)
	assert.Equal(t, 1, len(invs))
}
