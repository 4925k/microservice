package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"testing"
)

// TestRoutes tests the routes for the api
func Test_routes_exist(t *testing.T) {
	testApp := Config{}

	testRoutes := testApp.routes()

	chiRoutes := testRoutes.(chi.Router)

	routes := []string{"/authenticate"}

	for _, route := range routes {
		routeExists(t, chiRoutes, route)
	}
}

// routeExists checks if a route exists in the chi router
func routeExists(t *testing.T, routes chi.Router, route string) {
	found := false

	_ = chi.Walk(routes, func(method string, foundRoute string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if route == foundRoute {
			found = true
		}

		return nil
	})

	if !found {
		t.Errorf("did not find %s in registered routes\n", route)
	}
}
