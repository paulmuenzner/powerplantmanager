package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (client *Client) UpdateOneInMongo(databaseName string, filter bson.M, update bson.M, collection string) (*mongo.UpdateResult, error) {
	var result *mongo.UpdateResult
	// Create a session for the database
	session, err := client.MongoDB.StartSession()
	if err != nil {
		return result, err
	}
	defer session.EndSession(context.Background())

	// Select the database and collection
	db := client.MongoDB.Database(databaseName)
	col := db.Collection(collection)

	// Insert the data into the collection
	result, errUpdate := col.UpdateOne(context.TODO(), filter, update)
	if errUpdate != nil {
		return nil, fmt.Errorf("Error when updating documents in collection '%s' of database '%s' in 'UpdateOneInMongo()' using 'UpdateOne()'. Error: %v", collection, databaseName, err)
	}

	// Return the ID of the inserted document
	return result, nil
}
