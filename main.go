package main

import (
	"microservice/config"
	"microservice/handler"
	"microservice/logging"
	"microservice/server"
	"net/http"
)

func main() {
	// read config file
	cfg := config.Make()

	// set up logging with standard flags and file name
	logger := logging.New(cfg.Log)

	// creating a router to route the incoming requets
	mux := http.NewServeMux()

	// get a handler with logger passed
	handler := handler.New(logger)

	// set up routes
	handler.Routes(mux)

	srv := server.New(mux, cfg.Addr)
	logger.Printf("[INFO] starting server at %s\n", cfg.Addr)
	logger.Fatalf("[FATAL] server failed to start: %v", srv.ListenAndServeTLS(cfg.Crt, cfg.Key))
}
