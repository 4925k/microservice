package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/4925k/microservice/broker/event"
	"github.com/4925k/microservice/broker/logs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"net/rpc"
	"time"
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
		app.logViaRPC(w, req.Log) // logging via rpc service
		//app.logViaRabbit(w, req.Log) // log view rabbit mq queue
		//app.log(w, req.Log) // call directly to the logger service
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
//func (app *Config) log(w http.ResponseWriter, data LogPayload) {
//	// get data
//	jsonData, _ := json.Marshal(data)
//
//	// call the service
//	req, err := http.NewRequest("POST", "http://logger/log", bytes.NewBuffer(jsonData))
//	if err != nil {
//		_ = app.writeError(w, err)
//		return
//	}
//
//	// get response from auth service
//	client := &http.Client{}
//	res, err := client.Do(req)
//	if err != nil {
//		_ = app.writeError(w, err)
//		return
//	}
//	defer res.Body.Close()
//
//	// read status
//	if res.StatusCode != http.StatusAccepted {
//		_ = app.writeError(w, errors.New("error calling logger service"))
//		return
//	}
//
//	// return response
//	_ = app.writeJSON(w, http.StatusOK, serviceResponse{
//		Error:   false,
//		Message: "Logged",
//	})
//}

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

// logViaRabbit will log data via rabbit
func (app *Config) logViaRabbit(w http.ResponseWriter, data LogPayload) {
	err := app.pushToQueue(data.Name, data.Data)
	if err != nil {
		_ = app.writeError(w, err)
		return
	}

	var payload serviceResponse
	payload.Error = false
	payload.Message = "Logged via RabbitMQ"

	_ = app.writeJSON(w, http.StatusAccepted, payload)
}

// pushToQueue will push data to queue
func (app *Config) pushToQueue(name, msg string) error {
	// connect to rabbit
	emitter, err := event.NewEmitter(app.Rabbit)
	if err != nil {
		return err
	}

	// prepare payload
	payload := LogPayload{
		Name: name,
		Data: msg,
	}

	j, _ := json.Marshal(payload)

	// log event
	err = emitter.Push(string(j), "log.INFO")
	if err != nil {
		return err
	}

	return nil
}

type RPCPayload struct {
	Name string
	Data string
}

// logViaRPC will log data via RPC
func (app *Config) logViaRPC(w http.ResponseWriter, data LogPayload) {
	// connect to rpc
	client, err := rpc.Dial("tcp", "logger:5001")
	if err != nil {
		_ = app.writeError(w, err)
		return
	}

	// prepare payload
	rpcPayload := RPCPayload{
		Name: data.Name,
		Data: data.Data,
	}

	// call rpc
	var result string
	err = client.Call("RPCServer.LogInfo", rpcPayload, &result)
	if err != nil {
		_ = app.writeError(w, err)
		return
	}

	// return response
	_ = app.writeJSON(w, http.StatusAccepted, serviceResponse{
		Error:   false,
		Message: "Logged via RPC: " + result,
	})
}

// logViaGRPC will log data via gRPC
func (app *Config) logViaGRPC(w http.ResponseWriter, r *http.Request) {

	// read json
	var requestPayload RequestPayload
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		_ = app.writeError(w, err)
		return
	}

	// call grpc
	conn, err := grpc.NewClient("logger:50001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		_ = app.writeError(w, err)
		return
	}
	defer conn.Close()

	c := logs.NewLogServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = c.WriteLog(ctx, &logs.LogRequest{
		LogEntry: &logs.Log{
			Name: requestPayload.Log.Name,
			Data: requestPayload.Log.Data,
		},
	})
	if err != nil {
		_ = app.writeError(w, err)
		return
	}

	// return response
	_ = app.writeJSON(w, http.StatusAccepted, serviceResponse{
		Error:   false,
		Message: "Logged via gRPC",
	})
}
