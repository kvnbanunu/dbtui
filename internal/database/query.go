package database

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

// Table metadata
type TableInfo struct {
	Name        string
	RowCount    int
	ColumnCount int
	Type        string // table or view
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
func (m *Manager) ListTables() ([]string, error) {
	var tables []string

	query := `SELECT name FROM sqlite_master
	WHERE type ='table' AND name NOT LIKE 'sqlite_%'
	ORDER BY name;`

	rows, err := m.db.Query(query)
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

func (m *Manager) GetTableInfo(tableName string) (*TableInfo, error) {
	info := &TableInfo{Name: tableName}

	var tableType string
	err := m.db.QueryRow(`SELECT type FROM sqlite_master WHERE name = ?`,
		tableName).Scan(&tableType)
	if err != nil {
		return nil, fmt.Errorf("Failed to get table type: %w", err)
	}
	info.Type = tableType

	cols, err := m.GetTableSchema(tableName)
	if err != nil {
		return nil, err
	}
	info.ColumnCount = len(cols)

	// get row count, (doesn't apply to views)
	if tableType == "table" {
		var count int
		query := fmt.Sprintf("Select COUNT(*) FROM %s", quoteIdentifier(tableName))
		if err := m.db.QueryRow(query).Scan(&count); err != nil {
			return nil, fmt.Errorf("Failed to get row count: %w", err)
		}
		info.RowCount = count
	}

	return info, nil
}

// Returns all columns of table
func (m *Manager) GetTableSchema(tableName string) ([]Column, error) {
	var cols []Column

	query := fmt.Sprintf("PRAGMA table_info(%s)", tableName)

	rows, err := m.db.Query(query)
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

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error iterating column rows: %w", err)
	}

	if len(cols) == 0 {
		return nil, fmt.Errorf("Table %s not found or has no columns", tableName)
	}

	return cols, nil
}

func (m *Manager) GetTableData(tableName string, limit, offset int) ([][]string, error) {
	query := fmt.Sprintf("SELECT * FROM %s LIMIT ? OFFSET ?", quoteIdentifier(tableName))

	rows, err := m.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("Failed to query table data: %w", err)
	}
	defer rows.Close()

	// col names
	cols, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("Failed to get columns: %w", err)
	}

	res, err := extractRows(rows, cols)
	if err != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error iterating table rows: %w", err)
	}

	return res, nil
}

// returns total # of rows in a table
func (m *Manager) GetRowCount(tableName string) (int, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", quoteIdentifier(tableName))

	var count int
	err := m.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("Failed to get row count: %w", err)
	}

	return count, nil
}

// search rows in a table
func (m *Manager) SearchTable(tableName, term string, limit, offset int) ([][]string, error) {
	if term == "" {
		return m.GetTableData(tableName, limit, offset)
	}

	cols, err := m.GetTableSchema(tableName)
	if err != nil {
		return nil, err
	}

	var conditions []string
	for _, col := range cols {
		// only look in TEXT type cols
		if strings.Contains(strings.ToUpper(col.Type), "TEXT") ||
			strings.Contains(strings.ToUpper(col.Type), "CHAR") ||
			col.Type == "" { // *** SQLITE allows untyped cols
			conditions = append(conditions, fmt.Sprintf("%s LIKE ?", quoteIdentifier(col.Name)))
		}
	}

	// no text cols
	if len(conditions) == 0 {
		return m.GetTableData(tableName, limit, offset)
	}

	where := strings.Join(conditions, " OR ")
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s LIMIT ? OFFSET ?",
		quoteIdentifier(tableName),
		where,
	)

	pattern := "%" + term + "%"
	args := make([]any, len(conditions)+2)
	for i := 0; i < len(conditions); i++ {
		args[i] = pattern
	}
	args[len(conditions)] = limit
	args[len(conditions)+1] = offset

	rows, err := m.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("Failed to search table: %w", err)
	}
	defer rows.Close()

	cols2, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("Failed to get columns: %w", err)
	}

	res, err := extractRows(rows, cols2)
	if err != nil {
		return nil, err
	}
	return res, rows.Err()
}

func (m *Manager) GetDBInfo() (map[string]string, error) {
	info := make(map[string]string)
	var pageSize int
	var pageCount int
	var encoding string
	var fkEnabled int

	info["path"] = m.path

	if err := m.db.QueryRow("PRAGMA page_size").Scan(&pageSize); err != nil {
		return nil, err
	}
	info["page_size"] = fmt.Sprintf("%d bytes", pageSize)

	if err := m.db.QueryRow("PRAGMA page_count").Scan(&pageCount); err != nil {
		return nil, err
	}
	info["page_count"] = fmt.Sprintf("%d", pageCount)

	dbSize := pageSize * pageCount
	info["size"] = formatBytes(dbSize)

	if err := m.db.QueryRow("PRAGMA encoding").Scan(&encoding); err != nil {
		return nil, err
	}
	info["encoding"] = encoding

	if err := m.db.QueryRow("PRAGMA foreign_keys").Scan(&fkEnabled); err != nil {
		return nil, err
	}
	if fkEnabled == 1 {
		info["foreign_keys"] = "enabled"
	} else {
		info["foreign_keys"] = "disabled"
	}

	tables, err := m.ListTables()
	if err != nil {
		return nil, err
	}
	info["table_count"] = fmt.Sprintf("%d", len(tables))

	return info, nil
}

