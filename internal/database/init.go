package database

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type Manager struct {
	db   *sql.DB
	path string
}

func NewManager(path string) (*Manager, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("Failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("Failed to ping database: %w", err)
	}

	return &Manager{
		db: db,
		path: path,
	}, nil
}

func (m *Manager) Close() error {
	return m.db.Close()
}
