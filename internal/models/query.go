package models

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type queryModel struct {
	input   textinput.Model
	table   table.Model
	results [][]string
	err     error
	width   int
	height  int
}

func newQueryModel() queryModel {
	ti := textinput.New()
	ti.Placeholder = "SELECT * FROM table_name"
	ti.Focus()
	ti.Width = 50

	return queryModel{
		input: ti,
		table: table.New(),
	}
}

func (q *queryModel) setSize(width, height int) {
	q.width = width
	q.height = height
	q.input.Width = width - 4
}

func (m *queryModel) setResults(columns []string, rows [][]string, err error) {
	m.err = err
	m.results = rows

	if err == nil && len(rows) > 0 {
		tableCols := make([]table.Column, len(columns))
		for i, col := range columns {
			width := 15
			if len(col) > width {
				width = len(col) + 2
			}

			tableCols[i] = table.Column{
				Title: col,
				Width: width,
			}
		}

		// convert rows to bubbles
		tableRows := make([]table.Row, len(rows))
		for i, row := range rows {
			tableRows[i] = row
		}

		m.table.SetColumns(tableCols)
		m.table.SetRows(tableRows)
		m.table.SetHeight(m.height - 5)

	}
}

func (q queryModel) Init() tea.Cmd {
	return nil
}

func (q queryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+u":
			q.input.Reset()
			return q, nil
		}
	}

	q.input, cmd = q.input.Update(msg)
	return q, cmd
}

func (q queryModel) View() string {
	view := "SQL Query\n\n"
	view += q.input.View() + "\n\n"

	if q.err != nil {
		view += fmt.Sprintf("Error: %s\n", q.err)
	} else if len(q.results) > 0 {
		view += q.table.View()
	}

	return view
}
