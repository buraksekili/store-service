package emailer

import (
	"fmt"
	"net/smtp"

	"github.com/buraksekili/store-service/pkg/logger"
)

// Emailer consists of adequate functionalities to send emails
// with plain auth.
type Emailer struct {
	config *SMTPConfig
	auth   smtp.Auth
	log    logger.Logger
}

// New returns a new email to send emails.
func New(sc *SMTPConfig) *Emailer {
	e := &Emailer{config: sc}
	return e
}

// LoginAuth returns a new plain auth for SMTP.
// CRAM-MD5 auth is not supported at the moment.
func LoginAuth(identity, username, password, host string) smtp.Auth {
	return smtp.PlainAuth(identity, username, password, host)
}

// SendEmail sends an email based on Emailer's config field.
func (e *Emailer) SendEmail(content Email) error {
	smtpHost := e.config.SMTPHost
	smtpPort := e.config.SMTPPort
	e.config.Content = content
	message := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", content.To, content.Subject, content.Message))

	e.auth = LoginAuth("", e.config.From, e.config.Password, e.config.SMTPHost)
	return smtp.SendMail(fmt.Sprintf("%s:%s", smtpHost, smtpPort), e.auth, e.config.From, []string{content.To}, message)
}
