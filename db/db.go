package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB_CONN *pgxpool.Pool

func InitializeDatabase() {
	var err error

	const DATABASE_CONN_URL = "postgresql://harinarayananks@localhost:5432/postgres"

	DB_CONN, err = pgxpool.New(context.Background(), DATABASE_CONN_URL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("connected to db")

}
