package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
	webPort = "80"
)

type Config struct {
}

func main() {
	log.Println("mail service started")

	app := Config{}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	log.Println("Starting mail service on port", webPort)
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}
