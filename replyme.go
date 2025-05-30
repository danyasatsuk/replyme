package replyme

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Run starts the REPL
func Run(app *App) error {
	err := i18nInit()
	if err != nil {
		return err
	}
	app.setHelpFlags()
	_, err = tea.NewProgram(createModel(app), tea.WithAltScreen(), tea.WithMouseAllMotion()).Run()
	return err
}
