package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"net/http"
)

// Consumer helps to consume events
type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

// NewConsumer helps to create new consumer
func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	c := Consumer{
		conn: conn,
	}

	err := c.setup()
	if err != nil {
		return Consumer{}, err
	}

	return c, nil
}

// setup helps to set up consumer
func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	return declareExchange(channel)
}

// Payload will help to return payload
type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (consumer *Consumer) Listen(topic []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}

	for _, s := range topic {
		err = ch.QueueBind(
			q.Name,
			s,
			"logs_topic",
			false,
			nil,
		)
		if err != nil {
			return err
		}

	}

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	forever := make(chan bool)

	go func() {
		for d := range messages {
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)

			go handlePayload(payload)
		}
	}()

	fmt.Printf("waiting for message [Exchange, Queue] [logs_topic, %s]\n", q.Name)

	<-forever

	return nil
}

// handlePayload helps to handle payload
func handlePayload(payload Payload) {
	switch payload.Name {
	case "log", "event":
		// log
		err := logEvent(payload)
		if err != nil {
			fmt.Println("failed to log event", err)
		}
	case "auth":
		// authenticate
	default:
		err := logEvent(payload)
		if err != nil {
			fmt.Println("failed to log event", err)
		}
	}
}

// logEvent helps to log event
func logEvent(entry Payload) error {
	// get data
	jsonData, _ := json.Marshal(entry)

	// call the service
	req, err := http.NewRequest("POST", "http://logger/log", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	// get response from auth service
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err

	}
	defer res.Body.Close()

	// read status
	if res.StatusCode != http.StatusAccepted {
		return err

	}

	return nil
}
