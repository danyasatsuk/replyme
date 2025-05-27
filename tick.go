package replyme

import (
	tea "github.com/charmbracelet/bubbletea"
	"time"
)

type Tick time.Time

func Ticker() tea.Cmd {
	return tea.Tick(time.Millisecond*50, func(t time.Time) tea.Msg {
		return Tick(t)
	})
}
