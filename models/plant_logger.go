package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Simplified place holder model
// More parameters can be added as per needs
// In some cases it can make sense to break up the structure into more models/collections
type PlantLogger struct {
	ID                 primitive.ObjectID `bson:"_id"`
	VoltageOutput      float64            `bson:"voltage_output" json:"voltage_output" validate:"required"`   // Unit: Volt (V), Symbol: Vdc
	CurrentOutput      float64            `bson:"current_output" json:"current_output" validate:"required"`   // Unit: Ampere (A), Symbol: Idc
	PowerOutput        float64            `bson:"power_output" json:"power_output" validate:"required"`       // Unit: Wattage (W), Symbol: Pdc
	SolarRadiation     float64            `bson:"solar_radiation" json:"solar_radiation" validate:"required"` // Unit: W/m2, Symbol: G
	AmbientTemperature float64            `bson:"t_ambient" json:"t_ambient" validate:"required"`             // Unit: °C, Symbol: Tamb
	ModuleTemperature  float64            `bson:"t_module" json:"t_module" validate:"required"`               // Unit: °C, Symbol: Tmod
	RelativeHumidity   float64            `bson:"rel_humidity" json:"rel_humidity" validate:"required"`       // Rel. humidity a measurement range of 0 to 100% RH
	WindSpeed          float64            `bson:"wind_speed" json:"wind_speed" validate:"required"`           // Unit: m/s, Symbol: Sw
	CreatedAt          time.Time          `bson:"created_at" json:"created_at" validate:"required"`
}
