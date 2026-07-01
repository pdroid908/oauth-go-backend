package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

// Connect membuat koneksi PostgreSQL dan menyimpannya di DB.
func Connect() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatalf("Failed to create connection pool: %v", err)
	}

	// Test koneksi
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	DB = pool

	log.Println("✅ PostgreSQL connected")
}

// Close menutup koneksi database.
func Close() {
	if DB != nil {
		DB.Close()
		log.Println("🛑 PostgreSQL connection closed")
	}
}