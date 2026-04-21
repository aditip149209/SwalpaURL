package storage

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(dsn string) (*Storage, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening db: %w", err)
	}

	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to db: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) SaveAvailableKeys(keys []string) error {
	if len(keys) == 0 {
		return nil
	}

	placeholders := make([]string, len(keys))
	values := make([]interface{}, len(keys))

	for i, key := range keys {
		placeholders[i] = "(?)"
		values[i] = key
	}

	query := fmt.Sprintf(
		"INSERT INTO available_keys (short_key) VALUES %s",
		strings.Join(placeholders, ","),
	)

	_, err := s.db.Exec(query, values...)
	if err != nil {
		return fmt.Errorf("failed to batch insert keys: %w", err)
	}

	return nil
}
