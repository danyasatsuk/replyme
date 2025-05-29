package replyme

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Run starts the REPL
func Run(app *App) error {
	err := I18nInit()
	if err != nil {
		return err
	}
	_, err = tea.NewProgram(CreateModel(app), tea.WithAltScreen()).Run()
	return err
}
