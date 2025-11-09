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
	row []string
}

type editSubmitMsg struct {
	tableName string
	id        string
	columns   []database.Column
	row       []string
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

		return tableDataLoadedMsg{columns, rows, tableName}
	}
}

func selectRowCmd(row []string) tea.Cmd {
	return func() tea.Msg {
		return rowSelectedMsg{
			row: row,
		}
	}
}

func editSubmitCmd(tableName, id string, columns []database.Column, row []string) tea.Cmd {
	return func() tea.Msg {
		return editSubmitMsg{
			tableName: tableName,
			id: id,
			columns: columns,
			row: row,
		}
	}
}

func execEditCmd(m *database.Manager, msg editSubmitMsg) tea.Cmd {
	return func() tea.Msg {
		err := m.EditRow(msg.tableName, msg.id, msg.columns, msg.row)
		if err != nil {
			return errMsg{err: err}
		}
		return tableSelectedMsg{msg.tableName}
	}
}

func execQueryCmd(m *database.Manager, query string) tea.Cmd {
	return func() tea.Msg {
		columns, rows, err := m.ExecuteQuery(query)

		return queryResultMsg{columns, rows, err}
	}
}
