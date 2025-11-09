package models

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tab int

const (
	dataTab tab = iota
	infoTab
	queryTab
)

type tabsModel struct {
	tabs      []string
	activeTab tab
}

func newTabsModel() tabsModel {
	return tabsModel{
		tabs:      []string{"Data", "Info", "Query"},
		activeTab: dataTab,
	}
}

func (t tabsModel) Init() tea.Cmd {
	return nil
}

func (t tabsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Left):
			t.activeTab = max(t.activeTab-1, 0)
			return t, switchTabCmd(t.activeTab)
		case key.Matches(msg, keys.Right):
			t.activeTab = min(t.activeTab+1, tab(len(t.tabs)-1))
			return t, switchTabCmd(t.activeTab)
		}
	}
	return t, nil
}

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

var (
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	docStyle          = lipgloss.NewStyle().Padding(1, 2, 1, 2)
	highlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Border(activeTabBorder, true)
	windowStyle       = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(2, 0).Align(lipgloss.Center).Border(lipgloss.NormalBorder()).UnsetBorderTop()
)

func (t tabsModel) View() string {
	doc := strings.Builder{}

	var renderedTabs []string

	for i, tab := range t.tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(t.tabs)-1, i == int(t.activeTab)
		if isActive {
			style = activeTabStyle
		} else {
			style = inactiveTabStyle
		}
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "│"
		} else if isLast && !isActive {
			border.BottomRight = "┤"
		}
		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(tab))
	}
	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")
	// doc.WriteString(windowStyle.Width((lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize())).Render(t.tabs[t.activeTab]))
	return docStyle.Render(doc.String())
}
