package models

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

func (m *model) formView() string {
	switch m.form.State {
	case huh.StateCompleted:
		return "Success!"
	default:
		v := strings.TrimSuffix(m.form.View(), "\n\n")
		form := lipgloss.NewStyle().Margin(1, 0).Render(v)

		errors := m.form.Errors()
		header := "Edit Entry"
		if len(errors) > 0 {
			header = fmt.Sprintf("Error: %v", errors[0])
		}
		return header + "\n" + form
	}
}

func (m *model) onEditSuccess(t tab) {
	m.tabs = nil
	m.tabs = []string{"Data", "Info", "Query"}
	m.activeTab = t
	m.selectedRow = nil
	m.toEdit = nil
	m.form = nil
}

func (m *model) onRowSelect(row []string) {
	m.selectedRow = row

	m.tabs = append(m.tabs, "Edit")
	m.activeTab = editTab

	m.toEdit = row

	var inputs []huh.Field
	for i, col := range m.columns {
		inputs = append(
			inputs,
			huh.NewInput().
				Key(col.Name).
				Title(col.Name).
				Description(col.Type).
				Placeholder(row[i]).
				Value(&m.toEdit[i]).
				Validate(func(str string) error {
					switch col.Type {
					case "TEXT":
						return nil
					case "INTEGER":
						if str == "NULL" {
							return nil
						}
						if _, err := strconv.Atoi(str); err != nil {
							return errors.New("Not an integer")
						}
					default:
						return nil
					}
					return nil
				}),
		)
	}
	inputs = append(inputs, huh.NewConfirm().Title("Save").Value(&m.confirmEdit))

	m.form = huh.NewForm(
		huh.NewGroup(inputs...),
	)
}
