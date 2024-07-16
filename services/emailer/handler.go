package emailer

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/sync/errgroup"
	"log"
	"log/slog"
	"os"
	"time"
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
	Channel           channel           `json:"communication_channel" validate:"required"`
	Subject           string            `json:"subject" validate:"required"`
	ReplyTo           string            `json:"reply_to" validate:"required,email"`
	Recipient         Recipient         `json:"recipient" validate:"required"`
	MetaData          map[string]string `json:"metadata"`           // data not needed for success, but useful for logging/observability
	MessageDataFields map[string]any    `json:"message_datafields"` // datafields specific to a communication_type
}

type ApiResponse struct {
	Msg   string `json:"msg"`
	Error string `json:"error"`
}

var v = *validator.New(validator.WithRequiredStructEnabled())

func (s SendEmailRequest) Validate() error {
	err := v.Struct(s)
	if err != nil {
		return err
	}
	return nil
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
		return c.JSON(400, ApiResponse{Msg: "Invalid request details", Error: "invalid request body"})
	}

	err = req.Validate()
	if err != nil {
		return c.JSON(400, ApiResponse{Msg: "Invalid request details", Error: err.Error()})
	}

	b, err := h.templater.Template(c.Request().Context(), commType, req.MessageDataFields)
	if err != nil {
		h.logger.Warn("send email", "err", err)
		return c.JSON(500, ApiResponse{
			Msg:   "Could not parse template",
			Error: "error parsing template",
		})
	}

	uid := uuid.NewString()
	var location = "failed to upload"
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
	start := time.Now()
	err = h.sender.Send(b, req.Recipient.Email, req.Subject)
	log.Printf("diff: %s", time.Since(start))
	if err != nil {
		h.logger.Warn("send email", "err", err)
		return c.JSON(500, ApiResponse{
			Msg:   "Could not send message",
			Error: "error while sending message",
		})
	}

	if err = wg.Wait(); err != nil {
		h.logger.Warn("Send", "err", err) //error comes from wg.Go
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
		}
	case push:
		h.logger.Info("Push is not implemented")
	case sms:
		h.logger.Info("SMS is not implemented")
	default:
		h.logger.Warn("handler", "err", fmt.Sprintf("%s is not a recognised channel", req.Channel))
		return fmt.Errorf("%s is not a recognised channel", req.Channel)
	}

	return c.JSON(200, ApiResponse{
		Msg:   id,
		Error: "",
	})
}
