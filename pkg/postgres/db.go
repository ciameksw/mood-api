package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type PostgresDB struct {
	DB *sql.DB
}

const (
	maxRetries        = 5
	initialBackoff    = 1 * time.Second
	backoffMultiplier = 2.0
)

func Connect(host, port, user, password, dbname, sslmode string) (*PostgresDB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	var db *sql.DB
	var err error
	backoff := initialBackoff

	for attempt := 1; attempt <= maxRetries; attempt++ {
		db, err = attemptConnect(connStr, 10*time.Second)
		if err == nil {
			log.Printf("Successfully connected to PostgreSQL on attempt %d", attempt)
			// Set connection pool settings
			db.SetMaxOpenConns(25)
			db.SetMaxIdleConns(5)
			db.SetConnMaxLifetime(5 * time.Minute)
			return &PostgresDB{DB: db}, nil
		}

		if attempt < maxRetries {
			log.Printf("Failed to connect to PostgreSQL (attempt %d/%d): %v. Retrying in %v...", attempt, maxRetries, err, backoff)
			time.Sleep(backoff)
			backoff = time.Duration(float64(backoff) * backoffMultiplier)
		}
	}

	return nil, fmt.Errorf("failed to connect to PostgreSQL after %d attempts: %w", maxRetries, err)
}

func attemptConnect(connStr string, timeout time.Duration) (*sql.DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func (p *PostgresDB) Disconnect(ctx context.Context) {
	// Wait for all connections to be returned to the pool
	p.DB.SetMaxIdleConns(0)
	p.DB.SetMaxOpenConns(0)

	err := p.DB.Close()
	if err != nil {
		log.Printf("Error closing database connection: %v", err)
	}

	log.Printf("Successfully disconnected from PostgreSQL")

	// Ensure context is used
	<-ctx.Done()
}
