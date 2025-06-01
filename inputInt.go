package replyme

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type inputInt struct {
	input       textinput.Model
	IsValidated bool
	IsExit      bool
	Value       int
	params      TUIInputIntParams
	isCLI       bool
}

func inputIntNew(isCLI ...bool) inputInt {
	var cli bool
	if len(isCLI) > 0 && isCLI[0] {
		cli = true
	}

	return inputInt{
		input: textinput.New(),
		isCLI: cli,
	}
}

func (m inputInt) SetParams(p TUIInputIntParams) {
	m.params = p
	m.input.Placeholder = L(i18n_inputint_placeholder)
	m.input.Focus()
}

func (m inputInt) Init() tea.Cmd {
	return nil
}

func (m inputInt) Focus() {
	m.input.Focus()
}

func (m inputInt) Blur() {
	m.input.Blur()
}

func (m inputInt) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			return m.onEnter()
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)

	return m, cmd
}

func (m inputInt) View() string {
	return fmt.Sprintf(`%s

%s

%s`, m.params.Name, m.params.Description, m.input.View())
}

func (m inputInt) onEnter() (inputInt, tea.Cmd) {
	strVal := m.input.Value()

	intVal, err := strconv.Atoi(strVal)
	if err != nil {
		return m, nil
	}

	if m.params.MinValue != 0 && intVal < m.params.MinValue {
		return m, nil
	}

	if m.params.MaxValue != 0 && intVal > m.params.MaxValue {
		return m, nil
	}

	if m.params.Validate != nil && !m.params.Validate(strVal) {
		return m, nil
	}

	m.IsValidated = true
	m.Value = intVal
	m.input.Reset()

	if m.isCLI {
		return m, tea.Quit
	}

	return m, nil
}
