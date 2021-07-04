package email

import "os"

// SMTPConfig represents required fields to instantiate SMTP.
type SMTPConfig struct {
	From     string
	To       []string
	Password string
	SMTPHost string
	SMTPPort string
	Subject  string
	Message  string
}

// ExtractSMTPConfig reads environment variables to instantiate SMTP
// for email service.
func ExtractSMTPConfig() *SMTPConfig {
	sc := &SMTPConfig{}
	if v := os.Getenv("from"); v != "" {
		sc.From = v
	}
	if v := os.Getenv("to"); v != "" {
		sc.To = append(sc.To, v)
	}
	if v := os.Getenv("smtpHost"); v != "" {
		sc.SMTPHost = v
	}
	if v := os.Getenv("smtpPort"); v != "" {
		sc.SMTPPort = v
	}
	if v := os.Getenv("subject"); v != "" {
		sc.Subject = v
	}
	if v := os.Getenv("message"); v != "" {
		sc.Message = v
	}
	return sc
}
