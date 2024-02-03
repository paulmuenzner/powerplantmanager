package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

// Check if value for field name (key) exists in collection (eg. is email address xyz already registered)
func (client *Client) IsValueInCollection(databaseName, collectionName, fieldName, fieldValue string) (bool, error) {
	db := client.MongoDB.Database(databaseName)
	col := db.Collection(collectionName)
	filter := bson.M{fieldName: fieldValue}
	count, err := col.CountDocuments(context.Background(), filter)
	return count != 0, err
}
