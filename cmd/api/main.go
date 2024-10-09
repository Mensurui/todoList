package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Mensurui/todoList/internal/data"
	_ "github.com/lib/pq"
)

type config struct {
	port int
	db   struct {
		dns string
	}
}

type application struct {
	config config
	logger *log.Logger
	models data.Models
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 8080, "The port for running the server")
	flag.StringVar(&cfg.db.dns, "dns", os.Getenv("TODOS_DB_DSN"), "DNS for the database")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	db, err := openDB(cfg)

	if err != nil {
		logger.Fatal("Error while connecting to the db")
	}

	defer db.Close()

	logger.Printf("Connected to the db: %s", cfg.db.dns)

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	srv := http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		ReadTimeout:  10 * time.Minute,
		IdleTimeout:  10 * time.Minute,
		WriteTimeout: 10 * time.Minute,
	}

	log.Printf("starting server on: %d", cfg.port)
	err = srv.ListenAndServe()
	log.Fatal(err)
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dns)

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil

}
