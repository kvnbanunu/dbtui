package models

import (
	"dbtui/internal/database"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type App struct {
	store          *database.Manager
	state          State // focused
	help           help.Model
	err            error
	tableListModel tableList
	tableDataModel tableData // table content
	queryModel     queryModel
	tabsModel      tabsModel // disply tableInfo or tableData or query
	width          int
	height         int
	ready          bool
	// tableInfoModel tableInfo // table metadata
}

func NewApp(m *database.Manager) App {
	tl := newTableList()
	td := newTableData()
	q := newQueryModel()
	tabs := newTabsModel()
	help := help.New()
	help.ShowAll = false

	return App{
		store:          m,
		state:          stateTableList,
		help:           help,
		tableListModel: tl,
		tableDataModel: td,
		queryModel:     q,
		tabsModel:      tabs,
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
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.ready = true

		helpHeight := lipgloss.Height(a.help.View(keys))
		tabsHeight := lipgloss.Height(a.tabsModel.View()) + 2
		contentHeight := msg.Height - helpHeight - 2

		listWidth := msg.Width * 30 / 100
		contentWidth := msg.Width - listWidth - 2

		// update all sub models with new dimensions
		a.tableListModel.setSize(listWidth, contentHeight)
		a.tableDataModel.setSize(contentWidth, contentHeight - tabsHeight)
		a.queryModel.setSize(contentWidth, contentHeight - tabsHeight)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			return a, tea.Quit

		case key.Matches(msg, keys.Tab):
			a.state = a.nextState()
			return a, nil

		case key.Matches(msg, keys.Back):
			if a.state != stateTableList {
				a.state = stateTableList
			}
			return a, nil
		case key.Matches(msg, keys.Help):
			a.help.ShowAll = !a.help.ShowAll
			return a, nil
		}
	case tablesLoadedMsg:
		a.tableListModel.setTables(msg.tables)

	case tableSelectedMsg:
		a.state = stateTableData
		a.tableDataModel.setTable(msg.tableName)
		cmds = append(cmds, loadTableDataCmd(a.store, msg.tableName, 0))

	case switchTabMsg:
		a.state = State(msg.activeTab)

	case tableDataLoadedMsg:
		a.tableDataModel.setData(msg.columns, msg.rows)

	case queryResultMsg:
		a.queryModel.setResults(msg.columns, msg.rows, msg.err)

	case errMsg:
		a.err = msg.err
	}


	var model tea.Model
	var cmd tea.Cmd
	switch a.state {
	case stateTableList:
		model, cmd = a.tableListModel.Update(msg)
		a.tableListModel = model.(tableList)
		cmds = append(cmds, cmd)

	case stateTableData, stateTableInfo:
		model, cmd = a.tabsModel.Update(msg)
		a.tabsModel = model.(tabsModel)
		cmds = append(cmds, cmd)
		model, cmd = a.tableDataModel.Update(msg)
		a.tableDataModel = model.(tableData)
		cmds = append(cmds, cmd)

	case stateQuery:
		model, cmd = a.tabsModel.Update(msg)
		a.tabsModel = model.(tabsModel)
		cmds = append(cmds, cmd)
		model, cmd = a.queryModel.Update(msg)
		a.queryModel = model.(queryModel)
		cmds = append(cmds, cmd)

	// case stateTableInfo:
	// 	model, cmd = a.tabsModel.Update(msg)
	// 	a.tabsModel = model.(tabsModel)
	// 	cmds = append(cmds, cmd)
	}

	return a, tea.Batch(cmds...)
}

func (a App) View() string {
	if !a.ready {
		return "Loading..."
	}

	var content string
	var selected string

	switch a.state {
	case stateTableData:
		selected = a.tableDataModel.View()
	case stateTableInfo:
		selected = a.tableDataModel.View()
	case stateQuery:
		selected = a.queryModel.View()
	}

	content = lipgloss.JoinHorizontal(
		lipgloss.Top,
		a.tableListModel.View(),
		lipgloss.NewStyle().
			Padding(0, 1).
			Render(lipgloss.JoinVertical(
				lipgloss.Left,
				a.tabsModel.View(),
				selected,
			)),
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		content,
		a.help.View(keys),
	)
}

func (a App) nextState() State {
	switch a.state {
	case stateTableList:
		return stateTableData
	case stateTableData:
		return stateTableData
	default:
		return stateTableList
	}
}
