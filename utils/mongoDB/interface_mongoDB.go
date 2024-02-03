package mongodb

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// //////////////////////////////////////////////////////////////////////
// Setup interface for database repository utilizing Dependency Injection
// ///////////////////
type Repository interface {
	InsertOneToMongo(databaseName string, data interface{}, collection string) (string, error)
	IsValueInCollection(databaseName string, collectionName string, fieldName, fieldValue string) (bool, error)
	UpdateOneInMongo(databaseName string, filter bson.M, update bson.M, collection string) (*mongo.UpdateResult, error)
	FindOneInMongo(databaseName string, filter bson.M, collection string, sort bson.D, result interface{}) (foundOne bool, err error)
	FindManyInMongo(databaseName string, filter bson.M, collection string, sort bson.D, result interface{}) error
	DeleteDocumentMongo(databaseName string, filter bson.M, collection string) (interface{}, error)
	DeleteCollectionMongo(databaseName string, collection string) error
	CreateNewCollection(databaseName, collectionName string) error
	CountDocumentsInMongo(databaseName string, collection string, result interface{}) (int, error)
	CreateUniqueIndex(collectionName string, databaseName string, fieldName string, unique bool) error
	StartSession() (session mongo.Session, err error)
}

type Client struct {
	MongoDB *mongo.Client
}

type ClientConfigData struct {
	Scheme   string
	Username string
	Password string
	Host     string
	Port     string
}

type MethodInterface struct {
	RepositoryInterface Repository
}

func NewMongoDBMethodInterface(mongoClient *Client) *MethodInterface {
	return &MethodInterface{RepositoryInterface: mongoClient}
}
