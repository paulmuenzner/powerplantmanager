package config

import "go.mongodb.org/mongo-driver/mongo"

const (
	DatabaseNameUserAuth          string = "PlantDB"
	DatabaseNameFiles             string = "PlantDB"
	DatabaseNamePlants            string = "PlantDB"
	DatabaseNamePlantLoggerConfig string = "PlantDB"
	DatabaseNamePlantLogger       string = "PlantDBLogger"
	// Client config production
	MongoDatabaseSchemeEnv   string = "MONGODB_SCHEME"
	MongoDatabaseUsernameEnv string = "MONGODB_USERNAME"
	MongoDatabasePasswordEnv string = "MONGODB_PASSWORD"
	MongoDatabaseHostdEnv    string = "MONGODB_HOST"
	MongoDatabasePortEnv     string = "MONGODB_PORT"
	// Collection names
	UserAuthCollectionName          string = "user_auth"
	CollectionNameFiles             string = "files"
	CollectionNamePhotovoltaicPlant string = "pv_plants"
	CollectionNamePlantLoggerConfig string = "plant_logger_config"
)

// AppConfig holds the application configuration; here for the mongo connection
type AppConfig struct {
	MongoClient *mongo.Client
}
