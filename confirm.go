package replyme

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type confirm struct {
	IsValidated bool
	IsExit      bool
	Value       bool // true = Yes, false = No

	params TUIConfirmParams
	cursor int // 0 = Yes, 1 = No
	c      chan TUIResponse
	close  chan bool
}

func confirmNew(c chan bool) *confirm {
	return &confirm{
		close: c,
	}
}

func (m *confirm) SetParams(p TUIConfirmParams, c chan TUIResponse) {
	m.params = p
	m.cursor = 0
	m.c = c
}

func (m *confirm) Init() tea.Cmd {
	return nil
}

func (m *confirm) Update(msg tea.Msg) (*confirm, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.IsExit = true

			return m, nil

		case "left", "h":
			m.cursor = 0

			return m, nil

		case "right", "l":
			m.cursor = 1

			return m, nil

		case "y":
			m.IsValidated = true
			m.Value = true

			return m, nil

		case "n":
			m.IsValidated = true
			m.Value = false

			return m, nil

		case "enter":
			m.IsValidated = true
			m.Value = m.cursor == 0
			m.c <- TUIResponse{m.Value, nil}
			m.close <- true

			return m, nil
		}
	}

	return m, nil
}

func (m *confirm) View() string {
	var yes, no string

	if m.cursor == 0 {
		yes = fmt.Sprintf("[> %s <]", L(i18n_confirm_view_yes))
		no = fmt.Sprintf("[  %s  ]", L(i18n_confirm_view_no))
	} else {
		yes = fmt.Sprintf("[  %s  ]", L(i18n_confirm_view_yes))
		no = fmt.Sprintf("[> %s <]", L(i18n_confirm_view_no))
	}

	return fmt.Sprintf(`%s

%s

%s %s`, m.params.Name, m.params.Description, yes, no)
}
