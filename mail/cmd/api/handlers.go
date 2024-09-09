package main

import (
	"fmt"
	"net/http"
)

func (app *Config) SendMail(w http.ResponseWriter, r *http.Request) {
	// input struct
	type mailMessage struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	// parse request
	var requestPayload mailMessage
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		_ = app.writeError(w, err)
		return
	}

	// create message
	msg := Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	// send mail
	err = app.Mailer.SendSMTPMessage(msg)
	if err != nil {
		_ = app.writeError(w, err)
		return
	}

	// create response
	payload := JSONResponse{
		Error:   false,
		Message: fmt.Sprintf("Email sent to %s", msg.To),
		Data:    nil,
	}

	_ = app.writeJSON(w, http.StatusAccepted, payload, nil)
}
