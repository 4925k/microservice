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

	// send mail
	mux.Post("/send", app.SendMail)

	return mux
}
