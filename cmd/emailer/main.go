package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/Ewan-Greer09/mailer/services/emailer"
)

type AppConfig struct {
	Address      string `json:"ADDRESS"`
	S3AccessKey  string `json:"S3_ACCESS_KEY"`
	S3SecretKey  string `json:"S3_SECRET_KEY"`
	JWTSecretKey string `json:"JWT_SECRET_KEY"`
	MongoURI     string `json:"MONGO_URI"`
}

func main() {
	cfg := loadConfig()

	e := echo.New()
	e.HideBanner = true

	store := emailer.NewMongoStore(cfg.MongoURI)
	defer store.Close(context.Background())

	templater := emailer.NewEmailTemplater()

	uploader := emailer.NewS3Uploader("http://minio:9000", "http://localhost:9000")

	handler := emailer.NewHandler(emailer.NewEmailService(), store, templater, uploader)
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
		log.Fatal(err)
	}

	return &AppConfig{
		Address:      os.Getenv("ADDRESS"),
		S3AccessKey:  os.Getenv("S3_ACCESS_KEY"),
		MongoURI:     os.Getenv("MONGO_URI"),
		JWTSecretKey: os.Getenv("JWT_SECRET_KEY"),
	}
}
