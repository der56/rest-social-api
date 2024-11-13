package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
)

var Pool *pgxpool.Pool
var mutex sync.Mutex
var initialized bool

func InitDB() error {
	mutex.Lock()
	defer mutex.Unlock()

	if initialized {
		return nil
	}

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file", err)
		return err
	}

	dbURL := os.Getenv("POSTGRES_URL")
	if dbURL == "" {
		log.Fatal("You must set your 'POSTGRES_URL' environment variable.")
		return fmt.Errorf("POSTGRES_URL is not set")
	}

	poolConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatal("Error parsing database URL: ", err)
		return err
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
		return err
	}

	Pool = pool
	initialized = true

	fmt.Println("Database connected successfully.")
	return nil
}

func CloseDB() {
	if Pool != nil {
		Pool.Close()
	}
}
