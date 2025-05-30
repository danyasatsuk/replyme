package replyme

import (
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type InputResizeMsg struct {
	Delta int
}

type TerminalInput struct {
	text          string
	lines         []string
	cursor        int
	width         int
	history       []string
	historyIx     int
	lastLineCount int
	running       bool

	viewport viewport.Model
}

func NewTerminalInput() TerminalInput {
	vp := viewport.New(70, 1)
	vp.SetContent("")
	return TerminalInput{
		cursor:    0,
		width:     70,
		viewport:  vp,
		historyIx: 0,
	}
}

func (m TerminalInput) Init() tea.Cmd {
	return nil
}

func (m TerminalInput) Update(msg tea.Msg) (TerminalInput, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyRunes, tea.KeySpace:
			r := msg.String()
			runes := []rune(m.text)
			m.text = string(runes[:m.cursor]) + r + string(runes[m.cursor:])
			m.cursor += utf8.RuneCountInString(r)
		case tea.KeyBackspace:
			if m.cursor > 0 {
				runes := []rune(m.text)
				m.text = string(runes[:m.cursor-1]) + string(runes[m.cursor:])
				m.cursor--
			}
		case tea.KeyLeft:
			if m.cursor > 0 {
				m.cursor--
			}
		case tea.KeyRight:
			if m.cursor < len([]rune(m.text)) {
				m.cursor++
			}
		case tea.KeyShiftUp:
			if len(m.history) == 0 {
				break
			}
			if m.historyIx > 0 {
				m.historyIx--
				m.text = m.history[m.historyIx]
				m.cursor = len([]rune(m.text))
			}

		case tea.KeyShiftDown:
			if len(m.history) == 0 {
				break
			}
			if m.historyIx < len(m.history)-1 {
				m.historyIx++
				m.text = m.history[m.historyIx]
				m.cursor = len([]rune(m.text))
			} else if m.historyIx == len(m.history)-1 {
				m.historyIx = len(m.history)
				m.text = ""
				m.cursor = 0
			}
		case tea.KeyEnter:
			if trimmed := strings.TrimSpace(m.text); trimmed != "" {
				m.history = append(m.history, trimmed)
			}
			m.historyIx = len(m.history)
			m.text = ""
			m.cursor = 0
		}

		m.recalculateLines()
		m.viewport.SetContent(m.render())
		m.viewport.Height = m.GetLines()

		if len(m.lines) > 10 {
			m.viewport.YOffset = len(m.lines) - 10
		} else {
			m.viewport.YOffset = 0
		}

		if len(m.lines) != m.lastLineCount {
			delta := len(m.lines) - m.lastLineCount
			m.lastLineCount = len(m.lines)
			return m, func() tea.Msg {
				return InputResizeMsg{Delta: delta}
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.viewport.Width = msg.Width
		m.recalculateLines()
		m.viewport.SetContent(m.render())

		if len(m.lines) > 10 {
			m.viewport.YOffset = len(m.lines) - 10
		} else {
			m.viewport.YOffset = 0
		}

		if len(m.lines) != m.lastLineCount {
			delta := len(m.lines) - m.lastLineCount
			m.lastLineCount = len(m.lines)
			return m, func() tea.Msg {
				return InputResizeMsg{Delta: delta}
			}
		}
	}
	return m, nil
}

func (m TerminalInput) View() string {
	if m.running {
		m.viewport.SetContent(GrayStyle(L(i18n_cmd_input_running)))
	}
	return m.viewport.View()
}

func (m TerminalInput) Value() string {
	return m.text
}

func (m TerminalInput) GetLines() int {
	if len(m.lines) > 10 {
		return 10
	}
	return len(m.lines)
}

func (m *TerminalInput) recalculateLines() {
	if strings.TrimSpace(m.text) == "" {
		m.lines = []string{GrayStyle("> " + L(i18n_cmd_input_command))}
	} else {
		wrapped := wrapLines(m.text, m.width-2)
		for i := range wrapped {
			wrapped[i] = GrayStyle("> ") + wrapped[i]
		}
		m.lines = wrapped
	}
}

func (m TerminalInput) render() string {
	cursorChar := "▌"
	runes := []rune(m.text)

	if m.cursor < 0 {
		m.cursor = 0
	} else if m.cursor > len(runes) {
		m.cursor = len(runes)
	}

	marked := string(runes[:m.cursor]) + cursorChar + string(runes[m.cursor:])

	if strings.TrimSpace(m.text) == "" {
		return GrayStyle("> " + L(i18n_cmd_input_command))
	}

	wrapped := wrapLines(marked, m.width-2)
	for i := range wrapped {
		wrapped[i] = GrayStyle("> ") + wrapped[i]
	}
	return strings.Join(wrapped, "\n")
}

func wrapLines(text string, width int) []string {
	if width <= 0 {
		return []string{text}
	}
	var lines []string
	var current strings.Builder
	count := 0
	for _, r := range text {
		current.WriteRune(r)
		count++
		if count >= width {
			lines = append(lines, current.String())
			current.Reset()
			count = 0
		}
	}
	if current.Len() > 0 {
		lines = append(lines, current.String())
	}
	return lines
}
