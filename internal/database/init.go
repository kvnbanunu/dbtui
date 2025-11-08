package database

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

type DB struct {
	*sql.DB
}

func Init(path string) (*DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	conn := &DB{db}

	return conn, nil
}
