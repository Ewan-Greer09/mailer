package emailer

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"
)

type Sender interface {
	Send(any) error
}

type EmailService struct{}

func NewEmailService() *EmailService {
	return &EmailService{}
}

func (es EmailService) Send(content any) error {
	slog.Info("EmailService", "method", "Send", "Content", fmt.Sprint(content))
	return nil
}

type Handler struct {
	sender    Sender
	logger    slog.Logger
	store     Storer
	templater Templater
}

func NewHandler(s Sender, storer Storer, templater Templater) *Handler {
	return &Handler{
		sender:    s,
		logger:    *slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})),
		store:     storer,
		templater: templater,
	}
}

type channel string

const (
	email channel = "email"
	sms   channel = "sms"
	push  channel = "push"
)

type Recipient struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type SendEmailRequest struct {
	Channel           channel           `json:"communication_channel"`
	ReplyTo           string            `json:"reply_to"`
	Recipient         Recipient         `json:"recipient"`
	MetaData          map[string]string `json:"metadata"`           // data not needed for success, but useful for logging/observability
	MessageDatafields map[string]any    `json:"message_datafields"` // datafields specific to a communication_type
}

type Templater interface {
	Template(string, map[string]any) ([]byte, error)
}

type EmailTemplater struct{}

func NewEmailTemplater() *EmailTemplater {
	return &EmailTemplater{}
}

func (EmailTemplater) Template(commType string, datafields map[string]any) ([]byte, error) {
	slog.Info("EmailService", "method", "Send", "Content", fmt.Sprint(datafields))
	return []byte{}, nil
}

func (h Handler) Retrieve(c echo.Context) error {
	commID := c.Param("communication_uuid")

	h.logger.Info("retrieve", "id", commID)

	email, err := h.store.GetEmail(commID)
	if err != nil {
		h.logger.Warn("retrieve", "err", err)
		return err
	}

	b, err := json.Marshal(email)
	if err != nil {
		h.logger.Warn("retrieve", "err", err)
		return err
	}

	c.Response().WriteHeader(200)
	c.Response().Write(b)
	return nil
}

func (h Handler) Send(c echo.Context) error {
	var req SendEmailRequest
	commType := c.Param("communication_type")

	err := c.Bind(&req)
	if err != nil {
		h.logger.Warn("send email", "err", err)
		return err
	}

	b, err := h.templater.Template(commType, req.MessageDatafields)
	if err != nil {
		h.logger.Warn("send email", "err", err)
		return err
	}

	err = h.sender.Send(b)
	if err != nil {
		h.logger.Warn("send email", "err", err)
		return err
	}

	var id string
	switch req.Channel {
	case email:
		id, err = h.store.SaveEmail(EmailRecord{
			CommType: commType,
			ViewURL:  "www.google.com",
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
