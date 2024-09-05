package main

import (
	"context"
	"fmt"
	"github.com/4925k/microservice/logger/data"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://localhost:27017"
	grpcPort = "50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	log.Println("Starting logger service")

	// connect to mongo db
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	client = mongoClient

	// create a context in order to disconnect from mongo
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// close connection
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	// start web server
	go app.server()

	select {}
}

func (app *Config) server() {
	log.Println("Starting service on port", webPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}

// connectToMongo will connect to the mongoDB
func connectToMongo() (*mongo.Client, error) {
	// create connection options
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	// connect
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to MongoDB!")

	return c, nil
}
