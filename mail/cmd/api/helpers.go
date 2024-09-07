package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type JSONResponse struct {
	Error   bool   `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func (app *Config) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	var maxBytes int64 = 1048567 // one megabyte

	// read from body
	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)

	// decode data from request
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must have only a single JSON value")
	}

	return nil
}

func (app *Config) writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	// marshal data
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// set headers
	if len(headers) > 0 {
		for key, v := range headers[0] {
			w.Header()[key] = v
		}
	}

	// set content type
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// respond to request
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

// writeError will help return errors to client
func (app *Config) writeError(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	payload := JSONResponse{
		Error:   true,
		Message: err.Error(),
	}

	return app.writeJSON(w, statusCode, payload)
}
