package models

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

var (
	baseStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240"))
	inactiveTabBorder  = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder    = tabBorderWithBottom("┘", " ", "└")
	remainderTabBorder = tabRemainderBorder()
	docStyle           = lipgloss.NewStyle()
	highlightColor     = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle   = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle     = inactiveTabStyle.Border(activeTabBorder, true)
	remainderTabStyle  = inactiveTabStyle.Border(remainderTabBorder, true)
	windowStyle        = lipgloss.NewStyle().
				BorderForeground(highlightColor).
				Padding(1, 0).
				Align(lipgloss.Center).
				Border(lipgloss.NormalBorder()).
				UnsetBorderTop()
)

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

func tabRemainderBorder() lipgloss.Border {
	border := lipgloss.HiddenBorder()
	border.BottomLeft = "─"
	border.Bottom = "─"
	border.BottomRight = "┐"
	return border
}

// returns table with set style
func newTable() table.Model {
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

	return t
}
