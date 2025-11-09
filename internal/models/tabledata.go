package models

import (
	"fmt"

	"dbtui/internal/database"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type tableData struct {
	table       table.Model
	tableName   string
	columns     []database.Column
	currentPage int
	width       int
	height      int
}

func (td tableData) Init() tea.Cmd { return nil }

func newTableData() tableData {
	t := table.New(
		table.WithFocused(true),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return tableData{table: t}
}

func (td *tableData) setTable(name string) {
	td.tableName = name
	td.currentPage = 0
}

func (td *tableData) setData(columns []database.Column, rows [][]string) {
	td.columns = columns

	// reset rows before changing cols
	td.table.SetRows([]table.Row{})

	// convert cols to bubbles
	tableCols := make([]table.Column, len(columns))
	for i, col := range columns {
		width := 15
		if len(col.Name) > width {
			width = len(col.Name) + 2
		}

		tableCols[i] = table.Column{
			Title: col.Name,
			Width: width,
		}
	}

	td.table.SetColumns(tableCols)

	// convert rows to bubbles
	tableRows := make([]table.Row, len(rows))
	for i, row := range rows {
		tableRows[i] = row
	}

	td.table.SetRows(tableRows)
	td.table.SetHeight(td.height - 5)

	td.table.GotoTop()

	td.table.Focus()
}

func (td *tableData) setSize(width, height int) {
	td.width = width
	td.height = height
	td.table.SetWidth(width)
	td.table.SetHeight(height - 5)
}

func (td tableData) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	td.table, cmd = td.table.Update(msg)
	return td, cmd
}

func (td tableData) View() string {
	if td.tableName == "" {
		return "No table selected"
	}

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("170")).
		Render(fmt.Sprintf("Table: %s (Page %d)", td.tableName, td.currentPage+1))

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		baseStyle.Render(td.table.View()),
	)
}
