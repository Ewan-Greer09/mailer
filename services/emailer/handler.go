package emailer

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"os"
)

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
	Subject           string            `json:"subject"`
	ReplyTo           string            `json:"reply_to"`
	Recipient         Recipient         `json:"recipient"`
	MetaData          map[string]string `json:"metadata"`           // data not needed for success, but useful for logging/observability
	MessageDataFields map[string]any    `json:"message_datafields"` // datafields specific to a communication_type
}

type Handler struct {
	sender    Sender
	logger    *slog.Logger
	store     Storer
	templater Templater
	uploader  Uploader
}

func NewHandler(s Sender, storer Storer, templater Templater, uploader Uploader) *Handler {
	return &Handler{
		sender:    s,
		logger:    slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})),
		store:     storer,
		templater: templater,
		uploader:  uploader,
	}
}

func (h Handler) Retrieve(c echo.Context) error {
	commID := c.Param("communication_uuid")
	res, err := h.store.GetEmail(commID)
	if err != nil {
		h.logger.Warn("retrieve", "err", err)
		return err
	}

	_ = c.JSON(200, res)
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

	b, err := h.templater.Template(c.Request().Context(), commType, req.MessageDataFields)
	if err != nil {
		h.logger.Warn("send email", "err", err)
		return err
	}

	uid := uuid.NewString()
	var location = "failed to Upload"

	var wg errgroup.Group
	wg.SetLimit(1)

	wg.Go(func() error {
		url, s3err := h.uploader.Upload(b, fmt.Sprintf("%s-%s.html", commType, uid))
		if s3err != nil {
			h.logger.Warn("send: upload", "err", s3err)
			return err
		}
		location = url
		return nil
	})

	err = h.sender.Send(b, req.Recipient.Email, req.Subject)
	if err != nil {
		h.logger.Warn("send email", "err", err)
		return err
	}

	err = wg.Wait()
	if err != nil {
		h.logger.Warn("Send", "err", err)
	}

	var id string
	switch req.Channel {
	case email:
		id, err = h.store.SaveEmail(EmailRecord{
			CommType: commType,
			ViewURL:  location,
		})
		if err != nil {
			h.logger.Warn("save email", "err", err)
			return err
		}
	case push:
		h.logger.Info("Push is not implemented")
	case sms:
		h.logger.Info("SMS is not implemented")
	default:
		h.logger.Warn("handler", "err", fmt.Sprintf("%s is not a recognised channel", req.Channel))
		return fmt.Errorf("%s is not a recognised channel", req.Channel)
	}

	_ = c.JSON(
		200,
		struct {
			Id string `json:"id"`
		}{
			Id: id,
		},
	)
	return nil
}
