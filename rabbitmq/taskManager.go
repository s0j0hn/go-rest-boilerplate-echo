package rabbitmq

import (
	"encoding/json"
	libUuid "github.com/google/uuid"
	"log"
)

type task struct {
	ID          libUuid.UUID `json:"id"`
	Description string       `json:"description"`
	Tags        []string     `json:"tags"`
	Status      string       `json:"status"`
	Progress    float32      `json:"progress"`
}

// TaskClient is a task manager.
type TaskClient struct {
	amqpClient *AMQPClient
}

// NewTaskManagerClient is used to create the task manager client.
func NewTaskManagerClient(client *AMQPClient) *TaskClient {
	taskManagerClient := TaskClient{
		amqpClient: client,
	}

	return &taskManagerClient
}

func taskToBytes(task task) []byte {
	taskJSON, err := json.Marshal(task)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(taskJSON)

	return taskJSON
}

// CreateNewTask is used to create a new task into bytes.
func CreateNewTask(tags []string, status string) []byte {
	newTask := task{
		ID:       libUuid.New(),
		Tags:     tags,
		Status:   status,
		Progress: .01,
	}

	return taskToBytes(newTask)
}

// PushNewTask is used to push a taks into pushQueue.
func (c *TaskClient) PushNewTask(task []byte) error  {
	err := c.amqpClient.Push(task)
	if err != nil {
		return err
	}

	return nil
}
