package aws

import (
	"fmt"

	"github.com/paulmuenzner/powerplantmanager/config"
	envHandler "github.com/paulmuenzner/powerplantmanager/utils/env"
)

// Retrieve configuration data (eg. aws region, access key) from .env file for production settings only
// Base parameter for dependency injection of aws client (production)
func S3ProductionConfig() (awsClientConfig *ClientConfigData, bucketName string,
	err error) {
	// Retrieve .env values by keys provided in config file

	// AWS REGION
	awsRegion, err := envHandler.GetEnvValue(config.S3RegionEnv, "")
	if err != nil {
		return nil, "", fmt.Errorf("Cannot retrieve .env value for aws region in 'S3ProductionConfig()'. Env key: %s. No default value has been employed. Error: %v", config.S3RegionEnv, err)
	}

	// AWS ACCESS KEY
	awsAccessKeyID, err := envHandler.GetEnvValue(config.S3AccessKeyEnv, "")
	if err != nil {
		return nil, "", fmt.Errorf("Cannot retrieve .env value for aws access key in 'S3ProductionConfig()'. Env key: %s. No default value has been employed. Error: %v", config.S3AccessKeyEnv, err)
	}

	// AWS SECRET KEY
	awsSecretKey, err := envHandler.GetEnvValue(config.S3SecretKeyEnv, "")
	if err != nil {
		return nil, "", fmt.Errorf("Cannot retrieve .env value for aws secret key in 'S3ProductionConfig()'. Env key: %s. No default value has been employed. Error: %v", config.S3SecretKeyEnv, err)
	}

	// Configure ClientConfigData structure
	awsClientConfig = &ClientConfigData{AwsRegion: awsRegion,
		AwsAccessKeyID: awsAccessKeyID,
		AwsSecretKey:   awsSecretKey}

	// S3 BUCKET
	bucketName, err = envHandler.GetEnvValue(config.S3BucketEnv, "")
	if err != nil {
		return nil, "", fmt.Errorf("Cannot retrieve .env value for Mongo URI in 'S3ProductionConfig()'. Env key: %s. No default value has been employed. Error: %v", config.S3BucketEnv, err)
	}

	return awsClientConfig, bucketName, nil
}
