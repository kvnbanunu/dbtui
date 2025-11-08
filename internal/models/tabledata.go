package models

import (
	"dbtui/internal/database"

	"github.com/charmbracelet/bubbles/table"
)

type tableData struct {
	table   table.Model
	columns []database.Column
}
