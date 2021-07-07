package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	golog "github.com/labstack/gommon/log"
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/config"
	"net/http"
)

type (
	// Handler is a default handler as there is no generics.
	Handler struct {
		amqpMessages chan []byte
	}
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
	return true
}}

func createHandler(amqpMessages chan []byte) *Handler {
	return &Handler{amqpMessages}
}

func (h Handler) getTaskEvents(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	for {
		select {
		case message := <-h.amqpMessages:
			err := ws.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				c.Logger().Error(err)
			}
		}
	}
}

// CreateServer Creates a web socket server
func CreateServer(amqpMessages chan []byte) {
	e := echo.New()

	// For more customizations: https://echo.labstack.com/guide/customization
	if l, ok := e.Logger.(*golog.Logger); ok {
		l.SetHeader("${time_rfc3339} ${level}")
	}

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}  uri=${uri}  status=${status}\n",
	}))

	e.Use(middleware.Recover())

	handler := createHandler(amqpMessages)
	e.GET("/", handler.getTaskEvents)

	e.Logger.Fatal(e.Start(config.GetWebSocketAddress()))
}
