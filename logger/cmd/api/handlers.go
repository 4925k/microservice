package main

import (
	"errors"
	"github.com/4925k/microservice/logger/data"
	"log"
	"net/http"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	log.Println("log request")
	// read JSON from request
	var request JSONPayload
	_ = app.readJSON(w, r, &request)

	// insert data
	err := app.Models.LogEntry.Insert(data.LogEntry{
		Name: request.Name,
		Data: request.Data,
	})
	if err != nil {
		_ = app.writeError(w, errors.New("failed to insert log: "+err.Error()))
		return
	}

	_ = app.writeJSON(w, http.StatusAccepted, JSONResponse{
		Error:   false,
		Message: "logged",
	}, nil)
}
