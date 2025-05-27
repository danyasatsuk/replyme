package replyme

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"log"
)

type SelectOne struct {
	listModel   list.Model
	IsValidated bool
	params      TUISelectOneParams
	Value       TUISelectItem
	IsExit      bool
}

func (m *SelectOne) SetParams(p TUISelectOneParams) {
	m.params = p

	log.Printf("%+v", p)

	items := make([]list.Item, len(p.Items))
	for i, item := range p.Items {
		items[i] = item
	}
	m.listModel.SetItems(items)
	m.listModel.Title = p.Name
	m.IsValidated = false
	m.IsExit = false
}

func (m *SelectOne) Init() tea.Cmd {
	return nil
}

func (m *SelectOne) Update(msg tea.Msg) (*SelectOne, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.IsExit = true
			return m, nil
		case "enter":
			selected := m.listModel.SelectedItem()
			if item, ok := selected.(TUISelectItem); ok {
				m.Value = item
				m.IsValidated = true
			}
			return m, nil
		}
	case tea.WindowSizeMsg:
		m.listModel.SetWidth(msg.Width)
		m.listModel.SetHeight(msg.Height)
	}
	var cmd tea.Cmd
	m.listModel, cmd = m.listModel.Update(msg)
	return m, cmd
}

func (m *SelectOne) View() string {
	return m.listModel.View()
}

func SelectOneNew() *SelectOne {
	m := &SelectOne{
		listModel: list.New([]list.Item{}, list.NewDefaultDelegate(), 70, 10),
	}
	return m
}
