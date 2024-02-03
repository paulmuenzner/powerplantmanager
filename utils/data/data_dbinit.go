package data

import (
	config "github.com/paulmuenzner/powerplantmanager/config"
	mongodb "github.com/paulmuenzner/powerplantmanager/utils/mongoDB"
	"log"
)

// CreateIndexes initializes database indexes
func CreateIndexes(mongoDBInterface *mongodb.MethodInterface) error {

	// Create a unique index on the email field on UserAuth collection
	if err := mongoDBInterface.RepositoryInterface.CreateUniqueIndex(config.UserAuthCollectionName, config.DatabaseNameUserAuth, "email", true); err != nil {
		log.Fatal("Error creating unique index for 'email' in user auth collection:", err)
		return err
	}

	// Create a unique index on the email field on plant collection
	if err := mongoDBInterface.RepositoryInterface.CreateUniqueIndex(config.CollectionNamePhotovoltaicPlant, config.DatabaseNamePlants, "name", true); err != nil {
		log.Fatal("Error creating unique index for 'name' in plant collection:", err)
		return err
	}

	// Create a unique index on the email field on PlantLoggerConfig collection
	if err := mongoDBInterface.RepositoryInterface.CreateUniqueIndex(config.CollectionNamePlantLoggerConfig, config.DatabaseNamePlantLoggerConfig, "collection_name_logger", true); err != nil {
		log.Fatal("Error creating unique index for 'collection_name_logger' in plant logger config collection:", err)
		return err
	}

	// Create a unique index on the public_file_id field on file collection
	if err := mongoDBInterface.RepositoryInterface.CreateUniqueIndex(config.CollectionNameFiles, config.DatabaseNameFiles, "public_file_id", true); err != nil {
		log.Fatal("Error creating unique index for 'public_file_id' in file collection:", err)
		return err
	}

	// Create a unique index on the slug field on file collection
	if err := mongoDBInterface.RepositoryInterface.CreateUniqueIndex(config.CollectionNameFiles, config.DatabaseNameFiles, "slug", true); err != nil {
		log.Fatal("Error creating unique index for 'slug' in file collection:", err)
		return err
	}

	return nil
}
