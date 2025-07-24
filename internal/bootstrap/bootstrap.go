package bootstrap

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

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

var (
	app *App
)

type Handler interface {
	Mount(router fiber.Router)
}

type paymentMidtrans struct {
	snap    snap.Client
	coreapi coreapi.Client
}

type App struct {
	http      *fiber.App
	config    *config.Env
	postgres  *sqlx.DB
	validator *validator.Validate
	handlers  []Handler
	jwt       *jwt.JWTStruct
	storage   *supabasestorageuploader.Client
	payment   paymentMidtrans
}

func Initialize() error {
	env, err := config.LoadEnv()
	if err != nil {
		return err
	}

	postgres, err := db.NewPostgres(env)
	if err != nil {
		return err
	}

	jwt := jwt.New(env.JWTSecret)
	val := _validator.New()
	http := http.NewFiber()
	storage := storage.New(env.StorageURL, env.StorageToken, env.StorageBucket)
	pSnap, pCore := payment.New(env.MidtransKey)

	app = &App{
		http:      http,
		config:    env,
		postgres:  postgres,
		validator: val,
		jwt:       jwt,
		storage:   storage,
		payment: paymentMidtrans{
			snap:    pSnap,
			coreapi: pCore,
		},
	}

	logger.New()
	app.InitHandlers()

	go shutdown()

	return app.http.Listen(fmt.Sprintf(":%d", env.AppPort))
}

func shutdown() {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
	<-signalCh

	log.Info().Msg("Received shutdown signal")
	log.Info().Msg("Shutting down...")

	_ = app.postgres.Close()
	_ = app.http.Shutdown()
}
