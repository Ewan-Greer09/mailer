package frontend

import (
	"log/slog"
	"os"

	"github.com/Ewan-Greer09/mailer/services/frontend/views/root"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	logger *slog.Logger
}

func New() *Handler {
	return &Handler{
		logger: slog.New(
			slog.NewJSONHandler(
				os.Stderr,
				&slog.HandlerOptions{
					AddSource: true,
				},
			),
		).With("service", "frontend"),
	}
}

func (h *Handler) Root(c echo.Context) error {
	h.logger.InfoContext(c.Request().Context(), "called 'Root()'")
	return root.Page().Render(c.Request().Context(), c.Response())
}
