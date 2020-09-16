package models

import "go.mongodb.org/mongo-driver/bson/primitive"

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
