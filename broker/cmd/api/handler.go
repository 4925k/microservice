package main

import (
	"net/http"
)

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	data := &brokerResponse{
		Error:   false,
		Message: "OK",
	}

	_ = app.writeJSON(w, r, http.StatusOK, data)
}
