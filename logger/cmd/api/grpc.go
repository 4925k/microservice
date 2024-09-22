package main

import (
	"context"
	"fmt"
	"github.com/4925k/microservice/logger/data"
	"github.com/4925k/microservice/logger/logs"
	"google.golang.org/grpc"
	"log"
	"net"
)

// LogServer is a log server
type LogServer struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}

// WriteLog will write a log
func (l LogServer) WriteLog(_ context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry()

	// write the log
	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	// insert the log
	err := l.Models.LogEntry.Insert(logEntry)
	if err != nil {
		res := &logs.LogResponse{Result: "failed"}
		return res, err
	}

	// return response
	res := &logs.LogResponse{Result: "success"}

	return res, nil
}

// gRPCListen will listen for gRPC
func (app *Config) gRPCListen() {
	// listen for gRPC
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	// create new gRPC server
	s := grpc.NewServer()

	// register service
	logs.RegisterLogServiceServer(s, &LogServer{Models: app.Models})

	log.Println("gRPC server started on port", grpcPort)

	// start server
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}
}
