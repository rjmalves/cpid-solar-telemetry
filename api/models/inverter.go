package models

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var inverterCollection = "inverters"

// Inverter : model of an inverter installed in the PV system
type Inverter struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Serial          string             `bson:"serial" json:"serial"`
	Power           float64            `bson:"power" json:"power"`
	Voltage         float64            `bson:"voltage" json:"voltae"`
	Frequency       float64            `bson:"frequency" json:"frequnc"`
	Communication   bool               `bson:"communication" json:"cmmuniation"`
	Status          bool               `bson:"status" json:"status"`
	Switch          bool               `bson:"switch" json:"switch"`
	EnergyToday     float64            `bson:"energyToday" json:"energyToday"`
	EnergyThisMonth float64            `bson:"energyThisMonth" json:"energyThisMonth"`
	EnergyThisYear  float64            `bson:"energyThisYear" json:"energyThisYear"`
	TotalEnergy     float64            `bson:"totalEnergy" json:"totalEnergy"`
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

// FromScrapper : fills the inverter with data from the HTML scrapper
func (i *Inverter) FromScrapper(e *colly.HTMLElement) error {
	// Variables to only acquire information once
	foundPower := false
	foundVoltage := false
	foundFreq := false
	foundComm := false
	foundStatus := false
	foundSwitch := false
	foundEnergyToday := false
	foundEnergyMonth := false
	foundEnergyYear := false
	foundTotalEnergy := false
	// Helper variables to find information in HTML
	const SERIAL1 = "text-center title-vertical-padding black-font "
	const SERIAL2 = "grey-strong-class font-14"
	const OTHER1 = "row no-gutters white-cover text-center weak-border-top"
	const OTHER2 = "row no-gutters white-cover text-center undefined"
	const OTHER3 = "row no-gutters align-items-center white-cover text-center weak-border-top"
	const OTHER4 = "container-vertical-padding"
	const OTHER5 = "container-vertical-padding-single"
	const OTHER6 = "grey-strong-class"
	const OTHER7 = "font-16 grey-primary-font bold"
	// Looks in all divs with classes
	e.ForEach("div[class]", func(_ int, el *colly.HTMLElement) {
		// Looks for the inverter serial
		if el.Attr("class") == SERIAL1 {
			el.ForEach("span", func(_ int, ele *colly.HTMLElement) {
				if ele.Attr("class") == SERIAL2 {
					s := strings.Split(ele.Text, " ")
					if len(s) == 2 {
						i.Serial = s[1]
					}
				}
			})
		}
		// Looks for other inverter attributes
		if el.Attr("class") == OTHER1 || el.Attr("class") == OTHER2 || el.Attr("class") == OTHER3 {
			el.ForEach("div", func(_ int, ele *colly.HTMLElement) {
				if ele.Attr("class") == OTHER4 || ele.Attr("class") == OTHER5 {
					divData := ""
					ele.ForEach("span", func(_ int, elem *colly.HTMLElement) {
						if elem.Attr("class") == OTHER6 {
							divData = elem.Text
						} else if elem.Attr("class") == OTHER7 {
							switch divData {
							case "Potência":
								if !foundPower {
									foundPower = true
									s := strings.Split(elem.Text, " ")
									if len(s) == 2 {
										if pow, err := strconv.ParseFloat(s[0], 64); err == nil {
											i.Power = pow
										}
									}
								}
							case "Tensão":
								if !foundVoltage {
									foundVoltage = true
									s := strings.Split(elem.Text, " ")
									if len(s) == 2 {
										if vol, err := strconv.ParseFloat(s[0], 64); err == nil {
											i.Voltage = vol
										}
									}
								}
							case "Frequência":
								if !foundFreq {
									foundFreq = true
									s := strings.Split(elem.Text, " ")
									if len(s) == 2 {
										if f, err := strconv.ParseFloat(s[0], 64); err == nil {
											i.Frequency = f
										}
									}
								}
							case "Comunic. c/ Servidor.":
								if !foundComm {
									foundComm = true
									if strings.Index(elem.Text, "S_OK") != -1 {
										i.Communication = true
									}
								}
							case "Status":
								if !foundStatus {
									foundStatus = true
									if strings.Index(elem.Text, "Produção") != -1 {
										i.Status = true
									}
								}
							case "Chave está":
								if !foundSwitch {
									foundSwitch = true
									if strings.Index(elem.Text, "On") != -1 {
										i.Switch = true
									}
								}
							case "Hoje":
								if !foundEnergyToday {
									foundEnergyToday = true
									s := strings.Split(elem.Text, " ")
									if len(s) == 2 {
										if en, err := strconv.ParseFloat(s[0], 64); err == nil {
											i.EnergyToday = en
										}
									}
								}
							case "Este Mês":
								if !foundEnergyMonth {
									foundEnergyMonth = true
									s := strings.Split(elem.Text, " ")
									if len(s) == 2 {
										if en, err := strconv.ParseFloat(s[0], 64); err == nil {
											i.EnergyThisMonth = en
										}
									}
								}
							case "Este Ano":
								if !foundEnergyYear {
									foundEnergyYear = true
									s := strings.Split(elem.Text, " ")
									if len(s) == 2 {
										if en, err := strconv.ParseFloat(s[0], 64); err == nil {
											i.EnergyThisYear = en
										}
									}
								}
							case "Total":
								if !foundTotalEnergy {
									foundTotalEnergy = true
									s := strings.Split(elem.Text, " ")
									if len(s) == 2 {
										if en, err := strconv.ParseFloat(s[0], 64); err == nil {
											i.TotalEnergy = en
										}
									}
								}
							default:
							}
						}
					})
				}
			})
		}
	})
	return nil
}
