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
	// Just return 'nil', which means "no I/O right now, please."
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:
		key := msg.String()
		switch m.state {
		case fullView:
			switch key {
			case "q":
				return m, tea.Quit
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}
