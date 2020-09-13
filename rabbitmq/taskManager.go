package database

import (
	"encoding/json"
	libUuid "github.com/google/uuid"
	"github.com/streadway/amqp"
	"log"
)

type Task struct {
	Id libUuid.UUID `json:"id"`
	Tags []string `json:"tags"`
	Status []string `json:"status"`
	Progress float32 `json:"progress"`
}

func taskToBytes(task Task) []byte {
	taskJson, err := json.Marshal(task)
	if err != nil {
		log.Fatal(err)
	}

	return taskJson
}

func CreateNewTask(tags []string, status string) error  {
	newTask := Task{
		Id: libUuid.New(),
		Tags: tags,
		Status: []string{status},
		Progress: .01,
	}

	err := Channel.Publish(
		"",     // exchange
		Queue.Name, // routing key
		false,  // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         taskToBytes(newTask),
		})

	pubAck, pubNack := Channel.NotifyConfirm(make(chan uint64, 1), make(chan uint64, 1))
	select {
	case <-pubAck:
		// fmt.Println("Ack")
	case <-pubNack:
		// fmt.Println("NAck")
	}


	return err
}
