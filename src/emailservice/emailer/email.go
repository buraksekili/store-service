package emailer

import (
	"errors"
	"fmt"
	"log"
	"net/smtp"

	"github.com/buraksekili/store-service/config/email"
)

// https://gist.github.com/andelf/5118732
type loginAuth struct {
	username, password string
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unkown fromServer")
		}
	}
	return nil, nil
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

type Emailer struct {
	Config email.SMTPConfig
}

func NewEmailer(sc email.SMTPConfig) Emailer {
	return Emailer{sc}
}

func (e *Emailer) SendEmail() error {
	to := e.Config.To
	from := e.Config.From
	password := e.Config.Password
	smtpHost := e.Config.SMTPHost
	smtpPort := e.Config.SMTPPort

	message := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", to[0], e.Config.Subject, e.Config.Message))

	auth := LoginAuth(from, password)
	log.Println("Sending an email to ", to)
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
}
