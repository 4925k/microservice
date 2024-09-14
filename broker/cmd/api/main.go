package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"math"
	"net/http"
	"os"
	"time"
)

const (
	webPort     = "80"
	rabbitMQURL = "amqp://guest:guest@rabbitmq"
)

type Config struct {
	Rabbit *amqp.Connection
}

func main() {
	// connect to rabbit mq
	rabbitConn, err := connectToRabbitMQ()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	app := Config{
		Rabbit: rabbitConn,
	}

	log.Printf("Starting broker server on %s\n", webPort)

	// define server
	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	// start server
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

// connectToRabbitMQ helps to connect to rabbit mq
func connectToRabbitMQ() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// don't continue until rabbit is ready
	for {
		c, err := amqp.Dial(rabbitMQURL)
		if err != nil {
			fmt.Println("RabbitMQ not yet ready...")
			counts++
		} else {
			log.Println("connected to RabbitMQ!")
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
