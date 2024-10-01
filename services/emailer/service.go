package emailer

import (
	"fmt"
	"net/smtp"
)

type Sender interface {
	Send([]byte, string, string) error
}

type EmailService struct {
	smtpEmail    string
	smtpPassword string
}

func NewEmailService(smtpEmail, smtpPassword string) *EmailService {
	return &EmailService{
		smtpEmail:    smtpEmail,
		smtpPassword: smtpPassword,
	}
}

var host = "smtp.gmail.com"
var port = "587"

func (s EmailService) Send(content []byte, to string, subject string) error {
	auth := smtp.PlainAuth("", s.smtpEmail, s.smtpPassword, host)
	return smtp.SendMail(host+":"+port, auth, s.smtpEmail, []string{to}, getMessageString(s.smtpEmail, to, subject, content))
}

func getMessageString(from, to, subject string, body []byte) []byte {
	return []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIMI-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n%s\r\n", from, to, subject, body))
}
