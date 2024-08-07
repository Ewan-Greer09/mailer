package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/Ewan-Greer09/mailer/services/emailer"
	"github.com/Ewan-Greer09/mailer/services/frontend"
	_ "github.com/a-h/templ"
)

type AppConfig struct {
	Address     string `json:"ADDRESS" validate:"required"`
	S3AccessKey string `json:"S3_ACCESS_KEY" validate:"required"`
	S3SecretKey string `json:"S3_SECRET_KEY" validate:"required"`
	S3ViewURL   string `json:"S3_VIEW_URL" validate:"required"`
	S3HostURL   string `json:"S3_HOST_URL" validate:"required"`

	JWTSecretKey string `json:"JWT_SECRET_KEY" validate:"required"`
	MongoURI     string `json:"MONGO_URI" validate:"required"`

	SmtpEmail    string `json:"SMTP_EMAIL" validate:"required"`
	SmtpPassword string `json:"SMTP_PASSWORD" validate:"required"`
}

func main() {
	cfg := loadConfig()

	err := cfg.Validate()
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.HideBanner = true

	store := emailer.NewMongoStore(cfg.MongoURI)
	defer store.Close(context.Background())

	templater := emailer.NewEmailTemplater()

	uploader := emailer.NewS3Uploader(cfg.S3HostURL, cfg.S3ViewURL)

	handler := emailer.NewHandler(emailer.NewEmailService(cfg.SmtpEmail, cfg.SmtpPassword), store, templater, uploader)

	frontendHandler := frontend.New()

	MountRoutes(e, handler, frontendHandler, *cfg)

	err = e.Start(cfg.Address)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("server stopped", "err", err)
}

func MountRoutes(e *echo.Echo, h *emailer.Handler, fh *frontend.Handler, cfg AppConfig) {
	e.Use(
		middleware.RequestID(),
		middleware.Logger(),
		middleware.Recover(),
		middleware.CORS(),
	)

	e.GET("/", fh.HandleRoot)
	e.POST("/recipient", fh.HandleAddRecipient)
	e.GET("/recipient/list", fh.HandleRecipientList)
	e.DELETE("/recipient/:email", fh.HandleDeleteRecipient)

	fh.Static(e)

	api := e.Group("/api") //echojwt.JWT([]byte(cfg.JWTSecretKey))
	api.POST("/send/:communication_type", h.Send)
	api.GET("/:communication_uuid", h.Retrieve)
	api.POST("/send/batch", h.SendBatch)
}

func loadConfig() *AppConfig {
	err := godotenv.Load(".env")
	if err != nil {
		slog.Warn("loadConfig", "err", err)
	}

	return &AppConfig{
		Address:      os.Getenv("ADDRESS"),
		S3AccessKey:  os.Getenv("S3_ACCESS_KEY"),
		S3SecretKey:  os.Getenv("S3_SECRET_KEY"),
		MongoURI:     os.Getenv("MONGO_URI"),
		JWTSecretKey: os.Getenv("JWT_SECRET_KEY"),
		S3ViewURL:    os.Getenv("S3_VIEW_URL"),
		S3HostURL:    os.Getenv("S3_HOST_URL"),
		SmtpEmail:    os.Getenv("SMTP_EMAIL"),
		SmtpPassword: os.Getenv("SMTP_PASSWORD"),
	}
}

func (cfg *AppConfig) Validate() error {
	return validator.New().Struct(cfg)
}
