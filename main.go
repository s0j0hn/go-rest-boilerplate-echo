package main

import (
	"context"
	"crypto/tls"
	"errors"
	"github.com/casbin/casbin/v2"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	logger "github.com/labstack/gommon/log"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/swaggo/echo-swagger"
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/config"
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/database"
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/database/migrate"
	tenantModel "gitlab.com/s0j0hn/go-rest-boilerplate-echo/database/models/tenant"
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/docs"
	_ "gitlab.com/s0j0hn/go-rest-boilerplate-echo/docs" // docs are generated by Swag CLI, you have to import it.
	tenantHandler "gitlab.com/s0j0hn/go-rest-boilerplate-echo/handlers"
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/policy"
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/rabbitmq"
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/websocket"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type (
	// CustomValidator of requests.
	CustomValidator struct {
		validator *validator.Validate
	}
	// PolicyEnforcer is casbin rules policy.
	PolicyEnforcer struct {
		enforcer *casbin.Enforcer
	}
)

// Validate is just a init
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func (e *PolicyEnforcer) checkPolicyAccessGuests(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := "guest" // All unauthenticated requests only
		method := c.Request().Method
		path := c.Request().URL.Path

		isGood, err := e.enforcer.Enforce(user, path, method)
		if err != nil {
			return echo.ErrForbidden
		}

		if isGood {
			return next(c)
		}
		return echo.ErrForbidden
	}
}

func createTenantPolicies(policyEnforcer *casbin.Enforcer) {
	policy.AddGetPolicy(policyEnforcer, "guest", "/swagger/*")
	policy.AddGetPolicy(policyEnforcer, "guest", "/tenants")
	policy.AddCreatePolicy(policyEnforcer, "guest", "/tenants")
	policy.AddUpdatePolicy(policyEnforcer, "guest", "/tenants")
	// policy.AddDeletePolicy(policyEnforcer, "guest", "/tenants")
}

// @title Swagger Boilerplate API
// @version 1.0
// @description This is a sample
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /
func main() {
	echoServer := echo.New()
	echoServer.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}  uri=${uri}  status=${status}\n",
	}))

	// For more customizations: https://echo.labstack.com/guide/customization
	if l, ok := echoServer.Logger.(*logger.Logger); ok {
		l.SetHeader("${time_rfc3339} ${level}")
	}

	echoServer.Validator = &CustomValidator{validator: validator.New()}

	// API cache
	// s := souin_echo.New(souin_echo.DevDefaultConfiguration)

	echoServer.Use(middleware.Recover())
	echoServer.Use(middleware.Secure())
	// echoServer.Use(s.Process)

	// Database.
	gormClient := database.Connect()
	migrate.RunMigrateDatabase()

	policyEnforcer, err := policy.InitPolicy(gormClient)
	if err != nil {
		echoServer.Logger.Fatal(err)
		os.Exit(1)
	}

	createTenantPolicies(policyEnforcer)

	zeroLogger := log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	doneChannel := make(chan bool)
	messagesChannel := make(chan []byte)
	amqpContext := context.Background()
	rabbitMQClient := rabbitmq.NewAMQPClient(config.GetAMQPQListenQueue(), config.GetAMQPPushQueue(), config.GetRabbitMQAccess(), zeroLogger, doneChannel, messagesChannel, true)
	doneChannel <- true
	taskManager := rabbitmq.NewTaskManagerClient(rabbitMQClient)

	go func() {
		for {
			err = rabbitMQClient.Stream(amqpContext)
			if errors.Is(err, rabbitmq.ErrDisconnected) {
				continue
			}
			break
		}
	}()

	tenantInstance := tenantModel.ModelTenant{}
	tenantHandlerInstance := tenantHandler.CreateHandlerTenant(tenantInstance, taskManager)

	policyCheck := PolicyEnforcer{enforcer: policyEnforcer}

	// Apply the policy for all routes.
	echoServer.Use(policyCheck.checkPolicyAccessGuests)
	echoServer.GET("/tenants/:id", tenantHandlerInstance.GetOneByID)
	echoServer.GET("/tenants", tenantHandlerInstance.GetAll)
	echoServer.POST("/tenants", tenantHandlerInstance.Create)
	echoServer.PUT("/tenants", tenantHandlerInstance.Update)
	echoServer.DELETE("/tenants/:id", tenantHandlerInstance.DeleteByID)

	docs.SwaggerInfo.Host = config.GetAddress()
	echoServer.GET("/swagger/*", echoSwagger.WrapHandler)

	echoServer.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	rateLimiterConfig := middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{Rate: 10, Burst: 30, ExpiresIn: 3 * time.Minute},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			id := ctx.RealIP()
			return id, nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(http.StatusForbidden, nil)
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return context.JSON(http.StatusTooManyRequests, nil)
		},
	}

	echoServer.Use(middleware.RateLimiterWithConfig(rateLimiterConfig))

	go websocket.CreateServer(messagesChannel)

	if config.IsProd() {
		autoTLSManager := autocert.Manager{
			Prompt: autocert.AcceptTOS,
			// Cache certificates to avoid issues with rate limits (https://letsencrypt.org/docs/rate-limits)
			Cache: autocert.DirCache("./.cache"),
			//HostPolicy: autocert.HostWhitelist("<DOMAIN>"),
		}

		echoServer.AutoTLSManager.Cache = autocert.DirCache("./.cache")
		echoServer.Pre(middleware.HTTPSRedirect())

		s := http.Server{
			Addr:    ":443",
			Handler: echoServer, // set Echo as handler
			TLSConfig: &tls.Config{
				//Certificates: nil, // <-- s.ListenAndServeTLS will populate this field
				GetCertificate:   autoTLSManager.GetCertificate,
				NextProtos:       []string{acme.ALPNProto},
				MinVersion:       tls.VersionTLS13,
				CurvePreferences: []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
				CipherSuites: []uint16{
					tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
					tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_RSA_WITH_AES_256_CBC_SHA,
				},
			},
			//ReadTimeout: 30 * time.Second, // use custom timeouts
		}
		if err := s.ListenAndServeTLS("", ""); err != http.ErrServerClosed {
			echoServer.Logger.Fatal(err)
		}
	} else {
		go func(c *echo.Echo) {
			if err := echoServer.Start(config.GetAddress()); err != nil && err != http.ErrServerClosed {
				echoServer.Logger.Fatal("shutting down the server")
			}
		}(echoServer)
	}

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := echoServer.Shutdown(ctx); err != nil {
		echoServer.Logger.Fatal(err)
	}
}
