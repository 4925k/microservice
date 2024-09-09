package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	webPort = "80"
)

type Config struct {
	Mailer Mail
}

func main() {
	log.Println("mail service started")

	// create app config
	app := Config{
		Mailer: createMail(),
	}

	// set up server
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

// createMail create a new mailer
func createMail() Mail {
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))

	m := Mail{
		Domain:     os.Getenv("MAIL_DOMAIN"),
		Host:       os.Getenv("MAIL_HOST"),
		Port:       port,
		Username:   os.Getenv("MAIL_USERNAME"),
		Password:   os.Getenv("MAIL_PASSWORD"),
		Encryption: os.Getenv("MAIL_ENCRYPTION"),
		SenderName: os.Getenv("MAIL_SENDER_NAME"),
		Sender:     os.Getenv("MAIL_SENDER_ADDRESS"),
	}

	return m
}
