package rabbitmq

import (
	"encoding/json"
	libUuid "github.com/google/uuid"
	"log"
)

// Task is a task info description.
type Task struct {
	ID          libUuid.UUID `json:"id"`
	Description string       `json:"description"`
	Tags        []string     `json:"tags"`
	Status      string       `json:"status"`
	Progress    float32      `json:"progress"`
}

// TaskClient is a Task manager.
type TaskClient struct {
	amqpClient *AMQPClient
}

// NewTaskManagerClient is used to create the Task manager client.
func NewTaskManagerClient(client *AMQPClient) *TaskClient {
	taskManagerClient := TaskClient{
		amqpClient: client,
	}

	return &taskManagerClient
}

func taskToBytes(task Task) []byte {
	taskJSON, err := json.Marshal(task)
	if err != nil {
		log.Fatal(err)
	}

	return taskJSON
}

// CreateNewTask is used to create a new Task into bytes.
func CreateNewTask(tags []string, description string) []byte {
	newTask := Task{
		ID:          libUuid.New(),
		Tags:        tags,
		Status:      "waiting",
		Progress:    .01,
		Description: description,
	}

	return taskToBytes(newTask)
}

// PushNewTask is used to push a taks into pushQueue.
func (c *TaskClient) PushNewTask(task []byte) error {
	err := c.amqpClient.Push(task)
	if err != nil {
		return err
	}

	return nil
}
