package mongodb

import (
	"context"

	logger "github.com/paulmuenzner/powerplantmanager/utils/logs"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (client *Client) FindManyInMongo(databaseName string, filter bson.M, collection string, sort bson.D, result interface{}) error {

	// Create a session for the database
	session, err := client.MongoDB.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(context.Background())

	// Select the database and collection
	db := client.MongoDB.Database(databaseName)
	col := db.Collection(collection)

	// Set options for Find
	options := options.Find().SetSort(sort)

	// Insert the data into the collection
	cur, err := col.Find(context.Background(), filter, options)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return err
		}
		logger.GetLogger().Error("Error when querying collection "+collection+" of database "+databaseName+" in 'FindManyInMongo'.  Error: ", err)
		return err
	}
	defer cur.Close(context.Background())

	// Decode the results into the provided result interface
	if err := cur.All(context.Background(), result); err != nil {
		logger.GetLogger().Error("Error when decoding documents "+collection+" of database "+databaseName+" in 'FindManyInMongo'. Error: ", err)
		return err
	}

	return nil
}
