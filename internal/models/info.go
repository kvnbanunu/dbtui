package models

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

func (m *model) infoView() string {
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
		baseStyle.Render(m.infoTable.View()),
	)
}

func (m *model) setInfoTable() {
	m.infoTable = newTable()

	tableCols := []table.Column{
		{Title: "CID", Width: 5},
		{Title: "Name", Width: 15},
		{Title: "Type", Width: 10},
		{Title: "NotNull", Width: 10},
		{Title: "DefaultValue", Width: 15},
		{Title: "PK", Width: 5},
	}

	m.infoTable.SetColumns(tableCols)
	
	rows := make([]table.Row, len(m.columns))
	for i, col := range m.columns {
		defStr := "NULL"
		if col.DefaultValue != nil {
			defStr = *col.DefaultValue
		}
		rows[i] = []string{
			fmt.Sprintf("%d", col.CID),
			col.Name,
			col.Type,
			fmt.Sprintf("%v", col.NotNull),
			defStr,
			fmt.Sprintf("%v", col.PK),
		}
	}
	m.infoTable.SetRows(rows)
	m.infoTable.GotoTop()
	m.infoTable.Focus()
}
