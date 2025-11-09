package models

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tableList struct {
	focus  bool
	list   list.Model // list of table names
	tables []string
	width  int
	height int
}

// implements list.Item interface
type tableItem string

func (t tableItem) FilterValue() string { return string(t) }
func (t tableItem) Title() string       { return string(t) }
func (t tableItem) Description() string { return "" }

func newTableList() tableList {
	list := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	list.SetShowHelp(false)
	return tableList{
		focus: true,
		list: list,
	}
}

func (tl *tableList) setFocus(f bool) {
	tl.focus = f
}

func (tl *tableList) setTables(tables []string) {
	tl.tables = tables

	items := make([]list.Item, len(tables))
	for i, t := range tables {
		items[i] = tableItem(t)
	}
	tl.list.SetItems(items)
	tl.list.Title = "Database Tables"
}

func (tl *tableList) setSize(width, height int) {
	tl.width = width
	tl.height = height
	tl.list.SetSize(width, height)
}

func (tl tableList) Init() tea.Cmd {
	return nil
}

func (tl tableList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Enter):
			if item, ok := tl.list.SelectedItem().(tableItem); ok {
				return tl, selectTableCmd(string(item))
			}
		}
	}

	tl.list, cmd = tl.list.Update(msg)
	return tl, cmd
}

func (tl tableList) View() string {
	style := lipgloss.NewStyle()
	
	if tl.focus {
		style = style.
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")). // Blue when focused
			Padding(0, 1)
	} else {
		style = style.
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")). // Gray when not focused
			Padding(0, 1)
	}
	
	return style.Render(tl.list.View())
}

// func (tl tableList) View() string {
// 	return tl.list.View()
// }
