package models

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type tableList struct {
	list     list.Model // list of table names
	tables   []string
	width    int
	height   int
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
		list: list,
	}
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
	return tl.list.View()
}
