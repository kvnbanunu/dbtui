package models

import (
	"dbtui/internal/database"

	tea "github.com/charmbracelet/bubbletea"
)

type tablesLoadedMsg struct {
	tables []string
}

type errMsg struct {
	err error
}

func loadTablesCmd(db *database.DB) tea.Cmd {
	return func() tea.Msg {
		tables, err := db.ListTables()
		if err != nil {
			return errMsg{err}
		}
		return tablesLoadedMsg{tables}
	}
}
