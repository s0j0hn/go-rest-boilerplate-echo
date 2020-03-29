package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/config"
	tenantModel "gitlab.com/s0j0hn/go-rest-boilerplate-echo/db/models/tenant"
	tenantHandler "gitlab.com/s0j0hn/go-rest-boilerplate-echo/handlers"
	"golang.org/x/crypto/acme/autocert"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
)

type (
	CustomValidator struct {
		validator *validator.Validate
	}
)

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func status(c echo.Context) error {

	message := `
  ____    __
  / __/___/ /  ___
 / _// __/ _ \/ _ \
/___/\__/_//_/\___/ 

`
	return c.String(http.StatusOK, message)
}

func main() {

	echoServer := echo.New()
	tenantInstance := tenantModel.TenantModel{}
	tenantHandlerInstance := tenantHandler.CreateHandler(tenantInstance)

	echoServer.Validator = &CustomValidator{validator: validator.New()}

	echoServer.Use(middleware.Recover())
	echoServer.Use(middleware.Logger())
	echoServer.Use(middleware.Secure())

	echoServer.GET("/", status)

	echoServer.GET("/tenants/:id", tenantHandlerInstance.GetOneById)
	echoServer.POST("/tenants", tenantHandlerInstance.Create)
	echoServer.PUT("/tenants", tenantHandlerInstance.Update)
	echoServer.DELETE("/tenants/:id", tenantHandlerInstance.DeleteById)

	if config.IsProd() {
		echoServer.AutoTLSManager.Cache = autocert.DirCache("./.cache")
		echoServer.Pre(middleware.HTTPSRedirect())
		go func(c *echo.Echo) {
			echoServer.Logger.Fatal(echoServer.Start(":80"))
		}(echoServer)
		echoServer.Logger.Fatal(echoServer.StartAutoTLS(":443"))
	} else {
		echoServer.Logger.Fatal(echoServer.Start(config.GetPort()))
	}
}
