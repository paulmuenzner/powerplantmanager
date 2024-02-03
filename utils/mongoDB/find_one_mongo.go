package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (client *Client) FindOneInMongo(databaseName string, filter bson.M, collection string, sort bson.D, result interface{}) (foundOne bool, err error) {
	// Create a session for the database
	session, err := client.MongoDB.StartSession()
	if err != nil {
		return false, err
	}
	defer session.EndSession(context.Background())

	// Select the database and collection
	db := client.MongoDB.Database(databaseName)
	col := db.Collection(collection)

	// Set options for FindOne
	options := options.FindOne().SetSort(sort)

	// Execute query
	err = col.FindOne(context.Background(), filter, options).Decode(result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, fmt.Errorf("%v", err)
	}

	return true, nil
}
