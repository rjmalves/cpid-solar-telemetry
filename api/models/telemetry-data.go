package models

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
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
	oid, _ := res.InsertedID.(primitive.ObjectID)
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

// FromScrapper : fills the telemetry with data from the HTML scrapper
func (t *TelemetryData) FromScrapper(e *colly.HTMLElement) error {
	// Variables to only acquire information once
	foundModule := false
	foundLastTelemetry := false
	foundLastOutputVoltage := false
	foundLastInputVoltage := false
	foundLastInputCurrent := false
	// Helper variables to find information in HTML
	const SERIAL1 = "row heading no-gutters justify-content-center"
	const SERIAL2 = "heading-info-font"
	const OTHER1 = "row no-gutters align-items-center row-margin"
	const OTHER2 = "col-5 title-font"
	const OTHER3 = "col-5 setting-font text-right"
	// Looks in all divs with classes
	e.ForEach("div[class]", func(_ int, el *colly.HTMLElement) {
		// Looks for the telemetry serial
		if el.Attr("class") == SERIAL1 {
			el.ForEach("span", func(_ int, ele *colly.HTMLElement) {
				if ele.Attr("class") == SERIAL2 {
					s := strings.Split(ele.Text, " ")
					if len(s) == 2 {
						t.Serial = s[1]
					}
				}
			})
		}
		// Looks for other telemetry attributes
		if el.Attr("class") == OTHER1 {
			divData := ""
			el.ForEach("div", func(_ int, ele *colly.HTMLElement) {
				if ele.Attr("class") == OTHER2 {
					divData = ele.Text
				} else if ele.Attr("class") == OTHER3 {
					switch divData {
					case "Módulo":
						if !foundModule {
							foundModule = true
							t.Module = ele.Text
						}
					case "Última telemetria":
						if !foundLastTelemetry {
							foundLastTelemetry = true
							layout := "Jan-02-2006, 15:04:05"
							lt, err := time.Parse(layout, ele.Text)
							if err != nil {
								return
							}
							t.LastTelemetryTime = lt.Unix()
						}
					case "Tensão de Saída":
						if !foundLastOutputVoltage {
							foundLastOutputVoltage = true
							s := strings.Split(ele.Text, " ")
							if len(s) == 2 {
								if v, err := strconv.ParseFloat(s[0], 64); err == nil {
									t.OutputVoltage = v
								}
							}
						}
					case "Tensão de Entrada":
						if !foundLastInputVoltage {
							foundLastInputVoltage = true
							s := strings.Split(ele.Text, " ")
							if len(s) == 2 {
								if v, err := strconv.ParseFloat(s[0], 64); err == nil {
									t.InputVoltage = v
								}
							}
						}
					case "Corrente de Entrada":
						if !foundLastInputCurrent {
							foundLastInputCurrent = true
							s := strings.Split(ele.Text, " ")
							if len(s) == 2 {
								if i, err := strconv.ParseFloat(s[0], 64); err == nil {
									t.InputCurrent = i
								}
							}
						}
					default:
					}
				}
			})
		}
	})
	return nil
}
