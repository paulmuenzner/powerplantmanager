package mongodb

import (
	"context"

	logger "github.com/paulmuenzner/powerplantmanager/utils/logs"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (client *Client) CountDocumentsInMongo(databaseName string, collection string, result interface{}) (int, error) {

	// Create a session for the database
	session, err := client.MongoDB.StartSession()
	if err != nil {
		return 0, err
	}
	defer session.EndSession(context.Background())

	// Select the database and collection
	db := client.MongoDB.Database(databaseName)
	col := db.Collection(collection)

	// Insert the data into the collection
	result, err = col.CountDocuments(context.Background(), bson.D{})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, err
		}
		logger.GetLogger().Errorf("Error when counting documents in collection '%s' of database '%s' Error: %v", collection, databaseName, err)
		return 0, err
	}

	return result.(int), nil
}
