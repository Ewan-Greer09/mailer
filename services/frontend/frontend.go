package frontend

import (
	"embed"
	"log"
	"log/slog"
	"os"

	"github.com/Ewan-Greer09/mailer/services/frontend/views/root"
	"github.com/labstack/echo/v4"
)

//go:embed public
var publicFS embed.FS

type Handler struct {
	logger *slog.Logger
}

func New() *Handler {
	return &Handler{
		logger: slog.New(
			slog.NewJSONHandler(
				os.Stderr,
				&slog.HandlerOptions{
					AddSource: false,
				},
			),
		).With("service", "frontend"),
	}
}

func (*Handler) Static(e *echo.Echo) {
	e.StaticFS("/public", echo.MustSubFS(publicFS, "public/"))
}

func (h *Handler) Root(c echo.Context) error {
	return root.Page(reciepients).Render(c.Request().Context(), c.Response())
}

func (h *Handler) Recipient(c echo.Context) error {
	email := c.FormValue("email")
	reciepients = append(reciepients, email)
	return root.SingleEmailForm().Render(c.Request().Context(), c.Response())
}

var reciepients = []string{"example@example.com"}

func (h *Handler) RecipientList(c echo.Context) error {
	return root.EmailList(reciepients).Render(c.Request().Context(), c.Response())
}

func (h *Handler) DeleteRecipient(c echo.Context) error {
	email := c.Param("email")
	for i, v := range reciepients {
		if v == email {
			reciepients = append(reciepients[:i], reciepients[i+1:]...)
		}
	}

	log.Println("-----------------------------")
	log.Println("Email: ", email)
	log.Println("Recipients: ", reciepients)
	log.Println("-----------------------------")

	return root.EmailList(reciepients).Render(c.Request().Context(), c.Response())
}

func (h *Handler) BatchSend(c echo.Context) error {

	/*
		go func(){
			h.Emailer.BatchSend(contentType, recipients)
		}()

		return root.EmailSendingConfirm().Render(c.Request().Context(), c.Response())
	*/

	return nil
}
