package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/swaggo/echo-swagger"
	"gorm.io/gorm"

	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/config"
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/database"
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/database/migrate"
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/docs"
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/internal/domain"
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/internal/handlers"
	customMiddleware "gitlab.com/s0j0hn/go-rest-boilerplate-echo/internal/middleware"
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/internal/repository"
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/internal/service"
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/policy"
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/rabbitmq"
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/websocket"
)

// App represents the main application
type App struct {
	config      *config.Config
	server      *echo.Echo
	db          *gorm.DB
	rabbitMQ    *rabbitmq.AMQPClient
	taskManager *rabbitmq.TaskClient
	logger      zerolog.Logger
}

// NewApp creates a new application instance
func NewApp(cfg *config.Config) *App {
	return &App{
		config: cfg,
		logger: log.Output(zerolog.ConsoleWriter{Out: os.Stderr}),
	}
}

// Initialize sets up all application dependencies
func (a *App) Initialize() error {
	// Setup database
	if err := a.setupDatabase(); err != nil {
		return fmt.Errorf("failed to setup database: %w", err)
	}

	// Setup RabbitMQ
	if err := a.setupRabbitMQ(); err != nil {
		return fmt.Errorf("failed to setup RabbitMQ: %w", err)
	}

	// Setup HTTP server
	if err := a.setupServer(); err != nil {
		return fmt.Errorf("failed to setup server: %w", err)
	}

	return nil
}

// Start starts the application
func (a *App) Start(ctx context.Context) error {
	// Start WebSocket server
	messagesChannel := make(chan []byte)
	go websocket.CreateServer(messagesChannel)

	// Start RabbitMQ streaming
	go func() {
		for {
			err := a.rabbitMQ.Stream(ctx)
			if errors.Is(err, rabbitmq.ErrDisconnected) {
				continue
			}
			break
		}
	}()

	// Start HTTP server
	go func() {
		if err := a.server.Start(a.config.GetAddress()); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.logger.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// Graceful shutdown
	return a.Shutdown(ctx)
}

// Shutdown gracefully shuts down the application
func (a *App) Shutdown(ctx context.Context) error {
	a.logger.Info().Msg("Shutting down application...")

	// Create shutdown context with timeout
	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := a.server.Shutdown(shutdownCtx); err != nil {
		a.logger.Error().Err(err).Msg("Failed to shutdown server gracefully")
		return err
	}

	// Close database connection
	if a.db != nil {
		sqlDB, err := a.db.DB()
		if err == nil {
			sqlDB.Close()
		}
	}

	// Close RabbitMQ connection
	if a.rabbitMQ != nil {
		a.rabbitMQ.Close()
	}

	a.logger.Info().Msg("Application shutdown complete")
	return nil
}

// setupDatabase initializes the database connection and runs migrations
func (a *App) setupDatabase() error {
	db := database.Connect()
	if err := migrate.RunMigrateDatabase(); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	// Migrate new domain models
	if err := db.AutoMigrate(&domain.Tenant{}); err != nil {
		return fmt.Errorf("failed to migrate domain models: %w", err)
	}

	a.db = db
	return nil
}

// setupRabbitMQ initializes RabbitMQ connection and task manager
func (a *App) setupRabbitMQ() error {
	doneChannel := make(chan bool)
	messagesChannel := make(chan []byte)

	rabbitMQClient := rabbitmq.NewAMQPClient(
		a.config.RabbitMQ.ListenQueue,
		a.config.RabbitMQ.PushQueue,
		a.config.GetRabbitMQDSN(),
		a.logger,
		doneChannel,
		messagesChannel,
		true,
	)

	doneChannel <- true
	taskManager := rabbitmq.NewTaskManagerClient(rabbitMQClient)

	a.rabbitMQ = rabbitMQClient
	a.taskManager = taskManager
	return nil
}

// setupServer initializes the HTTP server with all middleware and routes
func (a *App) setupServer() error {
	e := echo.New()

	// Setup validator
	e.Validator = &CustomValidator{validator: validator.New()}

	// Setup middleware
	a.setupMiddleware(e)

	// Setup authorization
	policyEnforcer, err := a.setupAuthorization()
	if err != nil {
		return fmt.Errorf("failed to setup authorization: %w", err)
	}

	// Setup routes
	a.setupRoutes(e, policyEnforcer)

	a.server = e
	return nil
}

// setupMiddleware configures all middleware
func (a *App) setupMiddleware(e *echo.Echo) {
	// Basic middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method} uri=${uri} status=${status} latency=${latency_human}\n",
	}))
	e.Use(middleware.Recover())
	e.Use(customMiddleware.ErrorHandler())

	// Security middleware
	e.Use(middleware.Secure())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"https://localhost:3000", "https://yourdomain.com"},
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPost, http.MethodDelete},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	// Timeout middleware
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: 30 * time.Second,
	}))

	// Rate limiting
	e.Use(middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{
				Rate:      10,
				Burst:     30,
				ExpiresIn: 3 * time.Minute,
			},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			return ctx.RealIP(), nil
		},
		ErrorHandler: func(ctx echo.Context, err error) error {
			return ctx.JSON(http.StatusTooManyRequests, map[string]string{
				"error": "Rate limit exceeded",
			})
		},
	}))
}

