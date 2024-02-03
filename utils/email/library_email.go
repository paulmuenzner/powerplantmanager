package email

import (
	"fmt"
	"time"

	"github.com/paulmuenzner/powerplantmanager/config"
	"github.com/paulmuenzner/powerplantmanager/utils/date"
	envHandler "github.com/paulmuenzner/powerplantmanager/utils/env"
	strings "github.com/paulmuenzner/powerplantmanager/utils/strings"
)

// ///////////////////////////////////////////////////////////////////////////
// EMAILS AUTH REGISTRATION, LOGIN PROCESS
// ///////////////////////////////////////

// Registration request for already verified account/email address. Send warning notification to account owner
func (client *MailClient) EmailRegistrationVerifiedAccount(timeStamp time.Time, email string) error {
	timeStampStringUs := date.TimeStampToUSFormat(timeStamp)
	senderEmailAddress, err := envHandler.GetEnvValue(config.EmailAddressSenderEnv, "") // Feel free to use default value via base_config
	// Log as error if no defaultValue provided in GetEnvValue()
	if err != nil {
		return fmt.Errorf("Error in 'EmailRegistrationVerifiedAccount()' utilizing 'GetEnvValue()' for 'senderEmailAddress'. Cannot retrieve env value. Error: %v", err)
	}

	// Subject
	subject := "Attention! Registration request."

	// Body
	bodyComponents := []string{"<html><body><h2>Attention! Registration request.</h2> <br/><br/> Date registration request: ", timeStampStringUs, "<br/> A registration request has been received using your email address. It appears that you already possess a verified account.", "<br/> <br/>Wasn't you? Please inform us and forward this message to: ", senderEmailAddress, "</body></html>"}
	body := strings.ConcatenateStrings(bodyComponents...)

	// Send
	err = client.SendEmail(senderEmailAddress, email, subject, body)
	if err != nil {
		return fmt.Errorf("Error in 'EmailRegistrationVerifiedAccount()' utilizing 'SendEmail()'. Error: %v", err)
	}
	return nil
}

// Send verification request for new registration request
func (client *MailClient) EmailNewRegistration(timeStamp time.Time, email, verifyLinkValidMinutes, encryptedVerifyToken string) error {
	timeStampStringUs := date.TimeStampToUSFormat(timeStamp)
	senderEmailAddress, err := envHandler.GetEnvValue(config.EmailAddressSenderEnv, "") // Feel free to use default value via base_config
	// Log as error if no defaultValue provided in GetEnvValue()
	if err != nil {
		return fmt.Errorf("Error in 'EmailNewRegistration()' utilizing 'GetEnvValue()' for 'senderEmailAddress'. Cannot retrieve env value. Error: %v", err)
	}

	// Subject
	subject := "Verify your new account."

	// Body
	verificationURL := config.URL + "/auth/verify/" + encryptedVerifyToken
	bodyComponents := []string{"<html><body><h2>Registration request</h2>", "<br/><br/>Request date: ", timeStampStringUs, "<br/><br/>Verify your new account within ", verifyLinkValidMinutes, " minutes. <br/><br/> Verification link: ", "<a href=\"", verificationURL, "\" target=\"_blank\">", verificationURL, "</a></body></html>"}
	body := strings.ConcatenateStrings(bodyComponents...)

	// Send
	err = client.SendEmail(senderEmailAddress, email, subject, body)
	if err != nil {
		return fmt.Errorf("Error in 'EmailNewRegistration()' utilizing 'SendEmail()'. Error: %v", err)
	}
	return nil
}

// Inform account owner on failed login attempt
func (client *MailClient) EmailInformUserFailedLogin(timeStamp time.Time, email string) error {
	timeStampStringUs := date.TimeStampToUSFormat(timeStamp)
	senderEmailAddress, err := envHandler.GetEnvValue(config.EmailAddressSenderEnv, "") // Feel free to use default value via base_config
	// Log as error if no defaultValue provided in GetEnvValue()
	if err != nil {
		return fmt.Errorf("Error in 'EmailInformUserFailedLogin()' utilizing 'GetEnvValue()' for 'senderEmailAddress'. Cannot retrieve env value. Error: %v", err)
	}

	// Subject
	subject := "Failed login attempt!"

	// Body
	bodyComponents := []string{"<html><body><h2>Failed login attempt.</h2> <br/><br/> Date login attempt: ", timeStampStringUs, "<br/> Registered email address: ", email, "<br/> <br/>Wasn't you? Please inform us and forward this message to: ", senderEmailAddress, "</body></html>"}
	body := strings.ConcatenateStrings(bodyComponents...)

	// Send
	err = client.SendEmail(senderEmailAddress, email, subject, body)
	if err != nil {
		return fmt.Errorf("Error in 'EmailInformUserFailedLogin()' utilizing 'SendEmail()'. Error: %v", err)
	}
	return nil
}

// Send email after successfully verified registration
func (client *MailClient) EmailRegistrationSuccess(timeStamp time.Time, email string) error {
	timeStampStringUs := date.TimeStampToUSFormat(timeStamp)

	// Subject
	subjectComponents := []string{"Successful backup: ", timeStampStringUs}
	subject := strings.ConcatenateStrings(subjectComponents...)

	// Body
	bodyComponents := []string{"<html><body><h2>Successful Registration.</h2> <br/><br/> Registration date: ", timeStampStringUs, "<br/> Registered email address: ", email, "<br/> <br/>You can finally register power plants and connect them to our logging database.", "</body></html>"}
	body := strings.ConcatenateStrings(bodyComponents...)

	senderEmailAddress, err := envHandler.GetEnvValue(config.EmailAddressSenderEnv, "") // Feel free to use default value via base_config
	// Log as error if no defaultValue provided in GetEnvValue()
	if err != nil {
		return fmt.Errorf("Error in 'EmailRegistrationSuccess()' utilizing 'GetEnvValue()' for 'senderEmailAddress'. Cannot retrieve env value. Error: %v", err)
	}

	// Send
	err = client.SendEmail(senderEmailAddress, email, subject, body)
	if err != nil {
		return fmt.Errorf("Error in 'EmailRegistrationSuccess()' utilizing 'SendEmail()'. Error: %v", err)
	}
	return nil
}
