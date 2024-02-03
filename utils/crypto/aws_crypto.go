package crypto

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
)

// InitializeAWSSessionAndKMSClient initializes an AWS session and a KMS client.
func InitializeAWSSessionAndKMSClientAWS() (*session.Session, *kms.KMS) {
	// Initialize AWS session
	sess := session.Must(session.NewSession())

	// Create a KMS client
	kmsClient := kms.New(sess)

	return sess, kmsClient
}

// Parameters
// - keySpec (string): Specify the desired key specification, eg "AES_256"
// - kmsClient refers to an instance of the AWS Key Management Service (KMS) client
func GenerateDataKeyAWS(kmsClient *kms.KMS, keySpec string) (*kms.GenerateDataKeyOutput, error) {
	// Generate a data key without specifying a specific key ID (using the default key)
	result, err := kmsClient.GenerateDataKey(&kms.GenerateDataKeyInput{
		KeySpec: aws.String(keySpec),
	})

	return result, err
}

// RetrieveEncryptionKey retrieves the encryption key from AWS Key Management Service (KMS) using the specified key ID.
func RetrieveEncryptionKeyAWS(keyID string) ([]byte, error) {
	// Create an AWS session
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	// Create a KMS client
	kmsClient := kms.New(sess)

	// Get the key policy for the specified key ID
	keyPolicyInput := &kms.GetKeyPolicyInput{
		KeyId:      aws.String(keyID),
		PolicyName: aws.String("default"), // Specify the default policy name
	}

	keyPolicyResult, err := kmsClient.GetKeyPolicy(keyPolicyInput)
	if err != nil {
		return nil, err
	}

	// Extract the key bytes from the key policy result
	keyBytes := []byte(*keyPolicyResult.Policy)
	return keyBytes, nil
}
