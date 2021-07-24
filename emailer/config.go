package emailer

import (
	"errors"
	"os"
)

var (
	// ErrMissingConfig shows that at least one of the environment
	// variables cannot be read from .env file. Please check
	// ./docker/.env file.
	ErrMissingConfig = errors.New("missing config provided")
)

// Email represent an email. It consists of fields that
// constitute content of an email.
type Email struct {
	// To represents the receiver's email address.
	To string

	// Subject is the subject of the email.
	Subject string

	// Message is the content of email message.
	Message string
}

// SMTPConfig represents required fields to instantiate SMTP.
type SMTPConfig struct {
	// From represents the sender's email address.
	From string

	// Password is the password of the email address
	// that is trying to send an email. So, this is the
	// password of the email of `From` field.
	Password string

	// SMTPHost and SMTPPort are required for email service
	// connection. For example, Outlook's SMTP Settings can
	// be found via:
	// https://support.microsoft.com/en-us/office/pop-imap-and-smtp-settings-8361e398-8af4-4e97-b147-6c6c4ac95353
	SMTPHost string
	SMTPPort string

	// Content is the email content that is desired to send.
	Content Email
}

// ExtractSMTPConfig reads S_SMTP_FROM, S_SMTP_HOST, S_SMTP_PORT
// and S_SMTP_PASSWORD environment variables to instantiate SMTP
// for email service. You can check and update the default values
// of these environment variables through ./docker/.env file.
func ExtractSMTPConfig() (*SMTPConfig, error) {
	sc := &SMTPConfig{}
	if v := os.Getenv("S_SMTP_FROM"); v != "" {
		sc.From = v
	} else {
		return nil, ErrMissingConfig
	}

	if v := os.Getenv("S_SMTP_PASSWORD"); v != "" {
		sc.Password = v
	} else {
		return nil, ErrMissingConfig
	}

	if v := os.Getenv("S_SMTP_HOST"); v != "" {
		sc.SMTPHost = v
	} else {
		return nil, ErrMissingConfig
	}

	if v := os.Getenv("S_SMTP_PORT"); v != "" {
		sc.SMTPPort = v
	} else {
		return nil, ErrMissingConfig
	}

	return sc, nil
}
