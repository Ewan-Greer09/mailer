package frontend

import (
	"embed"
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

var Recipients = []string{"example@example.com"} // todo: this should be a database (sqlite or mongo)

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

func (h *Handler) HandleRoot(c echo.Context) error {
	return root.Page(Recipients).Render(c.Request().Context(), c.Response())
}

func (h *Handler) HandleAddRecipient(c echo.Context) error {
	email := c.FormValue("email")
	Recipients = append(Recipients, email)
	return root.Email(email).Render(c.Request().Context(), c.Response())
}

func (h *Handler) HandleRecipientList(c echo.Context) error {
	return root.EmailList(Recipients).Render(c.Request().Context(), c.Response())
}

func (h *Handler) HandleDeleteRecipient(c echo.Context) error {
	email := c.Param("email")
	for i, v := range Recipients {
		if v == email {
			Recipients = append(Recipients[:i], Recipients[i+1:]...)
		}
	}

	return nil
}
