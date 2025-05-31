package replyme

import (
	tea "github.com/charmbracelet/bubbletea"
	"time"
)

type tick time.Time

const tickTime = 50

func ticker() tea.Cmd {
	return tea.Tick(time.Millisecond*tickTime, func(t time.Time) tea.Msg {
		return tick(t)
	})
}
