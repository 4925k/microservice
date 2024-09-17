package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"net/http"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	// allow all for CORS
	mux.Use(cors.AllowAll().Handler)

	// health check endpoint
	mux.Use(middleware.Heartbeat("/ping"))

	// root page
	mux.Post("/", app.Broker)

	// entry point
	mux.Post("/handle", app.HandleSubmission)

	// grpc endpoint
	mux.Post("/log-grpc", app.logViaGRPC)

	return mux
}
