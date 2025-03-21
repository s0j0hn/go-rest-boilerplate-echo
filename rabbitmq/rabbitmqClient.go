package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NeowayLabs/wabbit"
	"github.com/NeowayLabs/wabbit/amqptest"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"runtime"
	"sync"
	"time"
)

var (
	// ErrDisconnected the message error for disconnection
	ErrDisconnected = errors.New("disconnected from rabbitmq, trying to reconnect")
)

const (
	// When reconnecting to the server after connection failure
	reconnectDelay = 5 * time.Second
)

// AMQPClient holds necessery information for rabbitMQ
type AMQPClient struct {
	isReal             bool
	pushQueue          string
	listenQueue        string
	logger             zerolog.Logger
	connection         *amqp.Connection
	falseConnection    *amqptest.Conn
	falseChannel       wabbit.Channel
	amqpChannel        *amqp.Channel
	doneChannel        chan bool
	notifyClose        chan *amqp.Error
	notifyConfirm      chan amqp.Confirmation
	falseNotifyClose   chan wabbit.Error
	falseNotifyConfirm chan wabbit.Confirmation
	isConnected        bool
	alive              bool
	threads            int
	wg                 *sync.WaitGroup
	activeConsumers    []string
	messagesChannel    chan []byte
}

// NewAMQPClient is a constructor that takes address, push and listen queue names, logger, and a amqpChannel that will notify rabbitmq client on server shutdown. We calculate the number of threads, create the client, and start the connection process. Connect method connects to the rabbitmq server and creates push/listen channels if they don't exist.
func NewAMQPClient(listenQueue, pushQueue, addr string, l zerolog.Logger, done chan bool, messages chan []byte, isReal bool) *AMQPClient {
	threads := runtime.GOMAXPROCS(0)
	if numCPU := runtime.NumCPU(); numCPU > threads {
		threads = numCPU
	}

	client := AMQPClient{
		listenQueue:     listenQueue,
		logger:          l,
		threads:         threads,
		doneChannel:     done,
		pushQueue:       pushQueue,
		alive:           true,
		wg:              &sync.WaitGroup{},
		messagesChannel: messages,
		isReal:          isReal,
	}

	if isReal {
		client.wg.Add(threads)
	} else {
		client.wg.Add(1)
	}

	go client.handleReconnect(addr)
	return &client
}

// handleReconnect will wait for a connection error on
// notifyClose, and then continuously attempt to reconnect.
func (c *AMQPClient) handleReconnect(addr string) {
	if c.isReal {
		for c.alive {
			var retryCount int
			c.logger.Printf("Attempting to connect to rabbitMQ: %s", addr)

			c.isConnected = false
			t := time.Now()

			for !c.connect(addr) {
				if !c.alive {
					return
				}

				select {
				case <-c.doneChannel:
					c.logger.Printf("Received something into done amqpChannel")
					return
				case <-time.After(reconnectDelay + time.Duration(retryCount)*time.Second):
					c.logger.Printf("disconnected from rabbitMQ and failed to connect")
					retryCount++
				}
			}

			c.logger.Printf("Connected to rabbitMQ in: %vms", time.Since(t).Milliseconds())
			select {
			case <-c.doneChannel:
				return
			case <-c.notifyClose:
			}
		}
	} else {
		c.logger.Printf("Creating a fake client to rabbitMQ: %s", addr)
		for c.alive {
			c.logger.Printf("Attempting to connect to false rabbitMQ: %s", addr)

			c.isConnected = false
			t := time.Now()

			c.connect(addr)

			c.logger.Printf("Connected to rabbitMQ in: %vms", time.Since(t).Milliseconds())
			select {
			case <-c.doneChannel:
				return
			case <-c.notifyClose:
			}
		}
	}
}

