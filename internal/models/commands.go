package models

import (
	"dbtui/internal/database"

	tea "github.com/charmbracelet/bubbletea"
)

type tablesLoadedMsg struct {
	tables []string
}

type tableSelectedMsg struct {
	tableName string
}

type tableDataLoadedMsg struct {
	columns []database.Column
	rows    [][]string
}

type queryResultMsg struct {
	columns []string
	rows    [][]string
	err     error
}

type errMsg struct {
	err error
}

func loadTablesCmd(m *database.Manager) tea.Cmd {
	return func() tea.Msg {
		tables, err := m.ListTables()
		if err != nil {
			return errMsg{err}
		}
		return tablesLoadedMsg{tables}
	}
}

func selectTableCmd(name string) tea.Cmd {
	return func() tea.Msg {
		return tableSelectedMsg{tableName: name}
	}
}

func loadTableDataCmd(m *database.Manager, tableName string, offset int) tea.Cmd {
	return func() tea.Msg {
		columns, err := m.GetTableSchema(tableName)
		if err != nil {
			return errMsg{err}
		}

		rows, err := m.GetTableData(tableName, 100, offset)
		if err != nil {
			return errMsg{err}
		}

		return tableDataLoadedMsg{columns, rows}
	}
}

func execQueryCmd(m *database.Manager, query string) tea.Cmd {
	return func() tea.Msg {
		columns, rows, err := m.ExecuteQuery(query)

		return queryResultMsg{columns, rows, err}
	}
}
