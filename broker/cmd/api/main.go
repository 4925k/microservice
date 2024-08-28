package main

import (
	"fmt"
	"log"
	"net/http"
)

const webPort = "9998"

type Config struct {
}

func main() {
	app := Config{}

	log.Printf("Starting broker server on %s\n", webPort)

	// define server
	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	// start server
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
