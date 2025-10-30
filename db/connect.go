package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func Connect() *pgx.Conn {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found. using system environment variables.")
	}

	dbUrl := os.Getenv("DB_URL")

	if dbUrl == "" {
		log.Fatal("DB_URL not set.")
	}

	conn, err := pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	fmt.Println("Connected to PostgreSQL.")
	return conn
}