// connect will make a single attempt to connect to
// RabbitMq. It returns the success of the attempt.
func (c *AMQPClient) connect(addr string) bool {
	if c.isReal {
		conn, err := amqp.Dial(addr)
		if err != nil {
			c.logger.Printf("failed to dial rabbitMQ server: %v", err)
			return false
		}

		ch, err := conn.Channel()
		if err != nil {
			c.logger.Printf("failed connecting to amqpChannel: %v", err)
			return false
		}

		err = ch.Confirm(false)
		if err != nil {
			c.logger.Printf("failed to confirm amqpChannel: %v", err)
			return false
		}

		_, err = ch.QueueDeclare(
			c.listenQueue,
			true,  // Durable
			false, // Delete when unused
			false, // Exclusive
			false, // No-wait
			nil,   // Arguments
		)
		if err != nil {
			c.logger.Printf("failed to declare listen queue: %v", err)
			return false
		}

		_, err = ch.QueueDeclare(
			c.pushQueue,
			true,  // Durable
			false, // Delete when unused
			false, // Exclusive
			false, // No-wait
			nil,   // Arguments
		)

		if err != nil {
			c.logger.Printf("failed to declare push queue: %v", err)
			return false
		}

		c.changeRealConnection(conn, ch)
		c.isConnected = true

		return true
	}

	conn, err := amqptest.Dial(addr)

	ch, err := conn.Channel()
	if err != nil {
		c.logger.Printf("failed connecting to amqpChannel: %v", err)
		return false
	}

	err = ch.Confirm(false)
	if err != nil {
		c.logger.Printf("failed to confirm amqpChannel: %v", err)
		return false
	}

	options := wabbit.Option{
		"durable":    true,
		"autoDelete": false,
		"exclusive":  false,
		"noWait":     false,
	}

	_, err = ch.QueueDeclare(c.listenQueue, options)
	if err != nil {
		c.logger.Printf("failed to declare listen queue: %v", err)
		return false
	}

	_, err = ch.QueueDeclare(c.pushQueue, options)

	if err != nil {
		c.logger.Printf("failed to declare push queue: %v", err)
		return false
	}

	c.logger.Printf("Connected")

	c.changeConnection(conn, ch)
	c.isConnected = true

	return true
}

// changeRealConnection takes a new connection to the queue,
// and updates the amqpChannel listeners to reflect this.
func (c *AMQPClient) changeRealConnection(connection *amqp.Connection, channel *amqp.Channel) {
	c.connection = connection
	c.amqpChannel = channel
	c.notifyClose = make(chan *amqp.Error)
	c.notifyConfirm = make(chan amqp.Confirmation)

	c.amqpChannel.NotifyClose(c.notifyClose)
	c.amqpChannel.NotifyPublish(c.notifyConfirm)
}

// changeConnection takes a new connection to the queue, used for tests,
// and updates the amqpChannel listeners to reflect this.
func (c *AMQPClient) changeConnection(connection *amqptest.Conn, channel wabbit.Channel) {
	c.falseConnection = connection
	c.falseChannel = channel
	c.falseNotifyClose = make(chan wabbit.Error)
	c.falseNotifyConfirm = make(chan wabbit.Confirmation)

	c.falseChannel.NotifyClose(c.falseNotifyClose)
	c.falseChannel.NotifyPublish(c.falseNotifyConfirm)
}

// Push will push data onto the queue, and wait for a confirmation.
// If no confirms are received until within the resendTimeout,
// it continuously resends messages until a confirmation is received.
// This will block until the server sends a confirm.
func (c *AMQPClient) Push(data []byte) error {
	if !c.isConnected {
		return ErrDisconnected
	}

	for {
		err := c.UnsafePush(data)

		if err != nil {
			if err == ErrDisconnected {
				continue
			}
			return err
		}

		if c.isReal {
			select {
			case confirm := <-c.notifyConfirm:
				if confirm.Ack {
					return nil
				}
			case <-time.After(1 * time.Second):
			}
		} else {
			return nil
		}
	}
}

// UnsafePush will push to the queue without checking for
// confirmation. It returns an error if it fails to connect.
// No guarantees are provided for whether the server will
// receive the message.
func (c *AMQPClient) UnsafePush(data []byte) error {
	if !c.isConnected {
		return ErrDisconnected
	}

	if c.isReal {
		return c.amqpChannel.Publish(
			"",          // Exchange
			c.pushQueue, // Routing key
			false,       // Mandatory
			false,       // Immediate
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "text/plain",
				Body:         data,
			},
		)
	}

	c.logger.Printf("Pushing message for tests")
	return c.falseChannel.Publish(
		"",
		c.pushQueue,
		data,
		wabbit.Option{
			"deliveryMode": 2,
			"contentType":  "text/plain",
		},
	)
}

