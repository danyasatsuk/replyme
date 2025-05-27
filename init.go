package replyme

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Init is a method for initializing the BubbleTea model.
func (m *Model) Init() tea.Cmd {
	if m.app.Params.EnableInputBlinking {
		return tea.Batch(m.spinner.Tick, Ticker(), textinput.Blink)
	}
	return tea.Batch(m.spinner.Tick, Ticker())
}
