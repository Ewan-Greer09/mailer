package emailer

import (
	"fmt"
	"log/slog"
)

type Templater interface {
	Template(string, map[string]any) ([]byte, error)
}

type EmailTemplater struct{}

func NewEmailTemplater() *EmailTemplater {
	return &EmailTemplater{}
}

func (EmailTemplater) Template(commType string, datafields map[string]any) ([]byte, error) {
	slog.Info("EmailService", "method", "Send", "Content", fmt.Sprint(datafields))

	message := "stubbed template data"
	return []byte(message), nil
}
