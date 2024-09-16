package main

import (
	"context"
	"github.com/4925k/microservice/logger/data"
	"github.com/4925k/microservice/logger/logs"
)

// LogServer is a log server
type LogServer struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}

// WriteLog will write a log
func (l LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
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
