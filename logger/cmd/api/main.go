package main

import (
	"context"
	"fmt"
	"github.com/4925k/microservice/logger/data"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"time"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
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

	// start servers
	err = rpc.Register(new(RPCServer))
	go app.rpcListen()
	go app.server()
	go app.gRPCListen()

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

// rpcListen will start the RPC server
func (app *Config) rpcListen() {
	log.Println("Starting RPC server on port", rpcPort)

	// listen on tcp
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		log.Panic(err)
	}
	defer listen.Close()

	// create new server
	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}

		go rpc.ServeConn(rpcConn)
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

	log.Println("Pinging to MongoDB!")

	// Ping the database to check connection
	err = c.Ping(context.TODO(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Println("Pinging to MongoDB!")

	return c, nil
}
