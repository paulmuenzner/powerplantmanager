package email

import (
	"fmt"
	"time"

	"gopkg.in/gomail.v2"
)

// ///////////////////////////////////////////////////////////////////////
// Setup interface for email repository utilizing Dependency Injection
// /////////////////////
type Repository interface {
	EmailRegistrationSuccess(timeStamp time.Time, email string) error
	EmailInformUserFailedLogin(timeStamp time.Time, email string) error
	SendEmail(senderEmail, recipientEmail, subject, body string) error
	EmailRegistrationVerifiedAccount(timeStamp time.Time, email string) error
	EmailNewRegistration(timeStamp time.Time, email, verifyLinkValidMinutes, encryptedVerifyToken string) error
}

type MailClient struct {
	MyEmailClient *gomail.Dialer
}

type ClientConfigData struct {
	Host         string
	SMTPUsername string
	SMTPPassword string
	SMTPPort     int
}

type RepositoryInterface struct {
	RepositoryInterface Repository
}

func NewEmailMetodInterface(emailClient *MailClient) *RepositoryInterface {
	return &RepositoryInterface{RepositoryInterface: emailClient}
}

func GetEmailRepositoryInterface(emailClientConfig *ClientConfigData) (emailClientMethods *RepositoryInterface, err error) {
	// Setup email client dependency
	client, err := NewEmailClient(emailClientConfig)
	if err != nil {
		return nil, fmt.Errorf("Cannot create email client in 'EmailProductionClient()' with 'NewEmailClient()'. Error: %v", err)
	}
	emailClientMethods = NewEmailMetodInterface(client)

	return emailClientMethods, nil
}
