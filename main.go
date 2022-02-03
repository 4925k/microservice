package main

import (
	"microservice/handler"
	"microservice/logging"
	"microservice/server"
	"net/http"
)

const (
	addr = ":8080"
	crt  = "./certs/server.crt"
	key  = "./certs/server.key"
)

func main() {
	// set up logging with standard flags and file name
	logger := logging.New()

	// creating a router to route the incoming requets
	mux := http.NewServeMux()

	// get a handler with logger passed
	handler := handler.New(logger)

	// set up routes
	handler.Routes(mux)

	srv := server.New(mux, addr)
	logger.Printf("[INFO] starting server at %s\n", addr)
	logger.Fatalf("[FATAL] server failed to start: %v", srv.ListenAndServeTLS(crt, key))
}
