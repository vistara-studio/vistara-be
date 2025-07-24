package bootstrap

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/vistara-studio/vistara-be/internal/infra/ai"
	"github.com/vistara-studio/vistara-be/internal/infra/config"
	"github.com/vistara-studio/vistara-be/internal/infra/db"
	"github.com/vistara-studio/vistara-be/internal/infra/http"
	"github.com/vistara-studio/vistara-be/internal/infra/logger"
	"github.com/vistara-studio/vistara-be/internal/infra/payment"
	"github.com/vistara-studio/vistara-be/internal/infra/storage"
	"github.com/vistara-studio/vistara-be/pkg/jwt"
	_validator "github.com/vistara-studio/vistara-be/pkg/validator"
	
	supabasestorageuploader "github.com/adityarizkyramadhan/supabase-storage-uploader"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/rs/zerolog/log"
)

var app *App

// Handler interface for mounting routes
type Handler interface {
	Mount(router fiber.Router)
}

// paymentMidtrans wraps Midtrans payment clients
type paymentMidtrans struct {
	snap    snap.Client
	coreapi coreapi.Client
}

// App holds all application dependencies
type App struct {
	http      *fiber.App
	config    *config.Env
	postgres  *sqlx.DB
	validator *validator.Validate
	handlers  []Handler
	jwt       *jwt.JWTStruct
	storage   *supabasestorageuploader.Client
	payment   paymentMidtrans
	aiClient  *ai.Client
}

// Initialize starts the application with all dependencies
func Initialize() error {
	// Load configuration
	env, err := config.LoadEnv()
	if err != nil {
		return err
	}

	// Initialize database
	postgres, err := db.NewPostgres(env)
	if err != nil {
		return err
	}

	// Initialize dependencies
	jwt := jwt.New(env.JWTSecret)
	validator := _validator.New()
	httpServer := http.NewFiber()
	storage := storage.New(env.StorageURL, env.StorageToken, env.StorageBucket)
	paymentSnap, paymentCore := payment.New(env.MidtransKey)
	aiClient := ai.NewClient(env.VistaraAIURL, env.VistaraAIKey)

	// Create app instance
	app = &App{
		http:      httpServer,
		config:    env,
		postgres:  postgres,
		validator: validator,
		jwt:       jwt,
		storage:   storage,
		payment: paymentMidtrans{
			snap:    paymentSnap,
			coreapi: paymentCore,
		},
		aiClient: aiClient,
	}

	// Initialize logger and handlers
	logger.New()
	app.InitHandlers()

	// Start graceful shutdown listener
	go shutdown()

	// Start HTTP server
	return app.http.Listen(fmt.Sprintf(":%d", env.AppPort))
}

// shutdown handles graceful application shutdown
func shutdown() {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
	<-signalCh

	log.Info().Msg("Received shutdown signal")
	log.Info().Msg("Shutting down gracefully...")

	_ = app.postgres.Close()
	_ = app.http.Shutdown()
}
