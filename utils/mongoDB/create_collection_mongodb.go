package mongodb

import (
	"context"
)

func (client *Client) CreateNewCollection(databaseName, collectionName string) error {
	// Access the specified database
	database := client.MongoDB.Database(databaseName)

	// Create a new collection
	err := database.CreateCollection(context.Background(), collectionName)
	if err != nil {
		return err
	}

	return nil
}
