package database

import "fmt"

type Table struct {
	Name string `db:"name"`
}

func (db *DB) GetTables() ([]Table, error) {
	var tables []Table

	query := `SELECT name FROM sqlite_master
	WHERE type='table' and name != 'sqlite_sequence';
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("Failed to get table names: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var table Table
		err := rows.Scan(&table.Name)
		if err != nil {
			return nil, fmt.Errorf("Failed to scan table: %w", err)
		}
		tables = append(tables, table)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("Error iterating tables: %w", err)
	}

	return tables, nil
}
