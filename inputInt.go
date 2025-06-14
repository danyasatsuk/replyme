package replyme

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
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
	c           chan TUIResponse
	close       chan bool
	width       int
	height      int
}

func inputIntNew(c chan bool, isCLI ...bool) inputInt {
	var cli bool
	if len(isCLI) > 0 && isCLI[0] {
		cli = true
	}

	t := textinput.New()
	t.Width = standardWidth
	t.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))

	return inputInt{
		input: t,
		isCLI: cli,
		close: c,
	}
}

func (m inputInt) SetParams(p TUIInputIntParams, c chan TUIResponse) inputInt {
	m.params = p
	m.input.Placeholder = L(i18n_inputint_placeholder)
	m.input.Focus()
	m.c = c

	return m
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
	case tea.WindowSizeMsg:
		var cmd tea.Cmd

		m.width = msg.Width
		m.height = msg.Height
		m.input, cmd = m.input.Update(msg)

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

		case "enter":
			return m.onEnter()
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)

	return m, cmd
}

func (m inputInt) View() string {
	return inputContainer.Width(m.width - 2).Height(m.height - 2).Render(fmt.Sprintf(`%s

%s
%s

%s`, styles.InputTitle(m.params.Name), m.input.View(),
		styles.GrayStyle(fmt.Sprintf("<=%d | >=%d", m.params.MinValue, m.params.MaxValue)),
		styles.InputDescription(m.params.Description)))
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

	m.c <- TUIResponse{
		Value: m.Value,
		Err:   nil,
	}

	if m.isCLI {
		return m, tea.Quit
	}

	m.close <- true

	return m, nil
}
