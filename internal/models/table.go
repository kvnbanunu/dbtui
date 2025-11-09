package models

import (
	"strings"

	"dbtui/internal/database"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type tab int

const (
	dataTab tab = iota
	infoTab
	queryTab
	editTab
)

type model struct {
	focus       bool
	tabs        []string
	activeTab   tab
	name        string
	columns     []database.Column
	selectedRow []string
	dataTable   table.Model
	queryInput  textinput.Model
	queryTable  table.Model
	queryResult [][]string
	form        *huh.Form
	toEdit      []string
	confirmEdit bool
	currentPage int
	err         error
	width       int
	height      int
}

func newModel() model {
	ti := textinput.New()
	ti.Placeholder = "SELECT * FROM table_name"
	ti.Focus()
	ti.Width = 50

	return model{
		focus:      false,
		tabs:       []string{"Data", "Info", "Query"},
		activeTab:  dataTab,
		dataTable:  newTable(),
		queryInput: ti,
		queryTable: newTable(),
		form:       nil,
	}
}

func (m *model) setSize(w, h int) {
	m.width = w
	m.height = h
	m.dataTable.SetWidth(w - 2)
	m.queryInput.Width = w
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.activeTab {
		case dataTab, infoTab:
			switch {
			case key.Matches(msg, keys.Left):
				m.activeTab = max(m.activeTab-1, dataTab)
			case key.Matches(msg, keys.Right):
				m.activeTab = min(m.activeTab+1, tab(len(m.tabs)-1))
			case key.Matches(msg, keys.Tab):
				m.activeTab = m.nextTab()
			case key.Matches(msg, keys.Edit):
				return m, selectRowCmd(m.dataTable.SelectedRow())
			}
		case queryTab:
			switch {
			case key.Matches(msg, keys.Left):
				m.activeTab = max(m.activeTab-1, dataTab)
			case key.Matches(msg, keys.Right):
				m.activeTab = min(m.activeTab+1, tab(len(m.tabs)-1))
			case key.Matches(msg, keys.Tab):
				m.activeTab = m.nextTab()
			case key.Matches(msg, keys.Reset):
				m.queryInput.Reset()
				return m, nil
			}
		}

	case tableDataLoadedMsg:
		m.setDataTable(msg.tableName, msg.columns, msg.rows)

	case rowSelectedMsg:
		m.onRowSelect(msg.row)
		cmds = append(cmds, m.form.Init())

	case queryResultMsg:
		m.setQueryResult(msg.columns, msg.rows, msg.err)
	}

	switch m.activeTab {
	case dataTab, infoTab:
		m.dataTable, cmd = m.dataTable.Update(msg)
		cmds = append(cmds, cmd)
	case queryTab:
		m.queryInput, cmd = m.queryInput.Update(msg)
		cmds = append(cmds, cmd)
	case editTab:
		form, cmd := m.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form = f
			cmds = append(cmds, cmd)
		}

		if m.form.State == huh.StateCompleted {
			// m.activeTab = dataTab
			cmds = append(cmds, editSubmitCmd(m.name, m.selectedRow[0], m.columns, m.toEdit))
			m.onEditSuccess(dataTab)
		}
	}
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	doc := strings.Builder{}

	var renderedTabs []string

	tabWidth := 13

	for i, tab := range m.tabs {
		var style lipgloss.Style
		isFirst, isActive := i == 0, i == int(m.activeTab)
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
		}
		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(tab))
		tabWidth += style.GetHorizontalFrameSize()
	}

	remTab := remainderTabStyle
	renderedTabs = append(renderedTabs, remTab.Width((m.width - tabWidth)).Render(""))

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")

	var selected string
	switch m.activeTab {
	case dataTab, infoTab:
		selected = m.dataView()
	case queryTab:
		selected = m.queryView()
	case editTab:
		selected = m.formView()
	}

	doc.WriteString(windowStyle.
		Width((m.width)).
		Height((m.dataTable.Height())).
		Render(selected))
	return docStyle.Render(doc.String())
}

func (m model) nextTab() tab {
	switch m.activeTab {
	case dataTab:
		return infoTab
	case infoTab:
		return queryTab
	case queryTab:
		return dataTab
	default:
		return dataTab
	}
}
