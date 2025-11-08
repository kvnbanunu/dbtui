package models

import (
	"dbtui/internal/database"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type App struct {
	store *database.DB
	state State
	err   error
	// tabsModel      tabs      // disply tableInfo or tableData
	// tableInfoModel tableInfo // table metadata
	// tableDataModel tableData // table content
	tableListModel tableList
	width          int
	height         int
	ready          bool
}

func NewApp(db *database.DB) App {

	tl := newTableList()

	return App{
		store: db,
		state: stateTableList,
		tableListModel: tl,
		ready: false,
	}
}

func (a App) Init() tea.Cmd {
	return tea.Batch(
		loadTablesCmd(a.store),
		textinput.Blink,
	)
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.ready = true

		// update all sub models with new dimensions
		a.tableListModel.setSize(msg.Width, msg.Height)

	case tea.KeyMsg:
		switch a.state {
		case stateTableList:
			switch {
			case key.Matches(msg, keys.Quit):
				return a, tea.Quit
			}
		}
	case tablesLoadedMsg:
		a.tableListModel.setTables(msg.tables)

	case errMsg:
		a.err = msg.err
	}

	switch a.state {
	case stateTableList:
		var mod tea.Model
		mod, cmd = a.tableListModel.Update(msg)
		a.tableListModel = mod.(tableList)
		cmds = append(cmds, cmd)
	}

	return a, tea.Batch(cmds...)
}

func (a App) View() string {
	if !a.ready {
		return "Loading..."
	}

	var content string

	switch a.state {
	case stateTableList:
		content = a.tableListModel.View()
	}

	return content + "\n"
}
