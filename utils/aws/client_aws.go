package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func NewAwsClient(awsClientConfig *ClientConfigData) (client *S3Client, err error) {
	awsRegion := awsClientConfig.AwsRegion
	awsAccessKeyID := awsClientConfig.AwsAccessKeyID
	awsSecretKey := awsClientConfig.AwsSecretKey

	// Load AWS configuration from environment variables
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(awsRegion),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(awsAccessKeyID, awsSecretKey, "")),
	)
	if err != nil {
		return nil, fmt.Errorf("Error loading AWS configuration in 'NewS3Client': %v", err)
	}

	// Use the configuration to create an AWS service client (S3 client)
	return &S3Client{Client: s3.NewFromConfig(cfg)}, nil
}
