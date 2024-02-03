package mongodb

import (
	"context"

	logger "github.com/paulmuenzner/powerplantmanager/utils/logs"

	"go.mongodb.org/mongo-driver/mongo"
)

func (client *Client) DeleteCollectionMongo(databaseName string, collection string) error {

	// Create a session for the database
	session, err := client.MongoDB.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(context.Background())

	// Get a handle for your collection
	col := client.MongoDB.Database(databaseName).Collection(collection)

	// Delecte collection
	err = col.Drop(context.Background())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.GetLogger().Infof("No matching document found in collection '%s' of database '%s' Error: %v", collection, databaseName, err)
			return err
		}
		logger.GetLogger().Errorf("Error when querying collection '%s' of database '%s' Error:  %v", collection, databaseName, err)
		return err
	}

	return nil
}
