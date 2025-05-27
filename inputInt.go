package replyme

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type InputInt struct {
	input       textinput.Model
	IsValidated bool
	IsExit      bool
	Value       int
	params      TUIInputIntParams
}

func InputIntNew() *InputInt {
	return &InputInt{
		input: textinput.New(),
	}
}

func (m *InputInt) SetParams(p TUIInputIntParams) {
	m.params = p
	m.input.Placeholder = "Введите число"
	m.input.Focus()
}

func (m *InputInt) Init() tea.Cmd {
	return nil
}

func (m *InputInt) Focus() {
	m.input.Focus()
}

func (m *InputInt) Blur() {
	m.input.Blur()
}

func (m *InputInt) Update(msg tea.Msg) (*InputInt, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.IsExit = true
			return m, nil

		case "enter":
			strVal := m.input.Value()
			intVal, err := strconv.Atoi(strVal)
			if err != nil {
				return m, nil
			}

			// Проверка на диапазон
			if m.params.MinValue != 0 && intVal < m.params.MinValue {
				return m, nil
			}
			if m.params.MaxValue != 0 && intVal > m.params.MaxValue {
				return m, nil
			}

			// Валидация функцией
			if m.params.Validate != nil && !m.params.Validate(strVal) {
				return m, nil
			}

			m.IsValidated = true
			m.Value = intVal
			m.input.Reset()
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m *InputInt) View() string {
	return fmt.Sprintf(`%s

%s

%s`, m.params.Name, m.params.Description, m.input.View())
}
