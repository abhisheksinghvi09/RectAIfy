package schema

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations.sql
var migrationSQL string

// Migrate runs database migrations
func Migrate(ctx context.Context, db *pgxpool.Pool) error {
	conn, err := db.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Execute migration SQL
	_, err = tx.Exec(ctx, migrationSQL)
	if err != nil {
		return fmt.Errorf("failed to execute migrations: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit migrations: %w", err)
	}

	return nil
}

// InitDatabase initializes a new database connection pool
func InitDatabase(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	// Configure pool settings
	config.MaxConns = 10
	config.MinConns = 1

	db, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test the connection
	if err := db.Ping(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// CreateDatabase creates the database if it doesn't exist
func CreateDatabase(ctx context.Context, dsn, dbName string) error {
	// Parse DSN to connect to postgres database instead of target database
	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		return fmt.Errorf("failed to parse DSN: %w", err)
	}

	// Connect to postgres database to create target database
	config.Database = "postgres"
	conn, err := pgx.ConnectConfig(ctx, config)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres database: %w", err)
	}
	defer conn.Close(ctx)

	// Check if database exists
	var exists bool
	err = conn.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", dbName).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if database exists: %w", err)
	}

	if !exists {
		// Create database
		_, err = conn.Exec(ctx, fmt.Sprintf("CREATE DATABASE %s", dbName))
		if err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		}
	}

	return nil
}
