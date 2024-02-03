package mongodb

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

func (client *Client) StartSession() (session mongo.Session, err error) {
	// Create a session for the database
	session, err = client.MongoDB.StartSession()
	if err != nil {
		return nil, fmt.Errorf("Error creating session in 'StartSession()'. Error: %v", err)
	}

	return session, nil
}
