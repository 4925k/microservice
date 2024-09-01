package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"net/http"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	// allow all for cors
	mux.Use(cors.AllowAll().Handler)

	// health check endpoint
	mux.Use(middleware.Heartbeat("/ping"))

	mux.Get("/", func(w http.ResponseWriter, e *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Hello from auth"))
	})

	mux.Post("/authenticate", app.Authenticate)

	return mux
}
