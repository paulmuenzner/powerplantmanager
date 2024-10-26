package config

// General configuration parameter

const (
	URL string = "https://www.example.com"
	//http.Server configuration
	WriteTimeout        int = 20
	ReadTimeout         int = 20
	IdleTimeout         int = 60
	DeleteLogsAfterDays int = 5
	// Email
	EmailSendNotifications   bool   = false // If false, no email notifications at all (error & success)
	EmailProviderUserNameEnv string = "EMAIL_PROVIDER_USERNAME"
	EmailProviderPasswordEnv string = "EMAIL_PROVIDER_PASSWORD"
	EmailProviderSMTPPortEnv string = "EMAIL_PROVIDER_SMTP_PORT"
	EmailProviderHostEnv     string = "EMAIL_PROVIDER_HOST"
	EmailAddressSenderEnv    string = "EMAIL_ADDRESS_SENDER_BACKUP"
	EmailAddressReceiverEnv  string = "EMAIL_ADDRESS_RECEIVER_BACKUP"
	// AWS S3 Production config .env variable names
	S3BucketEnv    string = "AWS_S3_BUCKET_NAME"
	S3RegionEnv    string = "AWS_REGION"
	S3AccessKeyEnv string = "AWS_ACCESS_KEY_ID"
	S3SecretKeyEnv string = "AWS_SECRET_ACCESS_KEY"
	// Plant config
	PlantNameLength    int = 50
	IntervalSecDefault int = 15 * 60 // Default interval, in seconds, for enabling data logging to the plant logger.
)
