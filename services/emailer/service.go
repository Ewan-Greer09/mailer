package emailer

import (
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

// change
var host = "smtp.gmail.com"
var port = "587"

func (es EmailService) Send(content []byte, to string, subject string) error {
	auth := smtp.PlainAuth("", es.smtpEmail, es.smtpPassword, host)
	return smtp.SendMail(host+":"+port, auth, es.smtpEmail, []string{to}, getMessageString(es.smtpEmail, to, subject, string(content)))
}

func getMessageString(from, to, subject, body string) []byte {
	return []byte("From: " + from + "\r\n" + "To: " + to + "\r\n" + "Subject: " + subject + "\r\n" + "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" + body + "\r\n")
}
