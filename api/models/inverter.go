package models

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var inverterCollection = "inverters"

// Inverter : model of an inverter installed in the PV system
type Inverter struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Serial        string             `bson:"serial" json:"serial"`
	Power         float64            `bson:"power" json:"power"`
	Voltage       float64            `bson:"voltage" json:"voltage"`
	Frequency     float64            `bson:"frequency" json:"frequency"`
	Communication bool               `bson:"communication" json:"communication"`
	Status        bool               `bson:"status" json:"status"`
	Switch        bool               `bson:"switch" json:"switch"`
}

// AddInverterToDB : adds info about a inverter to the DB
func (i *Inverter) AddInverterToDB(db *mongo.Database) (primitive.ObjectID, error) {
	ctx := context.Background()
	res, err := db.Collection(inverterCollection).InsertOne(ctx, i)
	if err != nil {
		return primitive.NilObjectID, err
	}
	oid, _ := res.InsertedID.(primitive.ObjectID)
	return oid, nil
}

// UpdateInverterInDB : updates information of an inverter in the DB
func (i *Inverter) UpdateInverterInDB(db *mongo.Database) error {
	ctx := context.Background()
	filter := bson.M{"serial": i.Serial}
	update := bson.M{"$set": i}
	_, err := db.Collection(inverterCollection).UpdateOne(ctx, filter, update)
	return err
}

// DeleteInverterFromDB : deletes an inverter from the DB
func (i *Inverter) DeleteInverterFromDB(db *mongo.Database) error {
	ctx := context.Background()
	filter := bson.M{
		"serial": i.Serial,
	}
	res, err := db.Collection(inverterCollection).DeleteOne(ctx, filter)
	if res.DeletedCount < 1 {
		return fmt.Errorf("Inversor não encontrado")
	}
	return err
}
