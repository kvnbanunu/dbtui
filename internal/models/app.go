package models

import (
	"dbtui/internal/database"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	listView int = iota
	tableView
)

type App struct {
	store          *database.Manager
	focus          int // focused
	help           help.Model
	err            error
	tableListModel tableList
	tableModel     model
	width          int
	height         int
	ready          bool
}

func NewApp(m *database.Manager) App {
	help := help.New()
	help.ShowAll = true

	return App{
		store:          m,
		focus:          listView,
		help:           help,
		tableListModel: newTableList(),
		tableModel:     newModel(m),
		ready:          false,
	}
}

func (a App) Init() tea.Cmd {
	return tea.Batch(
		loadTablesCmd(a.store),
		textinput.Blink,
	)
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var mod tea.Model
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.ready = true

		helpHeight := lipgloss.Height(a.help.View(keys))
		contentHeight := msg.Height - helpHeight - 2

		listWidth := msg.Width * 30 / 100
		contentWidth := msg.Width - listWidth

		// update all sub models with new dimensions
		a.tableListModel.setSize(listWidth, contentHeight)
		a.tableModel.setSize(contentWidth, contentHeight)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			a.store.Close()
			return a, tea.Quit

		case key.Matches(msg, keys.Back):
			if a.focus != listView {
				a.focus = listView
				a.tableListModel.setFocus(true)
			}
			return a, nil

		case key.Matches(msg, keys.Help):
			a.help.ShowAll = !a.help.ShowAll
			return a, nil
		}
	case tablesLoadedMsg:
		a.tableListModel.setTables(msg.tables)

	case tableSelectedMsg:
		a.focus = tableView
		cmds = append(cmds, loadTableDataCmd(a.store, msg.tableName, 0))

	case editSubmitMsg:
		cmds = append(cmds, execEditCmd(a.store, msg))

	case errMsg:
		a.err = msg.err
	}

	switch a.focus {
	case listView:
		mod, cmd = a.tableListModel.Update(msg)
		a.tableListModel = mod.(tableList)
		a.tableListModel.setFocus(true)
		cmds = append(cmds, cmd)

	case tableView:
		mod, cmd = a.tableModel.Update(msg)
		a.tableModel = mod.(model)
		a.tableListModel.setFocus(false)
		cmds = append(cmds, cmd)
	}

	return a, tea.Batch(cmds...)
}

func (a App) View() string {
	if !a.ready {
		return "Loading..."
	}

	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		a.tableListModel.View(),
		a.tableModel.View(),
	)

	return lipgloss.JoinVertical(
		lipgloss.Center,
		content,
		a.help.View(keys),
	)
}
