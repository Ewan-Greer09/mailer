package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/Ewan-Greer09/mailer/services/emailer"
	_ "github.com/a-h/templ"
)

type AppConfig struct {
	Address     string `json:"ADDRESS"`
	S3AccessKey string `json:"S3_ACCESS_KEY"`
	S3SecretKey string `json:"S3_SECRET_KEY"`
	S3ViewURL   string `json:"S3_VIEW_URL"`
	S3HostURL   string `json:"S3_HOST_URL"`

	JWTSecretKey string `json:"JWT_SECRET_KEY"`
	MongoURI     string `json:"MONGO_URI"`

	SmtpEmail    string `json:"SMTP_EMAIL"`
	SmtpPassword string `json:"SMTP_PASSWORD"`
}

func main() {
	cfg := loadConfig()

	e := echo.New()
	e.HideBanner = true

	store := emailer.NewMongoStore(cfg.MongoURI)
	defer store.Close(context.Background())

	templater := emailer.NewEmailTemplater()

	uploader := emailer.NewS3Uploader(cfg.S3HostURL, cfg.S3ViewURL)

	handler := emailer.NewHandler(emailer.NewEmailService(cfg.SmtpEmail, cfg.SmtpPassword), store, templater, uploader)
	MountRoutes(e, handler, *cfg)

	err := e.Start(cfg.Address)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("server stopped", "err", err)
}

func MountRoutes(e *echo.Echo, h *emailer.Handler, cfg AppConfig) {
	e.Use(
		middleware.RequestID(),
		middleware.Logger(),
		middleware.Recover(),
		echojwt.JWT([]byte(cfg.JWTSecretKey)),
	)

	api := e.Group("/api")
	api.POST("/send/:communication_type", h.Send)
	api.GET("/:communication_uuid", h.Retrieve)
}

func loadConfig() *AppConfig {
	err := godotenv.Load(".env")
	if err != nil {
		slog.Warn("loadConfig", "err", err)
	}

	return &AppConfig{
		Address:      os.Getenv("ADDRESS"),
		S3AccessKey:  os.Getenv("S3_ACCESS_KEY"),
		MongoURI:     os.Getenv("MONGO_URI"),
		JWTSecretKey: os.Getenv("JWT_SECRET_KEY"),
		S3ViewURL:    os.Getenv("S3_VIEW_URL"),
		S3HostURL:    os.Getenv("S3_HOST_URL"),
		SmtpEmail:    os.Getenv("SMTP_EMAIL"),
		SmtpPassword: os.Getenv("SMTP_PASSWORD"),
	}
}
