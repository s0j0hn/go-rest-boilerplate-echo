package rabbitmq

import (
	"github.com/streadway/amqp"
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/config"
	"log"
	"sync"
)

var once sync.Once

var Client *amqp.Connection
var Queue amqp.Queue
var Channel *amqp.Channel

func Connect() {
	once.Do(func() {
		connection := config.GetRabbitMQAccess()
		client, err := amqp.Dial(connection)

		if err != nil {
			log.Fatal("Error AMQP connect:", err)
		}

		Client = client
		log.Println("We got a amqp client")
	})
}

func CreateDefaultQueue() {
	channel, err := Client.Channel()
	if err != nil {
		log.Fatal("Error AMQP channel: ", err)
	}

	queue, err := channel.QueueDeclare(
		"tasks_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Fatal("Error declaring AMQP queue: ", err)
	}

	err = channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		log.Fatal("Error declaring AMQP QoS:", err)
	}

	err = channel.Confirm(false)
	if err != nil {
		panic(err)
	}

	Queue = queue
	Channel = channel
}
