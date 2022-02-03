package logging

import (
	"log"
	"os"
)

func New() *log.Logger {
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	// If the file doesn't exist, create it or append to the file
	file, err := os.OpenFile("./logging/microservice.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	logger.SetOutput(file)

	return logger
}
