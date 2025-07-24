package bootstrap

import (
	aiHandler "github.com/vistara-studio/vistara-be/internal/domain/ai/handler/rest"
	"github.com/vistara-studio/vistara-be/internal/domain/local/handler/rest"
	localRepository "github.com/vistara-studio/vistara-be/internal/domain/local/repository"
	localService "github.com/vistara-studio/vistara-be/internal/domain/local/service"
	sessionHandler "github.com/vistara-studio/vistara-be/internal/domain/session/handler/rest"
	sessionRepository "github.com/vistara-studio/vistara-be/internal/domain/session/repository"
	sessionService "github.com/vistara-studio/vistara-be/internal/domain/session/service"
	userRepository "github.com/vistara-studio/vistara-be/internal/domain/user/repository"
	"github.com/vistara-studio/vistara-be/pkg/jwt"
	"github.com/gofiber/fiber/v2"
)

// InitHandlers initializes all application handlers and routes
func (app *App) InitHandlers() {
	app.registerRoutes(app.jwt)
	app.MountRoutes()
	app.registerHealthCheck()
}

// registerRoutes initializes repositories, services, and handlers
func (app *App) registerRoutes(jwt *jwt.JWTStruct) {
	// Initialize repositories
	userRepo := userRepository.New(app.postgres)
	sessionRepo := sessionRepository.New(app.postgres)
	localRepo := localRepository.New(app.postgres)

	// Initialize services
	authService := sessionService.New(userRepo, sessionRepo, jwt)
	localBusinessService := localService.New(localRepo, app.payment.snap, app.payment.coreapi)

	// Initialize handlers
	authHandler := sessionHandler.New(authService, app.validator)
	localHandler := rest.New(localBusinessService, app.validator, app.jwt)
	aiHandler := aiHandler.NewAIHandler(app.aiClient, app.validator, app.jwt)

	// Register handlers
	app.handlers = append(app.handlers, authHandler, localHandler, aiHandler)
}

// MountRoutes mounts all registered handlers on the router
func (app *App) MountRoutes() {
	routerGroup := app.http.Group("/api")
	for _, handler := range app.handlers {
		handler.Mount(routerGroup)
	}
}

// registerHealthCheck adds a simple health check endpoint
func (app *App) registerHealthCheck() {
	app.http.Get("/health", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Everything is good!")
	})
}
