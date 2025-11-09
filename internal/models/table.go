package models

import (
	"fmt"
	"strings"

	"dbtui/internal/database"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type tab int

const (
	dataTab tab = iota
	infoTab
	queryTab
)

var (
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	remainderTabBorder = tabRemainderBorder()
	docStyle          = lipgloss.NewStyle()
	highlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Border(activeTabBorder, true)
	remainderTabStyle = inactiveTabStyle.Border(remainderTabBorder, true)
	windowStyle       = lipgloss.NewStyle().
		BorderForeground(highlightColor).
		Padding(1, 0).
		Align(lipgloss.Center).
		Border(lipgloss.NormalBorder()).
		UnsetBorderTop()
)

type model struct {
	store       *database.Manager
	tabs        []string
	activeTab   tab
	name        string
	columns     []database.Column
	dataTable   table.Model
	queryInput  textinput.Model
	queryTable  table.Model
	queryResult [][]string
	currentPage int
	err         error
	width       int
	height      int
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Left):
			m.activeTab = max(m.activeTab-1, dataTab)
		case key.Matches(msg, keys.Right):
			m.activeTab = min(m.activeTab+1, tab(len(m.tabs)-1))
		}
		switch m.activeTab {
		case dataTab, infoTab:
			switch {
			case key.Matches(msg, keys.Edit):
				// TODO selectedRow
				return m, nil
			}
		case queryTab:
			if msg.String() == "ctrl+u" {
				m.queryInput.Reset()
				return m, nil
			}
		}
	case tableDataLoadedMsg:
		m.setDataTable(msg)

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
	}
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	doc := strings.Builder{}

	var renderedTabs []string
	
	tabWidth := 0

	for i, tab := range m.tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(m.tabs)-1, i == int(m.activeTab)
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
			// border.BottomRight = "│"
		} else if isLast && !isActive {
			// border.BottomRight = "┤"
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
	}

	doc.WriteString(windowStyle.
		Width((m.width)).
		Height((m.dataTable.Height())).
		Render(selected))
	return docStyle.Render(doc.String())
}

func (m *model) dataView() string {
	if m.name == "" {
		return "No table selected"
	}

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("170")).
		Render(fmt.Sprintf("Table: %s (Page %d)", m.name, m.currentPage+1))

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		baseStyle.Render(m.dataTable.View()),
	)
}

func (m *model) queryView() string {
	view := "SQL Query\n\n"
	view += m.queryInput.View() + "\n\n"

	if m.err != nil {
		view += fmt.Sprintf("Error: %s\n", m.err)
	} else if len(m.queryResult) > 0 {
		view += m.queryTable.View()
	}

	return view
}

func (m *model) setSize(w, h int) {
	m.width = w
	m.height = h
	m.dataTable.SetWidth(w - 2)
	m.queryInput.Width = w
}

func (m *model) setDataTable(msg tableDataLoadedMsg) {
	m.name = msg.tableName
	m.currentPage = 0
	m.columns = msg.columns

	// reset rows before changing cols
	m.dataTable.SetRows([]table.Row{})

	// convert cols to bubbles
	tableCols := make([]table.Column, len(msg.columns))
	for i, col := range msg.columns {
		width := 10
		if len(col.Name) > width {
			width = len(col.Name) + 2
		}

		tableCols[i] = table.Column{
			Title: col.Name,
			Width: width,
		}
	}

	m.dataTable.SetColumns(tableCols)

	// convert rows to bubbles
	tableRows := make([]table.Row, len(msg.rows))
	for i, row := range msg.rows {
		tableRows[i] = row
	}

	m.dataTable.SetRows(tableRows)

	m.dataTable.GotoTop()

	m.dataTable.Focus()
}

func (m *model) setQueryResult(columns []string, rows [][]string, err error) {
	m.err = err
	m.queryResult = rows

	if err == nil && len(rows) > 0 {
		tableCols := make([]table.Column, len(columns))
		for i, col := range columns {
			width := 15
			if len(col) > width {
				width = len(col) + 2
			}

			tableCols[i] = table.Column{
				Title: col,
				Width: width,
			}
		}

		// convert rows to bubbles
		tableRows := make([]table.Row, len(rows))
		for i, row := range rows {
			tableRows[i] = row
		}

		m.queryTable.SetColumns(tableCols)
		m.queryTable.SetRows(tableRows)

	}
}

func newModel(m *database.Manager) model {
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

	ti := textinput.New()
	ti.Placeholder = "SELECT * FROM table_name"
	ti.Focus()
	ti.Width = 50

	return model{
		store:      m,
		tabs:       []string{"Data", "Info", "Query"},
		activeTab:  dataTab,
		dataTable:  t,
		queryInput: ti,
		queryTable: t,
	}
}

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
