package models

import "go.mongodb.org/mongo-driver/bson/primitive"

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
