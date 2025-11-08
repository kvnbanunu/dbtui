package render

import (
	"dbtui/internal/database"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	fullView uint = iota
)

type Model struct {
	store  *database.DB
	state  uint
	tables []database.Table
}

func InitialModel(db *database.DB) (Model, error) {
	var m Model
	tables, err := db.GetTables()
	if err != nil {
		return m, err
	}

	m.store = db
	m.state = fullView
	m.tables = tables

	return m, nil
}

func (m Model) Init() tea.Cmd {
	return nil
}
