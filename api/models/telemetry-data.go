package models

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var telemetryDataCollection = "telemetryData"

// TelemetryData : PV system state captured by the data acquisition service
type TelemetryData struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Serial            string             `bson:"serial" json:"serial"`
	Module            string             `bson:"module" json:"module"`
	LastTelemetryTime int64              `bson:"lastTelemetryTime" json:"lastTelemetryTime"`
	OutputVoltage     float64            `bson:"outputVoltage" json:"outputVoltage"`
	InputVoltage      float64            `bson:"inputVoltage" json:"inputVoltage"`
	InputCurrent      float64            `bson:"inputCurrent" json:"inputCurrent"`
}

// AlreadyAcquired : checks if a given telemetry data is already in the DB
func (t *TelemetryData) AlreadyAcquired(db *mongo.Database) bool {
	ctx := context.Background()
	filter := bson.M{
		"serial":            t.Serial,
		"lastTelemetryTime": t.LastTelemetryTime,
	}
	res := db.Collection(telemetryDataCollection).FindOne(ctx, filter)
	return res.Err() != nil
}

// AddDataToDB : adds a telemetry read to the DB
func (t *TelemetryData) AddDataToDB(db *mongo.Database) (primitive.ObjectID, error) {
	ctx := context.Background()
	res, err := db.Collection(telemetryDataCollection).InsertOne(ctx, t)
	if err != nil {
		return primitive.NilObjectID, err
	}
	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, fmt.Errorf("Erro na conversão para OID")
	}
	return oid, nil
}

// DeleteDataFromDB : deletes a telemetry read from the DB
func (t *TelemetryData) DeleteDataFromDB(db *mongo.Database) error {
	ctx := context.Background()
	filter := bson.M{
		"serial":            t.Serial,
		"lastTelemetryTime": t.LastTelemetryTime,
	}
	res, err := db.Collection(telemetryDataCollection).DeleteOne(ctx, filter)
	if res.DeletedCount < 1 {
		return fmt.Errorf("Dado de telemetria não encontrado")
	}
	return err
}
