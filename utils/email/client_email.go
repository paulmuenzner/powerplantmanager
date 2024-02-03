package email

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

func NewEmailClient(emailClientConfig *ClientConfigData) (client *MailClient, err error) {
	dialer := gomail.NewDialer(emailClientConfig.Host, emailClientConfig.SMTPPort, emailClientConfig.SMTPUsername, emailClientConfig.SMTPPassword)
	if dialer == nil {
		return nil, fmt.Errorf("Failed to create dialer in 'NewEmailClient()' using email client config data.")
	}

	return &MailClient{MyEmailClient: dialer}, nil
}
