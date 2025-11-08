package database

import (
	"database/sql"
	"fmt"
)

type Table struct {
	Name string `db:"name"`
}

// Column metadata (SQLite specific)
type Column struct {
	CID          int     // position in table
	Name         string  // title
	Type         string  // data type
	NotNull      bool    // not null constraint
	DefaultValue *string // nil if none
	PK           bool    // Primary key
}

// Returns list of table names
func (db *DB) ListTables() ([]string, error) {
	var tables []string

	query := `SELECT name FROM sqlite_master
	WHERE type ='table' AND name NOT LIKE 'sqlite_%'
	ORDER BY name;`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("Failed to get table names: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			return nil, fmt.Errorf("Failed to scan table: %w", err)
		}
		tables = append(tables, name)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("Error iterating tables: %w", err)
	}

	return tables, nil
}

// Returns all columns of table
func (db *DB) GetTableSchema(table string) ([]Column, error) {
	var cols []Column

	query := fmt.Sprintf("PRAGMA table_info(%s)", table)

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("Error getting table info: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var col Column
		var defaultVal sql.NullString
		var notNullInt int
		var pkInt int

		err := rows.Scan(
			&col.CID,
			&col.Name,
			&col.Type,
			&notNullInt,
			&defaultVal,
			&pkInt,
		)
		if err != nil {
			return nil, fmt.Errorf("Error scanning column: %w", err)
		}

		// SQLite bools are ints so we need to convert
		col.NotNull = notNullInt == 1
		col.PK = pkInt > 0

		if defaultVal.Valid {
			col.DefaultValue = &defaultVal.String
		}

		cols = append(cols, col)
	}

	return cols, rows.Err()
}
