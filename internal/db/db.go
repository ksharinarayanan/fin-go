package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB_CONN *pgxpool.Pool

func InitializeDatabase() {
	var err error

	port, isPresent := os.LookupEnv("POSTGRES_PORT")

	if !isPresent {
		port = POSTGRES_DEFAULT_PORT
	}

	DATABASE_CONN_URL := fmt.Sprintf("postgresql://harinarayananks@localhost:%v/postgres", port)

	DB_CONN, err = pgxpool.New(context.Background(), DATABASE_CONN_URL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v\n", err)
	}

	log.Println("Connected to db")
}
