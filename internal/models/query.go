package models

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
)

func (m *model) queryView() string {
	view := "SQL Query\n\n"
	view += m.queryInput.View() + "\n\n"

	if m.err != nil {
		view += fmt.Sprintf("Error: %s\n", m.err)
	} else if len(m.queryResult) > 0 {
		view += m.queryTable.View()
	}

	return view
}

func (m *model) setQueryResult(columns []string, rows [][]string, err error) {
	m.err = err
	m.queryResult = rows

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

		m.queryTable.SetColumns(tableCols)
		m.queryTable.SetRows(tableRows)

	}
}
