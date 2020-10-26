package rabbitmq

import (
	"github.com/rs/zerolog/log"
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/config"
	"os"
	"testing"
)

var amqpClient *AMQPClient

func TestMain(m *testing.M) {
	goChan := make(chan os.Signal, 1)
	amqpClient = NewAMQPClient("listenqueue", "pushqueue", config.GetRabbitMQAccess(), log.Logger, goChan)

	os.Exit(m.Run())
}