// Stream is used to listen on queue and parse the messages.
func (c *AMQPClient) Stream(cancelCtx context.Context) error {
	for {
		if c.isConnected {
			break
		}
		time.Sleep(1 * time.Second)
	}

	err := c.amqpChannel.Qos(1, 0, false)
	if err != nil {
		return err
	}

	var connectionDropped bool

	c.logger.Printf("Starting to wait for rabbitmq events ...")
	for i := 1; i <= c.threads; i++ {
		messages, err := c.amqpChannel.Consume(
			c.listenQueue,
			consumerName(i), // Consumer
			false,           // Auto-Ack
			false,           // Exclusive
			false,           // No-local
			false,           // No-Wait
			nil,             // Args
		)
		if err != nil {
			return err
		}

		c.activeConsumers = append(c.activeConsumers, consumerName(i))

		go func() {
			defer c.wg.Done()
			for {
				select {
				case <-cancelCtx.Done():
					return
				case message, ok := <-messages:
					if !ok {
						connectionDropped = true
						return
					}
					c.parseEvent(message)
				}
			}
		}()

	}

	c.wg.Wait()

	if connectionDropped {
		return ErrDisconnected
	}

	return nil
}

func (c *AMQPClient) parseEvent(msg amqp.Delivery) {
	var evt Task

	l := c.logger.Log().Timestamp()
	startTime := time.Now()

	go func() {
		c.messagesChannel <- msg.Body
	}()

	err := json.Unmarshal(msg.Body, &evt)
	if err != nil {
		logAndNack(msg, l, startTime, "unmarshalling body: %s - %s", string(msg.Body), err.Error())
		return
	}

	if evt.Status == "" {
		logAndNack(msg, l, startTime, "received event without data")
		return
	}

	// TODO: implement a away to store the tasks
	switch evt.Status {
	case "waiting":
	case "failed":
	case "running":
	case "completed":
	default:
		err = msg.Reject(false)
		if err != nil {
			logAndNack(msg, l, startTime, err.Error())
			return
		}
		return
	}

	l.Str("level", "info").Int64("took-ms", time.Since(startTime).Milliseconds()).Msgf("%s parsed successfully", evt.Description)

	err = msg.Ack(false)
	if err != nil {
		logAndNack(msg, l, startTime, err.Error())
		return
	}
}

func logAndNack(msg amqp.Delivery, l *zerolog.Event, t time.Time, errorMessage string, args ...interface{}) {
	err := msg.Nack(false, false)
	if err != nil {
		panic(err)
		return
	}
	l.Int64("took-ms", time.Since(t).Milliseconds()).Str("level", "error").Msg(fmt.Sprintf(errorMessage, args...))
}

// Close is used to destroy all tcp connection to rabbitmq.
func (c *AMQPClient) Close() error {
	if !c.isConnected {
		return nil
	}

	c.alive = false
	c.logger.Printf("Waiting for current messages to be processed...")

	go func() {
		if c.threads == 0 {
			defer c.wg.Done()
		}
		for i := 1; i <= len(c.activeConsumers); i++ {
			err := c.amqpChannel.Cancel(consumerName(i), false)
			if err != nil {
				c.logger.Printf("error canceling consumer %s: %v", consumerName(i), err)
			}
		}
	}()

	c.activeConsumers = nil

	if c.isReal {
		err := c.amqpChannel.Close()
		if err != nil {
			return err
		}

		err = c.connection.Close()
		if err != nil {
			return err
		}
	} else {
		err := c.falseChannel.Close()
		if err != nil {
			return err
		}

		err = c.falseConnection.Close()
		if err != nil {
			return err
		}
	}

	c.isConnected = false
	c.logger.Printf("gracefully stopped rabbitMQ connection")
	return nil
}

func consumerName(i int) string {
	return fmt.Sprintf("go-consumer-%v", i)
}
