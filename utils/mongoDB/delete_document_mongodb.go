package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (client *Client) DeleteDocumentMongo(databaseName string, filter bson.M, collection string) (interface{}, error) {

	// Create a session for the database
	session, err := client.MongoDB.StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(context.Background())

	// Select the database and collection
	db := client.MongoDB.Database(databaseName)
	col := db.Collection(collection)

	// Delete document from collection
	result, err := col.DeleteOne(context.TODO(), filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("No matching document found in collection '%s' of database '%s' in 'DeleteDocumentMongo()' using 'DeleteOne()'. Filter: %+v Error: %v", collection, databaseName, filter, err)
		}
		return nil, fmt.Errorf("Error when deleting document in collection '%s' of database '%s' in 'DeleteDocumentMongo()' using 'DeleteOne()'. Filter: %+v Error: %v", collection, databaseName, filter, err)
	}

	return result, nil
}
