package controllers

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/gocolly/colly"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Server : the base elements that make the service
type Server struct {
	DB                 *mongo.Database
	InverterCollector  *colly.Collector
	TelemetryCollector *colly.Collector
}

// Initialize : prepares the service to launch
func (s *Server) Initialize(DBHost, DBPort, DBUser, DBPassword, DBDatabase string) error {
	// Connects with the database
	mongoURI := fmt.Sprintf("mongodb://%v:%v@%v:%v/%v", DBUser, DBPassword, DBHost, DBPort, DBDatabase)
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return err
	}
	s.DB = client.Database(DBDatabase)
	// Checks if first connection
	colls, err := s.DB.ListCollectionNames(ctx, bson.D{})
	if len(colls) == 0 {
		if err := s.DBSetup(ctx); err != nil {
			return err
		}
	}
	// Creates the collectors
	s.InverterCollector = colly.NewCollector()
	s.TelemetryCollector = colly.NewCollector()
	return nil
}

// Terminate : closes connections and ends the service
func (s *Server) Terminate() error {
	// Disconnects from DB
	if err := s.DB.Client().Disconnect(context.Background()); err != nil {
		return err
	}
	return nil
}

// Run : runs the service and recovers errors
func (s *Server) Run() {
	defer s.Terminate()
	// Runs collector routines
	go s.InverterAcquisition()
	go s.TelemetryDataAcquisition()
	// Exits on SIGINT
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
}

// DBSetup : setups the DB collections in the first launch
func (s *Server) DBSetup(ctx context.Context) error {
	// Creates collections with existence rules
	opts := options.CreateCollection()
	opts.SetCapped(true)
	opts.SetSizeInBytes(1e11)
	if err := s.DB.CreateCollection(ctx, "inverters", opts); err != nil {
		return err
	}
	if err := s.DB.CreateCollection(ctx, "telemetryData", opts); err != nil {
		return err
	}
	iCol := s.DB.Collection("inverters")
	tCol := s.DB.Collection("telemetryData")
	// Creates unique indexes
	iMod := mongo.IndexModel{
		Keys: bson.M{
			"serial": -1,
		},
		Options: options.Index().SetUnique(true),
	}
	tMod := mongo.IndexModel{
		Keys: bson.M{
			"serial":            -1,
			"lastTelemetryTime": -1,
		},
		Options: options.Index().SetUnique(true),
	}
	if _, err := iCol.Indexes().CreateOne(ctx, iMod); err != nil {
		return err
	}
	if _, err := tCol.Indexes().CreateOne(ctx, tMod); err != nil {
		return err
	}
	return nil
}
