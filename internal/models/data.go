package models

import (
	"fmt"

	"dbtui/internal/database"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

func (m *model) dataView() string {
	if m.name == "" {
		return "No table selected"
	}

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("170")).
		Render(fmt.Sprintf("Table: %s (Page %d)", m.name, m.currentPage+1))

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		baseStyle.Render(m.dataTable.View()),
	)
}

func (m *model) setDataTable(tableName string, columns []database.Column, rows [][]string) {
	m.name = tableName
	m.currentPage = 0
	m.columns = columns
	m.activeTab = dataTab

	m.dataTable = newTable()

	// convert cols to bubbles
	tableCols := make([]table.Column, len(columns))
	for i, col := range columns {
		width := 10
		if len(col.Name) > width {
			width = len(col.Name) + 2
		}

		tableCols[i] = table.Column{
			Title: col.Name,
			Width: width,
		}
	}

	m.dataTable.SetColumns(tableCols)

	// convert rows to bubbles
	tableRows := make([]table.Row, len(rows))
	for i, row := range rows {
		tableRows[i] = row
	}

	m.dataTable.SetRows(tableRows)
	m.dataTable.GotoTop()
	m.dataTable.Focus()
}
