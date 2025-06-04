package replyme

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"strings"
)

func cliRunSelectOne(t TUIRequest, c chan<- bool) {
	m := selectOneNew(make(chan bool), true)
	m = m.SetParams(t.Payload.(TUISelectOneParams), t.Response)
	_, err := tea.NewProgram(m, tea.WithAltScreen()).Run()

	if err != nil {
		t.Response <- TUIResponse{
			Err: err,
		}

		return
	}

	c <- true
}

func cliRunInputText(t TUIRequest, c chan<- bool) {
	m := inputTextNew(make(chan bool), true)
	m = m.SetParams(t.Payload.(TUIInputTextParams), t.Response)
	m = m.Focus()
	_, err := tea.NewProgram(m, tea.WithAltScreen()).Run()

	if err != nil {
		t.Response <- TUIResponse{
			Err: err,
		}

		return
	}

	c <- true
}

func cliRunInputInt(t TUIRequest, c chan<- bool) {
	m := inputIntNew(make(chan bool), true)
	m = m.SetParams(t.Payload.(TUIInputIntParams), t.Response)
	_, err := tea.NewProgram(m, tea.WithAltScreen()).Run()

	if err != nil {
		t.Response <- TUIResponse{
			Err: err,
		}

		return
	}

	c <- true
}

func cliRunInputFile(t TUIRequest, c chan<- bool) {
	m := inputFileNew(make(chan bool), true)
	m = m.SetParams(t.Payload.(TUIInputFileParams), t.Response)
	_, err := tea.NewProgram(m, tea.WithAltScreen()).Run()

	if err != nil {
		t.Response <- TUIResponse{
			Err: err,
		}

		return
	}

	c <- true
}

func cliRunConfirm(t TUIRequest, c chan<- bool) {
	m := confirmNew(make(chan bool), true)
	m = m.SetParams(t.Payload.(TUIConfirmParams), t.Response)
	_, err := tea.NewProgram(m, tea.WithAltScreen()).Run()

	if err != nil {
		t.Response <- TUIResponse{
			Err: err,
		}

		return
	}

	c <- true
}

func runCLITUI(t TUIRequest, c chan<- bool) {
	switch t.Type {
	case tuiTypeSelectOne:
		cliRunSelectOne(t, c)
	case tuiTypeInputText:
		cliRunInputText(t, c)
	case tuiTypeInputInt:
		cliRunInputInt(t, c)
	case tuiTypeInputFile:
		cliRunInputFile(t, c)
	case tuiTypeConfirm:
		cliRunConfirm(t, c)
	}
}

func cliRunner(app *App) error {
	cmd := strings.Join(os.Args[1:], " ")
	logsChan := make(chan log)

	go func() {
		for d := range logsChan {
			l := log{
				logTypeLog,
				cmd,
				d.Message,
				d.Error,
				d.Time,
			}
			fmt.Println(l.Render())
		}
	}()

	if cmd == "" {
		h, err := helpApp(app)

		if err != nil {
			return err
		}

		fmt.Println(h)

		return nil
	}

	err := fullRunCommand(fullRunCommandParams{
		cmd, app, logsChan, os.Stdout, os.Stderr, func(msg logMsg) {
			logsChan <- log{
				logTypeLog,
				cmd,
				msg.Content,
				nil,
				msg.Time,
			}
		}, nil, runCLITUI, true,
	})

	return err
}
