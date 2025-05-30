package replyme

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type selectOne struct {
	listModel   list.Model
	IsValidated bool
	params      TUISelectOneParams
	Value       tuiSelectItem
	IsExit      bool
}

func (m *selectOne) SetParams(p TUISelectOneParams) {
	m.params = p

	items := make([]list.Item, len(p.Items))
	for i, item := range p.Items {
		items[i] = item
	}
	m.listModel.SetItems(items)
	m.listModel.Title = p.Name
	m.IsValidated = false
	m.IsExit = false
}

func (m *selectOne) Init() tea.Cmd {
	return nil
}

func (m *selectOne) Update(msg tea.Msg) (*selectOne, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.IsExit = true
			return m, nil
		case "enter":
			selected := m.listModel.SelectedItem()
			if item, ok := selected.(tuiSelectItem); ok {
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

func (m *selectOne) View() string {
	return m.listModel.View()
}

func selectOneNew() *selectOne {
	m := &selectOne{
		listModel: list.New([]list.Item{}, list.NewDefaultDelegate(), 70, 10),
	}
	return m
}
