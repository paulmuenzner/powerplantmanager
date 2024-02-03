package mongobtest

import (
	"context"
	"fmt"
	mongodb "github.com/paulmuenzner/powerplantmanager/utils/mongoDB"
	"log"
	"os"
	"time"

	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func TestMain(m *testing.M) {

	// Load environment variables from .env file before running tests
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)

	}

	// Run the tests
	exitCode := m.Run()

	// Exit with the appropriate exit code
	os.Exit(exitCode)
}

// Connection to real database is tested instead of mock
// Provide .env file in this folder (/utils/mongoDB/test) to enable connection providing common parameter (MONGODB_SCHEME, MONGODB_HOST, ...)
func TestConnectToMongoDB_Success(t *testing.T) {

	// Replace this with your actual MongoDB container configuration
	mongodbURI, err := mongodb.ClientConfig()
	if err != nil {
		fmt.Printf("Test error in 'TestConnectToMongoDB_Success()'. Error: %v", err)
	}

	client, err := mongodb.ConnectToMongoDB(mongodbURI)
	if err != nil {
		t.Fatalf("Error connecting to MongoDB: %v", err)
	}
	// Check if client is not nil
	if client == nil || client.MongoDB == nil {
		t.Fatalf("Invalid client or MongoDB instance")
	}

	// Assert that there is no error during connection setup
	assert.NoError(t, err)

	// Assert that the client is not nil
	assert.NotNil(t, client)

	pingCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Set a reasonable timeout
	defer cancel()

	err = client.MongoDB.Ping(pingCtx, readpref.Primary())
	if err != nil {
		fmt.Printf("Error in 'ConnectToMongoDB()' utilizing 'client.Ping()'. Cannot connect to MongoDB. Error: %v", err)
	}

	// Clean up
	defer client.MongoDB.Disconnect(context.Background())

}
