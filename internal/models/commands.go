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
	columns   []database.Column
	rows      [][]string
	tableName string
}

type queryResultMsg struct {
	columns []string
	rows    [][]string
	err     error
}

type rowSelectedMsg struct {
	tableName string
	row       []string
}

type switchTabMsg struct {
	activeTab tab // should match states
}

type errMsg struct {
	err error
}

func switchTabCmd(t tab) tea.Cmd {
	return func() tea.Msg {
		return switchTabMsg{activeTab: t}
	}
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

		return tableDataLoadedMsg{columns, rows, tableName}
	}
}

func selectRowCmd(tableName string, row []string) tea.Cmd {
	return func() tea.Msg {
		return rowSelectedMsg{
			tableName: tableName,
			row:       row,
		}
	}
}

func execQueryCmd(m *database.Manager, query string) tea.Cmd {
	return func() tea.Msg {
		columns, rows, err := m.ExecuteQuery(query)

		return queryResultMsg{columns, rows, err}
	}
}
