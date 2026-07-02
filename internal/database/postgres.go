package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

// Connect membuat koneksi PostgreSQL dan menyimpannya di DB.
func Connect() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	// 1. Parse konfigurasi terlebih dahulu
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}

	// 2. PAKSA MENGGUNAKAN SIMPLE PROTOCOL
	// Ini akan menonaktifkan caching prepared statement yang menyebabkan error 42P05
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	// 3. Buat pool dengan konfigurasi yang sudah dimodifikasi
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("Failed to create connection pool: %v", err)
	}

	// Test koneksi
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	DB = pool
	log.Println("✅ PostgreSQL connected (Simple Protocol enabled)")
}

// Close menutup koneksi database.
func Close() {
	if DB != nil {
		DB.Close()
		log.Println("🛑 PostgreSQL connection closed")
	}
}