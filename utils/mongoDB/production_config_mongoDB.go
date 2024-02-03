package mongodb

import (
	"fmt"

	"github.com/paulmuenzner/powerplantmanager/config"

	envHandler "github.com/paulmuenzner/powerplantmanager/utils/env"
)

// Retrieve configuration data (eg.host, port, etc.) from .env file for production settings only
// Base parameter for dependency injection of mongodb client (production)
func ClientConfig() (mongoDBClientConfig *ClientConfigData, err error) {
	// Retrieve .env values by keys provided in config file

	// Scheme
	mongodbScheme, err := envHandler.GetEnvValue(config.MongoDatabaseSchemeEnv, "")
	if err != nil {
		return nil, fmt.Errorf("Cannot retrieve .env value for MongoDB scheme in 'ClientConfig()'. Env key: %s. No default value has been employed. Error: %v", config.MongoDatabaseSchemeEnv, err)
	}

	// USERNAME
	mongodbUsername, err := envHandler.GetEnvValue(config.MongoDatabaseUsernameEnv, "")
	if err != nil {
		return nil, fmt.Errorf("Cannot retrieve .env value for MongoDB smtp user name in 'ClientConfig()'. Env key: %s. No default value has been employed. Error: %v", config.MongoDatabaseUsernameEnv, err)
	}

	// PASSWORD
	mongodbPassword, err := envHandler.GetEnvValue(config.MongoDatabasePasswordEnv, "")
	if err != nil {
		return nil, fmt.Errorf("Cannot retrieve .env value for MongoDB smtp password in 'ClientConfig()'. Env key: %s. No default value has been employed. Error: %v", config.MongoDatabasePasswordEnv, err)
	}

	// HOST
	mngodbHost, err := envHandler.GetEnvValue(config.MongoDatabaseHostdEnv, "")
	if err != nil {
		return nil, fmt.Errorf("Cannot retrieve .env value for MongoDB smtp port in 'ClientConfig()'. Env key: %s. No default value has been employed. Error: %v", config.MongoDatabaseHostdEnv, err)
	}

	// PORT
	mngodbPort, err := envHandler.GetEnvValue(config.MongoDatabasePortEnv, "")
	if err != nil {
		return nil, fmt.Errorf("Cannot retrieve .env value for MongoDB smtp port in 'ClientConfig()'. Env key: %s. No default value has been employed. Error: %v", config.MongoDatabasePortEnv, err)
	}

	// Configure ClientConfigData structure
	mongodbClientConfig := &ClientConfigData{
		Scheme:   mongodbScheme,
		Username: mongodbUsername,
		Password: mongodbPassword,
		Host:     mngodbHost,
		Port:     mngodbPort,
	}

	return mongodbClientConfig, nil
}
