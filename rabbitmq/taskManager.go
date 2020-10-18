package rabbitmq

import (
	"encoding/json"
	libUuid "github.com/google/uuid"
	"github.com/streadway/amqp"
	"log"
)

type Task struct {
	Id libUuid.UUID `json:"id"`
	Tags []string `json:"tags"`
	Status string `json:"status"`
	Progress float32 `json:"progress"`
}

func taskToBytes(task Task) []byte {
	taskJson, err := json.Marshal(task)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(taskJson)

	return taskJson
}

func CreateNewTaskV2(tags []string, status string) []byte {
	newTask := Task{
		Id: libUuid.New(),
		Tags: tags,
		Status: status,
		Progress: .01,
	}

	return taskToBytes(newTask)
}

func CreateNewTask(tags []string, status string) error  {
	newTask := Task{
		Id: libUuid.New(),
		Tags: tags,
		Status: status,
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

	if err != nil {
		log.Fatalf("%s", err)
	}

	//pubAck, pubNack := Channel.NotifyConfirm(make(chan uint64, 1), make(chan uint64, 1))
	//select {
	//case <-pubAck:
	//	log.Println("Ack")
	//case <-pubNack:
	//	log.Println("NAck")
	//}


	return err
}