// execs a custom sql query and returns the results
func (m *Manager) ExecuteQuery(query string) ([]string, [][]string, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, nil, fmt.Errorf("Error empty query")
	}

	qUpper := strings.ToUpper(query)
	isSelect := strings.HasPrefix(qUpper, "SELECT") ||
	strings.HasPrefix(qUpper, "PRAGMA") ||
	strings.HasPrefix(qUpper, "EXPLAIN")

	if !isSelect {
		// insert, update, delete...
		res, err := m.db.Exec(query)
		if err != nil {
			return nil, nil, fmt.Errorf("Failed to execute query: %w", err)
		}
		rowsAffected, _ := res.RowsAffected()
		lastInsertId, _ := res.LastInsertId()

		cols := []string{"Result", "Rows Affected", "Last Insert ID"}
		row := []string{
			"Success",
			fmt.Sprintf("%d", rowsAffected),
			fmt.Sprintf("%d", lastInsertId),
		}
		return cols, [][]string{row}, nil
	}

	// select
	rows, err := m.db.Query(query)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to execute query: %w", err)
	}

	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to ge columns: %w", err)
	}

	res, err := extractRows(rows, cols)
	if err != nil {
		return nil, nil, err
	}
	return cols, res, nil
}

func (m *Manager) EditRow(tableName, id string, columns []Column, row []string) error {
	query := "UPDATE %s SET %s WHERE id = ?"

	parts := make([]string, len(columns))
	args := make([]any, len(columns) + 1)

	for i, col := range columns {
		parts[i] = fmt.Sprintf("%s = ?", quoteIdentifier(col.Name))
		args[i] = stringToValue(col.Type, row[i])
	}

	args[len(columns)] = id

	query = fmt.Sprintf(
		query, quoteIdentifier(tableName),
		strings.Join(parts, ", "),
	)

	res, err := m.db.Exec(query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("No rows affected")
	}
	
	return nil
}

func extractRows(rows *sql.Rows, cols []string) ([][]string, error) {
	var res [][]string
	for rows.Next() {
		values := make([]any, len(cols))
		valuePtrs := make([]any, len(cols))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		row := make([]string, len(cols))
		for i, val := range values {
			row[i] = valToString(val)
		}

		res = append(res, row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return res, nil
}

// replaces any quotes in the name with double quotes (SQLite escape)
func quoteIdentifier(name string) string {
	e := strings.ReplaceAll(name, `"`, `""`)
	return fmt.Sprintf(`"%s"`, e)
}

// convert db vals to strings
func valToString(val any) string {
	if val == nil {
		return "NULL"
	}

	switch v := val.(type) {
	case []byte: // blob represented as len or hex
		if len(v) > 50 {
			return fmt.Sprintf("<BLOB: %d bytes>", len(v))
		}
		return fmt.Sprintf("%x", v)
	case string:
		return v
	case int64:
		return fmt.Sprintf("%d", v)
	case float64:
		return fmt.Sprintf("%g", v)
	case bool:
		if v {
			return "1"
		}
		return "0"
	default: // shouldn't reach
		return fmt.Sprintf("%v", v)
	}
}

func stringToValue(colType string, value string) any {
	// Handle NULL
	if value == "" || value == "NULL" {
		return nil
	}
	
	typeUpper := strings.ToUpper(colType)
	
	// INTEGER types
	if strings.Contains(typeUpper, "INT") {
		if v, err := strconv.ParseInt(value, 10, 64); err == nil {
			return v
		}
	}
	
	// REAL/FLOAT types
	if strings.Contains(typeUpper, "REAL") || 
	   strings.Contains(typeUpper, "FLOAT") || 
	   strings.Contains(typeUpper, "DOUBLE") ||
	   strings.Contains(typeUpper, "DECIMAL") {
		if v, err := strconv.ParseFloat(value, 64); err == nil {
			return v
		}
	}
	
	// BOOLEAN (SQLite stores as INTEGER 0/1)
	if strings.Contains(typeUpper, "BOOL") {
		if value == "true" || value == "1" || value == "TRUE" {
			return 1
		}
		if value == "false" || value == "0" || value == "FALSE" {
			return 0
		}
	}
	
	// Default to TEXT
	return value
}

// formats bytes into readable format
func formatBytes(bytes int) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div := int64(unit)
	exp := 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
