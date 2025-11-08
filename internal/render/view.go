package render

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var fullViewStyle = lipgloss.NewStyle().Background(lipgloss.Color("99")).Padding(0, 1)

func (m Model) View() string {
	// The header
	s := "Database Tables\n"

	for i, t := range m.tables {
		s += fullViewStyle.Render(fmt.Sprintf("%d: %s\n", i, t.Name))
	}

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}
