package main

import (
	"github.com/casbin/casbin/v2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/swaggo/echo-swagger"
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/config"
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/database"
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/database/migrate"
	tenantModel "gitlab.com/s0j0hn/go-rest-boilerplate-echo/database/models/tenant"
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/docs"
	_ "gitlab.com/s0j0hn/go-rest-boilerplate-echo/docs" // docs is generated by Swag CLI, you have to import it.
	tenantHandler "gitlab.com/s0j0hn/go-rest-boilerplate-echo/handlers"
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/policy"
	"golang.org/x/crypto/acme/autocert"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"os"

	"log"
)

type (
	CustomValidator struct {
		validator *validator.Validate
	}
	PolicyEnforcer struct {
		enforcer *casbin.Enforcer
	}
)

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func (e *PolicyEnforcer) checkPolicyAccessGuests(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := "guest" // All unauthenticated requests only
		method := c.Request().Method
		path := c.Request().URL.Path

		isGood, err := e.enforcer.Enforce(user, path, method)
		if err != nil {
			log.Fatal(err)
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

	echoServer.Validator = &CustomValidator{validator: validator.New()}

	echoServer.Use(middleware.Recover())
	echoServer.Use(middleware.Secure())

	// AMQP.
	//amqpClient.Connect()
	//err := amqpClient.CreateNewTask([]string{"test", "test2"}, "Status is OK")
	//if err != nil {
	//	echoServer.Logger.Fatal(err)
	//	os.Exit(1)
	//}

	// Database.
	gormClient := database.Connect()
	migrate.MigrateDatabase()

	policyEnforcer, err := policy.InitPolicy(gormClient)
	if err != nil {
		echoServer.Logger.Fatal(err)
		os.Exit(1)
	}

	createTenantPolicies(policyEnforcer)

	tenantInstance := tenantModel.TenantModel{}
	tenantHandlerInstance := tenantHandler.CreateHandler(tenantInstance)

	policyCheck := PolicyEnforcer{enforcer: policyEnforcer}

	// Apply the policy for all routes.
	echoServer.Use(policyCheck.checkPolicyAccessGuests)
	echoServer.GET("/tenants/:id", tenantHandlerInstance.GetOneById)
	echoServer.GET("/tenants", tenantHandlerInstance.GetAll)
	echoServer.POST("/tenants", tenantHandlerInstance.Create)
	echoServer.PUT("/tenants", tenantHandlerInstance.Update)
	echoServer.DELETE("/tenants/:id", tenantHandlerInstance.DeleteById)

	docs.SwaggerInfo.Host = config.GetAddress()
	echoServer.GET("/swagger/*", echoSwagger.WrapHandler)

	echoServer.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))


	if config.IsProd() {
		echoServer.AutoTLSManager.Cache = autocert.DirCache("./.cache")
		echoServer.Pre(middleware.HTTPSRedirect())
		go func(c *echo.Echo) {
			echoServer.Logger.Fatal(echoServer.Start(":80"))
		}(echoServer)
		echoServer.Logger.Fatal(echoServer.StartAutoTLS(":443"))
	} else {
		echoServer.Logger.Fatal(echoServer.Start(config.GetAddress()))
	}
}
