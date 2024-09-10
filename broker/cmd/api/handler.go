package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

// RequestPayload is the type for the request
type RequestPayload struct {
	Action string      `json:"action,omitempty"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

type MailPayload struct {
	From    string `json:"from,omitempty"`
	To      string `json:"to,omitempty"`
	Subject string `json:"subject,omitempty"`
	Message string `json:"message,omitempty"`
}

// AuthPayload is the type for the auth
type AuthPayload struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

// LogPayload is the type for the log
type LogPayload struct {
	Name string `json:"name,omitempty"`
	Data string `json:"data,omitempty"`
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
	case "logger":
		app.log(w, req.Log)
		return
	case "mail":
		app.sendMail(w, req.Mail)
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

	// log authentication
	err = app.logRequest(data.Email, data.Password)
	if err != nil {
		_ = app.writeError(w, err)
		return
	}

	// return response
	_ = app.writeJSON(w, http.StatusOK, serviceResponse{
		Error:   false,
		Message: "Authenticated",
		Data:    authResponse.Data,
	})
}

// log will log the request
func (app *Config) log(w http.ResponseWriter, data LogPayload) {
	// get data
	jsonData, _ := json.Marshal(data)

	// call the service
	req, err := http.NewRequest("POST", "http://logger/log", bytes.NewBuffer(jsonData))
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
	if res.StatusCode != http.StatusAccepted {
		_ = app.writeError(w, errors.New("error calling logger service"))
		return
	}

	// return response
	_ = app.writeJSON(w, http.StatusOK, serviceResponse{
		Error:   false,
		Message: "Logged",
	})
}

// logRequest will log the given name and data
func (app *Config) logRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = name
	entry.Data = data

	// create json to sent to auth service
	jsonData, _ := json.Marshal(entry)

	// logger service url
	logServiceURL := "http://logger/log"

	// create request
	req, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	// complete request
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

// mailRequest will send an email
func (app *Config) sendMail(w http.ResponseWriter, data MailPayload) {
	// parse data
	jsonData, _ := json.Marshal(data)

	// call mail service
	mailServiceURL := "http://mail/send"

	// create request
	req, err := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		_ = app.writeError(w, err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		_ = app.writeError(w, err)
		return
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusAccepted {
		_ = app.writeError(w, errors.New("error calling mail service"))
		return
	}

	// return response
	var payload serviceResponse
	payload.Error = false
	payload.Message = "Message sent to " + data.To

	_ = app.writeJSON(w, http.StatusAccepted, payload)
}
