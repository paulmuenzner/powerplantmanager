package mongodb

import (
	"context"
	"fmt"
	"net/url"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

func ConnectToMongoDB(mongoDBClientConfig *ClientConfigData) (*Client, error) {

	// Get Mongodb URI
	mongodbURI, err := mongoDBClientConfig.GetURI()
	if err != nil {
		return nil, fmt.Errorf("MongoDB URI setup failed in 'ConnectToMongoDB()' utilizing 'GetURI()'. Error: %v", err)
	}

	commandStarted := []string{}
	cmdMonitor := &event.CommandMonitor{
		Started: func(_ context.Context, evt *event.CommandStartedEvent) {
			commandStarted = append(commandStarted, evt.CommandName)
		},
	}

	clientOptions := options.Client().ApplyURI(mongodbURI).
		// Add your security settings here (e.g., clientOptions.SetAuth)
		SetMaxPoolSize(50).
		SetReadConcern(readconcern.Local()).
		SetWriteConcern(writeconcern.Majority()).
		SetRetryWrites(true).
		SetCompressors([]string{"zstd", "snappy"}).
		SetMonitor(cmdMonitor)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("Client setup MongoDB failed in 'ConnectToMongoDB()' utilizing 'mongo.Connect()'. Error: %v", err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, fmt.Errorf("Error in 'ConnectToMongoDB()' utilizing 'client.Ping()'. Cannot connect to MongoDB. Error: %v", err)
	}

	return &Client{MongoDB: client}, nil
}

// GetURI returns the MongoDB connection URI based on the configuration.
func (mongoDBClientConfig *ClientConfigData) GetURI() (string, error) {
	user := mongoDBClientConfig.Username
	password := mongoDBClientConfig.Password
	port := mongoDBClientConfig.Port
	host := mongoDBClientConfig.Host
	scheme := mongoDBClientConfig.Scheme

	var userName *url.Userinfo
	if user != "" {
		userName = url.UserPassword(user, password)
	} else {
		userName = nil
	}

	u := url.URL{
		Scheme: scheme,
		User:   userName,
		Host:   fmt.Sprintf("%s:%s", host, port),
	}

	uri := u.String()

	if uri == "" {
		return "", fmt.Errorf("MongoDB URI setup failed in 'GetURI()' utilizing 'url.URL{}'. Verify related config data (mongoDBClientConfig).")
	}

	return uri, nil
}
