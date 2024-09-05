package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action,omitempty"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

// Broker will handle the broker request
func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	data := &serviceResponse{
		Error:   false,
		Message: "OK",
	}

	_ = app.writeJSON(w, http.StatusOK, data)
}

// HandleSubmission will handle the broker request
func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var req RequestPayload

	err := app.readJSON(w, r, &req)
	if err != nil {
		_ = app.writeError(w, err)
		return
	}

	switch req.Action {
	case "auth":
		app.authenticate(w, req.Auth)
		return
	default:
		_ = app.writeError(w, errors.New("action not supported"), http.StatusBadRequest)
		return
	}
}

// authenticate will authenticate the user
func (app *Config) authenticate(w http.ResponseWriter, data AuthPayload) {
	// create json to sent to auth service
	jsonData, _ := json.Marshal(data)

	// call the service
	req, err := http.NewRequest("POST", "http://auth/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		_ = app.writeError(w, err)
		return
	}

	// get response from auth service
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		_ = app.writeError(w, err)
		return
	}
	defer res.Body.Close()

	// read status
	if res.StatusCode == http.StatusUnauthorized {
		_ = app.writeError(w, errors.New("invalid credentials"))
		return
	} else if res.StatusCode != http.StatusOK {
		_ = app.writeError(w, errors.New("error calling auth service"))
		return
	}

	// read body
	var authResponse serviceResponse
	err = json.NewDecoder(res.Body).Decode(&authResponse)
	if err != nil {
		_ = app.writeError(w, err)
		return
	}

	// check response
	if authResponse.Error {
		_ = app.writeError(w, errors.New(authResponse.Message), http.StatusUnauthorized)
		return
	}

	// return response
	_ = app.writeJSON(w, http.StatusOK, serviceResponse{
		Error:   false,
		Message: "Authenticated",
		Data:    authResponse.Data,
	})
}
