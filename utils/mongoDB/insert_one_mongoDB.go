package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (client *Client) InsertOneToMongo(databaseName string, data interface{}, collection string) (string, error) {

	// Create a session for the database
	session, err := client.MongoDB.StartSession()
	if err != nil {
		return "", err
	}
	defer session.EndSession(context.Background())

	// Select the database and collection
	db := client.MongoDB.Database(databaseName)
	col := db.Collection(collection)

	// Insert the data into the collection
	result, err := col.InsertOne(context.Background(), data)
	if err != nil {
		return "", err
	}

	// Convert the inserted ID to string
	var insertedID string
	switch id := result.InsertedID.(type) {
	case primitive.ObjectID:
		insertedID = id.Hex()
	case string:
		insertedID = id
	default:
		return "", fmt.Errorf("unexpected type for InsertedID in 'InsertOneToMongo()' using : %T", id)
	}

	// Return the ID of the inserted document
	return insertedID, nil
}
