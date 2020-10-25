package rabbitmq

import (
	"encoding/json"
	libUuid "github.com/google/uuid"
	"log"
)

type Task struct {
	Id libUuid.UUID `json:"id"`
	Description string `json:"description"`
	Tags []string `json:"tags"`
	Status string `json:"status"`
	Progress float32 `json:"progress"`
}

type TaskClient struct {
	amqpClient *AMQPClient
}

func NewTaskManagerClient(client *AMQPClient) *TaskClient {
	taskManagerClient := TaskClient{
		amqpClient: client,
	}

	return &taskManagerClient
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

func (c *TaskClient) PushNewTask(task []byte) error  {
	err := c.amqpClient.Push(task)
	if err != nil {
		return err
	}

	return nil
}
