package main

import (
	"encoding/json"
	"net/http"
)

type brokerResponse struct {
	Error   bool   `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	data := &brokerResponse{
		Error:   false,
		Message: "OK",
	}

	out, err := json.Marshal(data)
	if err != nil {
		// handle error
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	_, _ = w.Write(out)
}
