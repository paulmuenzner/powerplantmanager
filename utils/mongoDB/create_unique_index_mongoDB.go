package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (client *Client) CreateUniqueIndex(collectionName string, databaseName string, fieldName string, unique bool) error {
	// Create a unique index on the email field
	db := client.MongoDB.Database(databaseName)
	col := db.Collection(collectionName)

	// Set options based on unique parameter
	indexOptions := options.Index().SetUnique(unique)

	// Create the index model
	indexModel := mongo.IndexModel{
		Keys:    map[string]interface{}{fieldName: 1},
		Options: indexOptions,
	}

	_, err := col.Indexes().CreateOne(context.Background(), indexModel)
	return err
}
