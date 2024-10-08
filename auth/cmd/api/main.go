package main

import (
	"database/sql"
	"fmt"
	"github.com/microservice/auth/data"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const (
	webPort = "80"
)

var (
	dbCount int64
)

type Config struct {
	Repo data.Repository
}

func main() {
	log.Printf("Starting auth api at %s\n", webPort)

	// connect to DB
	conn := connectToDB()
	if conn == nil {
		log.Fatal("Unable to connect to database")
	}

	// set up config
	app := Config{}

	// server config
	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	// start server
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

// openDB helps to get a db connection using the given dsn
func openDB(dsn string) (*sql.DB, error) {
	// open connection
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	// test connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// connectToDB helps to connect to database
func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		db, err := openDB(dsn)
		if err != nil {
			log.Printf("Error connecting to database: %v\n", err)
			dbCount++
		} else {
			return db
		}

		if dbCount > 10 {
			log.Printf("Cannot connect to database after 10 retries\n")
			return nil
		}

		log.Printf("Waiting for database to become available...\n")
		time.Sleep(5 * time.Second)
		continue
	}
}

// setupRepo sets up the repository
func (app *Config) setupRepo(conn *sql.DB) {
	db := data.NewPostgresRepository(conn)

	app.Repo = &db
}
