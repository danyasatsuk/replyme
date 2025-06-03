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
	c           chan TUIResponse
	close       chan bool
	width       int
	height      int
}

func (m selectOne) SetParams(p TUISelectOneParams, c chan TUIResponse) selectOne {
	m.params = p

	items := make([]list.Item, len(p.Items))
	for i, item := range p.Items {
		items[i] = item
	}

	m.listModel.SetItems(items)
	m.listModel.Title = p.Name
	m.IsValidated = false
	m.IsExit = false
	m.c = c

	if !m.isCLI {
		m.listModel.SetWidth(m.width - 2)
		m.listModel.SetHeight(m.height - 2)
	}

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
				m.c <- TUIResponse{
					Value: TUISelectOneResult{
						SelectedID:   m.Value.ID,
						SelectedItem: m.Value,
					},
				}

				if m.isCLI {
					return m, tea.Quit
				}

				m.close <- true
			}

			return m, nil
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.listModel.SetWidth(msg.Width - 2)
		m.listModel.SetHeight(msg.Height - 2)
	}

	var cmd tea.Cmd
	m.listModel, cmd = m.listModel.Update(msg)

	return m, cmd
}

func (m selectOne) View() string {
	return inputContainer.Width(m.width - 2).Height(m.height - 2).Render(m.listModel.View())
}

func selectOneNew(c chan bool, isCLI ...bool) selectOne {
	var cli bool
	if len(isCLI) > 0 && isCLI[0] {
		cli = true
	}

	l := list.New([]list.Item{}, list.NewDefaultDelegate(), standardWidth, standardHeight)
	l.SetStatusBarItemName(L(i18n_tui_selectone_item), L(i18n_tui_selectone_items))
	l.SetShowHelp(false)

	m := selectOne{
		listModel: l,
		isCLI:     cli,
		close:     c,
	}

	return m
}
