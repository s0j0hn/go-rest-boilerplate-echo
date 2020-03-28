package main

import (
	userHandler "github.com/tkc/go-echo-server-sandbox/handlers"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/tkc/go-echo-server-sandbox/config"
	userModel "github.com/tkc/go-echo-server-sandbox/models/tenant"
	"golang.org/x/crypto/acme/autocert"
	"gopkg.in/go-playground/validator.v9"
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
	userModelInstance := userModel.TenantModel{}
	userHandlerInstance := userHandler.CreateHandler(userModelInstance)

	echoServer.Validator = &CustomValidator{validator: validator.New()}

	echoServer.Use(middleware.Recover())
	echoServer.Use(middleware.Logger())

	echoServer.GET("/", status)

	echoServer.GET("/tenant/:id", userHandlerInstance.GetOneById)
	echoServer.POST("/tenant", userHandlerInstance.Create)
	echoServer.PUT("/tenant", userHandlerInstance.Update)
	echoServer.DELETE("/tenant", userHandlerInstance.DeleteById)

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
