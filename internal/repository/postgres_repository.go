package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
)

// Repository defines the interface for data access operations
type Repository interface {
	Initkey_pool(ctx context.Context, keys []string) error
	GetNextAvailableKey(ctx context.Context, original_url string) (string, error)
	LinkURL(ctx context.Context, short_url, original_url string) error
	Getoriginal_url(ctx context.Context, short_url string) (string, error)
}

// PostgresRepository implements Repository interface using PostgreSQL
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository creates a new PostgresRepository instance
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// Initkey_pool bulk inserts all provided keys using PostgreSQL COPY for maximum efficiency
// COPY is 10-100x faster than individual INSERTs (50K keys inserted in ~100-200ms instead of 30+ seconds)
// All keys are inserted with is_used = false and current timestamp
func (pr *PostgresRepository) Initkey_pool(ctx context.Context, keys []string) error {
	if len(keys) == 0 {
		return fmt.Errorf("no keys provided")
	}

	tx, err := pr.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Prepare COPY statement
	// COPY is PostgreSQL's bulk insert mechanism - streams data directly into table
	stmt, err := tx.Prepare(pq.CopyIn("key_pool", "short_url", "is_used", "created_at"))
	if err != nil {
		return fmt.Errorf("failed to prepare COPY statement: %w", err)
	}
	defer stmt.Close()

	// Insert all keys via COPY
	now := time.Now()
	for _, key := range keys {
		_, err := stmt.ExecContext(ctx, key, false, now)
		if err != nil {
			return fmt.Errorf("failed to copy key %s: %w", key, err)
		}
	}

	// Flush the COPY buffer - this actually sends all data to database
	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to flush COPY buffer: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetNextAvailableKey atomically fetches an unused key and marks it as used
// Uses SELECT...FOR UPDATE to prevent race conditions in concurrent access
func (pr *PostgresRepository) GetNextAvailableKey(ctx context.Context, original_url string) (string, error) {
	tx, err := pr.db.BeginTx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Safety net: automatically rolls back if anything fails below

	// 2. Select and lock an unused key from the pool
	var short_url string
	err = tx.QueryRowContext(ctx,
		`SELECT "short_url" FROM key_pool 
         WHERE "is_used" = false 
         LIMIT 1 
         FOR UPDATE SKIP LOCKED`).Scan(&short_url)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("no available keys in pool")
		}
		return "", fmt.Errorf("failed to fetch key: %w", err)
	}

	// 3. Mark that key as used
	_, err = tx.ExecContext(ctx,
		`UPDATE key_pool SET "is_used" = true WHERE "short_url" = $1`, short_url)
	if err != nil {
		return "", fmt.Errorf("failed to mark key as used: %w", err)
	}

	// 4. Link it to the original URL in the key_link table right here!
	now := time.Now()
	_, err = tx.ExecContext(ctx,
		`INSERT INTO key_link ("short_url", "original_url", "created_at") 
         VALUES ($1, $2, $3)`,
		short_url, original_url, now)
	if err != nil {
		return "", fmt.Errorf("failed to link URL: %w", err)
	}

	// 5. Commit everything together smoothly
	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	return short_url, nil
}

// LinkURL creates a mapping between a short URL and the original URL
// This records the shortened URL to original URL relationship
func (pr *PostgresRepository) LinkURL(ctx context.Context, short_url string, original_url string) error {
	now := time.Now()
	_, err := pr.db.ExecContext(ctx,
		`INSERT INTO key_link (short_url, original_url, created_at) 
		 VALUES ($1, $2, $3)`,
		short_url, original_url, now)
	if err != nil {
		return fmt.Errorf("failed to link URL: %w", err)
	}
	return nil
}

// Getoriginal_url retrieves the original URL for a given short URL
func (pr *PostgresRepository) Getoriginal_url(ctx context.Context, short_url string) (string, error) {
	var original_url string
	err := pr.db.QueryRowContext(ctx,
		`SELECT original_url FROM key_link WHERE short_url = $1`,
		short_url).Scan(&original_url)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("short URL not found: %s", short_url)
		}
		return "", fmt.Errorf("failed to get original URL: %w", err)
	}
	return original_url, nil
}