// setupAuthorization initializes Casbin policy enforcement
func (a *App) setupAuthorization() (*casbin.Enforcer, error) {
	policyEnforcer, err := policy.InitPolicy(a.db)
	if err != nil {
		return nil, err
	}

	// Create default policies
	a.createDefaultPolicies(policyEnforcer)

	return policyEnforcer, nil
}

// setupRoutes configures all application routes
func (a *App) setupRoutes(e *echo.Echo, policyEnforcer *casbin.Enforcer) {
	// Initialize repositories
	tenantRepo := repository.NewTenantRepository(a.db)

	// Initialize services
	tenantService := service.NewTenantService(tenantRepo, a.taskManager)

	// Initialize handlers
	tenantHandler := handlers.NewTenantHandler(tenantService)
	healthHandler := handlers.NewHealthHandler(a.db, a.rabbitMQ)

	// Health check routes (no auth required)
	e.GET("/health", healthHandler.Health)
	e.GET("/ready", healthHandler.Ready)
	e.GET("/live", healthHandler.Live)

	// API routes group with authentication
	api := e.Group("/api/v1")

	// Apply policy middleware
	policyCheck := &PolicyEnforcer{enforcer: policyEnforcer}
	api.Use(policyCheck.checkPolicyAccessGuests)

	// Tenant routes
	api.GET("/tenants", tenantHandler.GetAllTenants)
	api.GET("/tenants/:id", tenantHandler.GetTenantByID)
	api.POST("/tenants", tenantHandler.CreateTenant)
	api.PUT("/tenants", tenantHandler.UpdateTenant)
	api.DELETE("/tenants/:id", tenantHandler.DeleteTenant)

	// Swagger documentation
	docs.SwaggerInfo.Host = a.config.GetAddress()
	e.GET("/swagger/*", echoSwagger.WrapHandler)
}

// createDefaultPolicies creates the default authorization policies
func (a *App) createDefaultPolicies(policyEnforcer *casbin.Enforcer) {
	// Allow access to health endpoints
	policy.AddGetPolicy(policyEnforcer, "guest", "/health")
	policy.AddGetPolicy(policyEnforcer, "guest", "/ready")
	policy.AddGetPolicy(policyEnforcer, "guest", "/live")

	// Allow access to Swagger
	policy.AddGetPolicy(policyEnforcer, "guest", "/swagger/*")

	// Tenant endpoints
	policy.AddGetPolicy(policyEnforcer, "guest", "/api/v1/tenants")
	policy.AddCreatePolicy(policyEnforcer, "guest", "/api/v1/tenants")
	policy.AddUpdatePolicy(policyEnforcer, "guest", "/api/v1/tenants")
	// policy.AddDeletePolicy(policyEnforcer, "guest", "/api/v1/tenants") // Uncomment to allow delete
}

// CustomValidator wraps the validator
type CustomValidator struct {
	validator *validator.Validate
}

// Validate validates the given struct
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// PolicyEnforcer wraps casbin enforcer
type PolicyEnforcer struct {
	enforcer *casbin.Enforcer
}

// checkPolicyAccessGuests middleware for checking guest access
func (p *PolicyEnforcer) checkPolicyAccessGuests(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := "guest"
		method := c.Request().Method
		path := c.Request().URL.Path

		allowed, err := p.enforcer.Enforce(user, path, method)
		if err != nil {
			return echo.ErrInternalServerError
		}

		if !allowed {
			return echo.ErrForbidden
		}

		return next(c)
	}
}