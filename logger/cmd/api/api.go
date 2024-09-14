package main

import (
	"context"
	"github.com/4925k/microservice/logger/data"
	"log"
)

// RPCServer is an RPC server
type RPCServer struct {
}

// RPCPayload is the type for the request
type RPCPayload struct {
	Name string
	Data string
}

// LogInfo will log the info message to the logs collection
func (r *RPCServer) LogInfo(payload RPCPayload, resp *string) error {
	log.Println("LogInfo was invoked with", payload)

	// connect to mongo db and insert data
	collection := client.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), data.LogEntry{
		Name: payload.Name,
		Data: payload.Data,
	})
	if err != nil {
		log.Println("error writing to logs", err)
		return err
	}

	// write response
	*resp = "Processed payload via RPC:" + payload.Name

	return nil
}
