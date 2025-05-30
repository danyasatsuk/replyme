package replyme

import (
	tea "github.com/charmbracelet/bubbletea"
	"time"
)

type tick time.Time

func ticker() tea.Cmd {
	return tea.Tick(time.Millisecond*50, func(t time.Time) tea.Msg {
		return tick(t)
	})
}
