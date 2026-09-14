package db

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect() *pgxpool.Pool {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		url = "postgres://xwitter:xwitter@localhost:5432/xwitter?sslmode=disable"
	}

	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatalf("invalid DATABASE_URL: %v", err)
	}

	var pool *pgxpool.Pool
	for attempt := 1; attempt <= 30; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		pool, err = pgxpool.NewWithConfig(ctx, cfg)
		if err == nil {
			err = pool.Ping(ctx)
		}
		cancel()
		if err == nil {
			break
		}
		log.Printf("waiting for database (%d/30): %v", attempt, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	return pool
}
