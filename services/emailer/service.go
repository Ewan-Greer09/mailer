package emailer

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"
)

type Sender interface {
	Send(emailType string) error
}

type EmailService struct{}

func NewEmailService() *EmailService {
	return &EmailService{}
}

func (EmailService) Send(emailType string) error { return nil }

type Handler struct {
	sender Sender
	logger slog.Logger
	store  Storer
}

func NewHandler(s Sender, storer Storer) *Handler {
	return &Handler{
		sender: s,
		logger: *slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})),
		store:  storer,
	}
}

type channel string

const (
	email channel = "email"
	sms   channel = "sms"
	push  channel = "push"
)

type SendEmailRequest struct {
	Channel  channel           `json:"communication_channel"`
	Subject  string            `json:"subject,omitempty"`
	To       string            `json:"to"`
	From     string            `json:"from"`
	ReplyTo  string            `json:"reply_to"`
	MetaData map[string]string `json:"metadata"`
}

func (h Handler) Send(c echo.Context) error {
	var req SendEmailRequest
	commType := c.Param("communication_type")

	err := c.Bind(&req)
	if err != nil {
		h.logger.Warn("Send Email", "err", err)
		return err
	}

	err = h.sender.Send(commType)
	if err != nil {
		h.logger.Warn("send email", "err", err)
		return err
	}

	var id string
	switch req.Channel {
	case email:
		id, err = h.store.SaveEmail(EmailRecord{
			Subject: req.Subject,
			ViewURL: "www.google.com",
		})
		if err != nil {
			h.logger.Warn("save email", "err", err)
			return err
		}
	default:
		h.logger.Warn("handler", "err", fmt.Sprintf("%s is not a recognised channel", req.Channel))
		return fmt.Errorf("%s is not a recognised channel", req.Channel)
	}

	c.Response().Write([]byte(fmt.Sprintf("id: %s", id)))
	c.Response().Status = 200
	return nil
}
