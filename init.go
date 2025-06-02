package replyme

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *model) Init() tea.Cmd {
	if m.app.Params.EnableInputBlinking {
		return tea.Batch(m.spinner.Tick, ticker(), textinput.Blink)
	}

	return tea.Batch(m.spinner.Tick, ticker(), m.inputFile.Init())
}
