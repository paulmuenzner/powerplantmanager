package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Configuration of the logging requirements for the API logger
// More parameters can be added as needed
type PlantLoggerConfig struct {
	ID                   primitive.ObjectID `bson:"_id"` // Same _id as in PhotovoltaicPlant
	PublicPlantID        string             `bson:"public_plant_id" json:"public_plant_id" validate:"required" unique:"true"`
	IntervalSec          int                `bson:"interval_sec" json:"interval_sec"  validate:"required" unique:"false"` // Important security feature prventing spamming ('rate limiting'). Time window for logging new data (document)
	Key                  string             `bson:"key" json:"key"`
	Secret               string             `bson:"secret" json:"secret"`
	URLID                string             `bson:"url_id" json:"url_id"`
	CollectionNameLogger string             `bson:"collection_name_logger" json:"collection_name_logger" validate:"required" unique:"true"` // Logging of plant measurements is realized with a separate database collection for each plant
	IPWhitelist          []string           `bson:"ip_whitelist" json:"ip_whitelist" unique:"false"`
	CreatedAt            time.Time          `bson:"created_at" json:"created_at" validate:"required"`
}
