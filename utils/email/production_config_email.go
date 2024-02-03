package email

import (
	"fmt"

	config "github.com/paulmuenzner/powerplantmanager/config"
	convert "github.com/paulmuenzner/powerplantmanager/utils/convert"
	envHandler "github.com/paulmuenzner/powerplantmanager/utils/env"
)

// Retrieve configuration data (eg. email provider, smtp port) from .env file for production settings only
// Base parameter for dependency injection of email client (production)
func ProductionConfig() (emailClientConfig *ClientConfigData, err error) {
	// Retrieve .env values by keys provided in config file

	// HOST
	host, err := envHandler.GetEnvValue(config.EmailProviderHostEnv, "")
	if err != nil {
		return nil, fmt.Errorf("Cannot retrieve .env value for host of email provider in 'ProductionConfig()'. Env key: %s. No default value has been employed. Error: %v", config.EmailProviderHostEnv, err)
	}

	// USERNAME
	smtpUsername, err := envHandler.GetEnvValue(config.EmailProviderUserNameEnv, "")
	if err != nil {
		return nil, fmt.Errorf("Cannot retrieve .env value for smtp user name in 'ProductionConfig()'. Env key: %s. No default value has been employed. Error: %v", config.EmailProviderUserNameEnv, err)
	}

	// PASSWORD
	smtpPassword, err := envHandler.GetEnvValue(config.EmailProviderPasswordEnv, "")
	if err != nil {
		return nil, fmt.Errorf("Cannot retrieve .env value for smtp password in 'ProductionConfig()'. Env key: %s. No default value has been employed. Error: %v", config.EmailProviderPasswordEnv, err)
	}

	// SMTP PORT
	smtpPortString, err := envHandler.GetEnvValue(config.EmailProviderSMTPPortEnv, "")
	if err != nil {
		return nil, fmt.Errorf("Cannot retrieve .env value for smtp port in 'ProductionConfig()'. Env key: %s. No default value has been employed. Error: %v", config.EmailProviderSMTPPortEnv, err)
	}
	smtPortInt, err := convert.StringToInt(smtpPortString)
	if err != nil {
		return nil, fmt.Errorf("Cannot convert smtpPortString to int in 'ProductionConfig()'. smtpPortString: %s. No default value has been employed. Error: %v", smtpPortString, err)
	}

	// Configure ClientConfigData structure
	emailClientConfig = &ClientConfigData{Host: host,
		SMTPUsername: smtpUsername,
		SMTPPassword: smtpPassword,
		SMTPPort:     smtPortInt}

	return emailClientConfig, nil
}
