package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

// Authenticate checks if the user exists and password matches
func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	// input structure to parse req parameters
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// read values from request
	err := app.readJSON(w, r, &input)
	if err != nil {
		_ = app.writeError(w, r, err, http.StatusBadRequest)
		return
	}

	// get related user from email
	user, err := app.Models.User.GetByEmail(input.Email)
	if err != nil {
		log.Println(err)
		_ = app.writeError(w, r, errors.New("invalid email"), http.StatusNotFound)
		return
	}

	// check for correct password
	valid, err := user.PasswordMatches(input.Password)
	if err != nil || !valid {
		_ = app.writeError(w, r, errors.New("invalid credentials"), http.StatusNotFound)
		return
	}

	// respond to client
	ret := authResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in User %s", user.Email),
		Data:    user,
	}

	_ = app.writeJSON(w, r, http.StatusOK, ret)
}
