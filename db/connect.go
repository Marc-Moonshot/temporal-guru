package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func Connect() *pgxpool.Pool {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found. using system environment variables.")
	}

	environment := os.Getenv("environment")

	if environment == "" {
		log.Fatal("environment not set.")
	}

	var dbUrl string
	dbUrlEnvString := "DB_URL"
	if environment == "development" {
		dbUrlEnvString = "DEV_DB_URL"
	}

	dbUrl = os.Getenv(dbUrlEnvString)

	if dbUrl == "" {
		log.Fatalf("%s not set.", dbUrlEnvString)
	}

	// conn, err := pgx.Connect(context.Background(), dbUrl)
	pool, err := pgxpool.New(context.Background(), dbUrl)

	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	fmt.Println("Connected to PostgreSQL.")
	return pool
}
