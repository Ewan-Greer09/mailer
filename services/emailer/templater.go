package emailer

import (
	"bytes"
	"context"
	"embed"
	"errors"
	"github.com/labstack/gommon/log"
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
	var buf = bytes.NewBuffer(nil)

	tmpl, err := getTemplate(commType)
	if err != nil {
		return nil, err
	}

	err = tmpl.Execute(buf, dataFields)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func getTemplate(commType string) (*template.Template, error) {
	log.Info(commType)
	tmpl, ok := templates[commType]
	if !ok {
		return nil, errors.New("no such template")
	}
	return tmpl, nil
}
