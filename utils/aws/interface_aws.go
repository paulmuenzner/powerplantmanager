package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// ///////////////////////////////////////////////////////////
// Setup interface for AWS S3 repository utilizing Dependency Injection
// /////////////////////
type S3Repository interface {
	UploadFile(bucketName string, objectKey string, fileBytes []byte) error
	BucketExists(bucketName string) (bucketExists bool, err error)
	DeleteObjects(bucketName string, objectKeys []string) error
	S3ObjectExists(objectKey, bucketName string) (objectExists bool, err error)
	ChangeObjectName(bucketName, oldObjectKey, newObjectKey string) error
}

type S3Client struct {
	Client *s3.Client // Requires AWS SDK setup for actual usage
}

type ClientConfigData struct {
	AwsRegion      string
	AwsAccessKeyID string
	AwsSecretKey   string
}

type MethodInterface struct {
	RepositoryInterfaceS3 S3Repository
}

func NewAwsMethodInterface(s3Client *S3Client) *MethodInterface {
	return &MethodInterface{RepositoryInterfaceS3: s3Client}
}

func GetAwsMethods(awsClientConfig *ClientConfigData) (awsClientMethods *MethodInterface, err error) {
	// Setup AWS S3 client dependency
	client, err := NewAwsClient(awsClientConfig)
	if err != nil {
		return nil, fmt.Errorf("Couldn't create S3 client in 'AwsProductionClient()' with 'NewAwsClient()'. Error: %v", err)
	}
	awsClientMethods = NewAwsMethodInterface(client)
	return awsClientMethods, nil
}
