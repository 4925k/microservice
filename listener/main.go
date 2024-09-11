package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"math"
	"os"
	"time"
)

func main() {
	// TODO connect to rabbit mq
	rabbitConn, err := connectToRabbitMQ()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()
	log.Println("connected to RabbitMQ!")

	// TODO start listening to messages

	// TODO create consumer

	// TODO watch the queue and consume events

}

// connectToRabbitMQ helps to connect to rabbit mq
func connectToRabbitMQ() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// don't continue until rabbit is ready
	for {
		c, err := amqp.Dial("amqp://guest:guest@localhost")
		if err != nil {
			fmt.Println("RabbitMQ not yet ready...")
			counts++
		} else {
			connection = c
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		time.Sleep(backOff)
		backOff = time.Duration(math.Pow(float64(counts), 2))
		log.Print("backing off...")

	}

	return connection, nil
}
