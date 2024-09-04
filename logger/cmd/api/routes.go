package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"net/http"
)

// routes for the api
func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.AllowAll().Handler)

	mux.Use(middleware.Heartbeat("/ping"))

	mux.Post("/log", app.WriteLog)

	return mux
}
