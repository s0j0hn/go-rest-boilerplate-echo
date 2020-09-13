package database

import (
	"github.com/streadway/amqp"
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/config"
	"log"
	"sync"
)

var once sync.Once

var Client *amqp.Connection
var Queue amqp.Queue

func Connect() *amqp.Connection {
	once.Do(func() {
		connection := config.GetRabbitMQAccess()
		client, err := amqp.Dial(connection)

		if err != nil {
			log.Fatal("Error AMQP connect:", err)
		}

		defer client.Close()

		Client = client
		CreateDefaultQueue()
	})

	return Client
}

func CreateDefaultQueue() {
	channel, err := Client.Channel()
	if err != nil {
		log.Fatal("Error AMQP channel: ", err)
	}
	defer channel.Close()

	queue, err := channel.QueueDeclare(
		"task_queue", // name
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

	Queue = queue
}
