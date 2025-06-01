package replyme

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type selectOne struct {
	listModel   list.Model
	IsValidated bool
	params      TUISelectOneParams
	Value       TUISelectItem
	IsExit      bool
	onExit      func(id string)
	isCLI       bool
}

func (m selectOne) SetParams(p TUISelectOneParams) selectOne {
	m.params = p

	items := make([]list.Item, len(p.Items))
	for i, item := range p.Items {
		items[i] = item
	}

	m.listModel.SetItems(items)
	m.listModel.Title = p.Name
	m.IsValidated = false
	m.IsExit = false

	return m
}

func (m selectOne) Init() tea.Cmd {
	return nil
}

func (m selectOne) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.IsExit = true

			if m.isCLI {
				return m, tea.Quit
			}

			return m, nil
		case "enter":
			selected := m.listModel.SelectedItem()
			if item, ok := selected.(TUISelectItem); ok {
				m.Value = item
				m.IsValidated = true

				if m.isCLI {
					return m, tea.Quit
				}
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

func (m selectOne) View() string {
	return m.listModel.View()
}

func selectOneNew(isCLI ...bool) selectOne {
	var cli bool
	if len(isCLI) > 0 && isCLI[0] {
		cli = true
	}

	m := selectOne{
		listModel: list.New([]list.Item{}, list.NewDefaultDelegate(), standardWidth, standardHeight),
		isCLI:     cli,
	}

	return m
}
