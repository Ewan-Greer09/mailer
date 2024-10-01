package emailer

import (
	"bytes"
	"context"
	"embed"
	"errors"
	"html/template"
)

var (
	//go:embed templates/**/*.html
	templatesFs embed.FS

	templates = map[string]*template.Template{
		"confirm-email": template.Must(template.ParseFS(templatesFs, "templates/base/confirm-email.html")),
		"chat-invite":   template.Must(template.ParseFS(templatesFs, "templates/base/chat-invite.html")),
		"welcome":       template.Must(template.ParseFS(templatesFs, "templates/base/welcome.html")),
	}
)

type Templater interface {
	Template(context.Context, string, map[string]any) ([]byte, error)
}

type EmailTemplater struct{}

func NewEmailTemplater() *EmailTemplater {
	return &EmailTemplater{}
}

func (EmailTemplater) Template(ctx context.Context, commType string, dataFields map[string]any) ([]byte, error) {
	tmpl, ok := templates[commType]
	if !ok {
		return nil, errors.New("no template found for: " + commType)
	}

	var buf bytes.Buffer
	err := tmpl.Execute(&buf, dataFields)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
