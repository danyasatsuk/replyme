package replyme

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type inputText struct {
	input       textinput.Model
	IsValidated bool
	IsExit      bool
	Value       string
	params      TUIInputTextParams
	isCLI       bool
	close       chan bool
	c           chan TUIResponse
}

func (m inputText) SetParams(p TUIInputTextParams, c chan TUIResponse) inputText {
	m.params = p
	if m.params.IsPassword {
		m.input.EchoMode = textinput.EchoPassword
		m.input.EchoCharacter = '*'
	}

	m.input.Placeholder = m.params.Placeholder

	m.c = c

	m = m.Focus()

	return m
}

func (m inputText) Focus() inputText {
	m.input.Focus()

	return m
}

func (m inputText) Blur() {
	m.input.Blur()
}

func (m inputText) Init() tea.Cmd {
	return nil
}

func (m inputText) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.IsExit = true

			if m.isCLI {
				return m, tea.Quit
			}

			return m, nil
		case "enter":
			if m.params.MaxLength > 0 && len(m.input.Value()) > m.params.MaxLength {
				return m, nil
			}

			if m.params.Validate == nil || m.params.Validate(m.input.Value()) {
				m.IsValidated = true
				m.Value = m.input.Value()
				m.input.Reset()

				m.c <- TUIResponse{
					Value: m.Value,
					Err:   nil,
				}
				m.close <- true

				if m.isCLI {
					return m, tea.Quit
				}

				return m, nil
			}
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)

	return m, cmd
}

func (m inputText) View() string {
	return fmt.Sprintf(`%s

%s

%s`, m.params.Name, m.params.Description, m.input.View())
}

func inputTextNew(c chan bool, isCLI ...bool) inputText {
	var cli bool
	if len(isCLI) > 0 && isCLI[0] {
		cli = true
	}

	t := textinput.New()
	t.Width = standardWidth

	m := inputText{
		input: t,
		isCLI: cli,
		close: c,
	}

	return m
}
