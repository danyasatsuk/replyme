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
	isCLI  bool
	width  int
	height int
}

func confirmNew(c chan bool, isCLI ...bool) confirm {
	var cli bool
	if len(isCLI) > 0 && isCLI[0] {
		cli = true
	}

	return confirm{
		close: c,
		isCLI: cli,
	}
}

func (m confirm) SetParams(p TUIConfirmParams, c chan TUIResponse) confirm {
	m.params = p
	m.cursor = 0
	m.c = c

	return m
}

func (m confirm) Init() tea.Cmd {
	return nil
}

func (m confirm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		var cmd tea.Cmd

		m.width = msg.Width
		m.height = msg.Height

		return m, cmd
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.IsExit = true

			if m.isCLI {
				return m, tea.Quit
			}

			m.close <- true

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

			if m.isCLI {
				return m, tea.Quit
			}

			m.close <- true

			return m, nil
		}
	}

	return m, nil
}

func (m confirm) View() string {
	var yes, no string

	if m.cursor == 0 {
		yes = styles.InputSelected(fmt.Sprintf("[> %s <]", L(i18n_confirm_view_yes)))
		no = fmt.Sprintf("[  %s  ]", L(i18n_confirm_view_no))
	} else {
		yes = fmt.Sprintf("[  %s  ]", L(i18n_confirm_view_yes))
		no = styles.InputSelected(fmt.Sprintf("[> %s <]", L(i18n_confirm_view_no)))
	}

	return inputContainer.Width(m.width - 2).Height(m.height - 2).Render(fmt.Sprintf(`%s

%s

%s %s`, styles.InputTitle(m.params.Name), styles.InputDescription(m.params.Description), yes, no))
}
