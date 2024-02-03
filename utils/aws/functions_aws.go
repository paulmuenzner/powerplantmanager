package aws

import (
	"bytes"
	"context"
	"fmt"

	files "github.com/paulmuenzner/powerplantmanager/utils/files"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go/aws/awserr"
)

// /////////////////////////////////////////////////////////////////////
// //// UPLOADER
// Upload object to S3
func (client *S3Client) UploadFile(bucketName string, objectKey string, fileBytes []byte) error {

	// Determine size of file to upload
	fileSize, err := files.GetSizeOfByteSlice(fileBytes)
	if err != nil {
		return fmt.Errorf("Couldn't define file size for object key '%s' in 'UploadFile()' with 'GetSizeOfByteSlice()'. Error: %v", objectKey, err)
	}

	// If file size larger than 11MB stream file in chunks with uploadLargeObjectToS3().
	// ! Minimum file size 5MB to be able to use uploadLargeObjectToS3() according to AWS
	if fileSize < 11*1024*1024 {
		err := client.uploadSmallObjectToS3(bucketName, objectKey, fileBytes)
		if err != nil {
			return fmt.Errorf("Error uploading file in 'UploadFile()' with 'uploadSmallObjectToS3' for object key '%s'. Error: %v", objectKey, err)
		}
	} else {
		err := client.uploadLargeObjectToS3(bucketName, objectKey, fileBytes)
		if err != nil {
			return fmt.Errorf("Error uploading file in 'UploadFile()' with 'uploadLargeObjectToS3' for object key '%s'. Error: %v", objectKey, err)
		}
	}
	return nil
}

// /////////////////////////////////////////////////////////////////////////////////
// Upload files breaks large data into parts and uploads the parts concurrently
// ////////////////////////////////////////////////////////////////////////////
func (client *S3Client) uploadLargeObjectToS3(bucketName string, objectKey string, fileBytes []byte) error {

	// Define chunk sizes
	var partMiBs int64 = 10
	uploader := manager.NewUploader(client.Client, func(u *manager.Uploader) {
		u.PartSize = partMiBs * 1024 * 1024
	})

	// Stream the file in chunks directly to the uploader
	_, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   bytes.NewReader(fileBytes), // Use the file directly as the input stream
	})

	if err != nil {
		return fmt.Errorf("Couldn't upload large object in 'uploadLargeObjectToS3()' with 'uploader.Upload()' to bucket %v with object key:%v. Here's why: %v",
			bucketName, objectKey, err)
	}

	return err
}

// ////////////////////////////////////////////////////////
// Upload file
// ///////////
func (client *S3Client) uploadSmallObjectToS3(bucketName string, objectKey string, fileBytes []byte) error {
	_, err := client.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   bytes.NewReader(fileBytes),
	})
	if err != nil {
		return fmt.Errorf("Couldn't upload file with object key %s to bucket %s. Error in 'uploadSmallObjectToS3' from 'awsS3.Client.PutObject()'. Error: %v", objectKey, bucketName, err)
	}
	return err
}

// ////////////////////////////////////////////////////////////
// Validate if S3 bucket exists
// ////////////////////////////////
func (client *S3Client) BucketExists(bucketName string) (bucketExists bool, err error) {
	_, err = client.Client.HeadBucket(context.TODO(), &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})

	bucketExists = true
	if err != nil {
		err = fmt.Errorf("Either no access to bucket %s or another error determined in 'BucketExists' with 'HeadBucket()'. Error: %v", bucketName, err)
		bucketExists = false
	}

	return bucketExists, err
}

// ////////////////////////////////////////////////////////////
// ////////////////////////////////////////////////////////////
// Delete object from aws S3 bucket
func (client *S3Client) DeleteObjects(bucketName string, objectKeys []string) error {

	var objectIds []types.ObjectIdentifier
	for _, key := range objectKeys {
		objectIds = append(objectIds, types.ObjectIdentifier{Key: aws.String(key)})
	}
	_, err := client.Client.DeleteObjects(context.TODO(), &s3.DeleteObjectsInput{
		Bucket: aws.String(bucketName),
		Delete: &types.Delete{Objects: objectIds},
	})
	if err != nil {
		return fmt.Errorf("Couldn't delete objects from bucket %v. Here's why: %v", bucketName, err)
	}
	return err
}

// ////////////////////////////////////////////////////////////
// Check if object in S3 bucket exists
// ///////////////////////////////////
func (client *S3Client) S3ObjectExists(objectKey, bucketName string) (objectExists bool, err error) {
	// Create a HeadObjectInput with the specified key and bucket
	input := &s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}

	// Execute the HeadObject operation
	_, err = client.Client.HeadObject(context.TODO(), input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "NotFound":
				return false, nil
			default:
				return false, err
			}
		}
		return false, err
	}
	return true, nil
}

// ///////////////////////////////////////////////////////////////////////
// List all virtual folders inside a virtual S3 folder (folderPrefix)
// ///////////////////////////////////
func (client *S3Client) ChangeObjectName(bucketName, oldObjectKey, newObjectKey string) error {
	// Copy the object to the new key.
	_, err := client.Client.CopyObject(context.TODO(), &s3.CopyObjectInput{
		Bucket:     aws.String(bucketName),
		CopySource: aws.String(bucketName + "/" + oldObjectKey),
		Key:        aws.String(newObjectKey),
	})
	if err != nil {
		return err
	}

	// Delete the old object.
	objectKeys := []string{oldObjectKey}
	if err := client.DeleteObjects(bucketName, objectKeys); err != nil {
		return err
	}

	fmt.Printf("Object name successfully changed in S3: %s -> %s\n", oldObjectKey, newObjectKey)
	return nil
}
