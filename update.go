package replyme

import (
	"errors"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/wordwrap"
)

const scrollLines = 3

func (m *model) tuiUpdater(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.runningTUI.Type {
	case tuiTypeSelectOne:
		var mod tea.Model
		mod, cmd = m.selectOne.Update(msg)
		m.selectOne = mod.(selectOne)
	case tuiTypeSelectSeveral:
		//nolint:godox
		// TODO(medium): Implement SelectSeveral and add
	case tuiTypeInputText:
		var mod tea.Model
		mod, cmd = m.inputText.Update(msg)
		m.selectOne = mod.(selectOne)
	case tuiTypeInputInt:
		var mod tea.Model
		mod, cmd = m.inputInt.Update(msg)
		m.selectOne = mod.(selectOne)
	case tuiTypeInputFile:
		var mod tea.Model
		mod, cmd = m.inputFile.Update(msg)
		m.selectOne = mod.(selectOne)
	case tuiTypeConfirm:
		var mod tea.Model
		mod, cmd = m.confirm.Update(msg)
		m.selectOne = mod.(selectOne)
	}

	return m, tea.Batch(cmd, ticker())
}

func (m *model) updateLogsHeight() {
	if m.isRunningTUI {
		m.logsViewport.Height = 0
		m.tuiViewport.Height = m.windowHeight
	} else {
		m.tuiViewport.Height = 0
		m.logsViewport.Height = m.windowHeight - m.input.GetLines()
	}
}

func (m *model) handleWindowSizeMsg(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	m.windowHeight = msg.Height
	m.windowWidth = msg.Width
	m.updateLogsHeight()
	m.input, _ = m.input.Update(msg)
	m.logsViewport.Width = msg.Width
	m.tuiViewport.Width = msg.Width
	m.logsViewport.GotoBottom()
	cmd := m.updateViewport(msg)

	return m, cmd
}

func (m *model) helpFunc(msg tea.Msg) (tea.Model, tea.Cmd) {
	help, err := helpApp(m.app)
	if err != nil {
		m.logs.Add(logTypeError, fmt.Sprintf("error: %s", err.Error()))
	}

	m.logs.Add(logTypeCommandSuccess, "help")
	m.logs.Add(logTypeMessage, help)
	m.logsViewport.SetContent(m.logs.Render())
	m.input.text = ""

	var cmd tea.Cmd

	var cmd2 tea.Cmd

	m.logsViewport, cmd = m.logsViewport.Update(msg)
	m.input, cmd2 = m.input.Update(msg)

	return m, tea.Batch(cmd, cmd2)
}

func (m *model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "enter":
		command := m.input.Value()
		if command == "" {
			return m, nil
		}

		if command == "exit" {
			return m, tea.Quit
		}

		if command == "help" {
			return m.helpFunc(msg)
		}

		m.logs.Add(logTypeCommandRunning, command)
		m.logsViewport.SetContent(wordwrap.String(m.logs.Render(), m.logsViewport.Width))
		m.logsViewport, _ = m.logsViewport.Update(msg)
		m.runningCommand = command
		m.input.running = true
		m.input, _ = m.input.Update(msg)

		go func() {
			err := m.runCommand(command)
			if err != nil {
				typeOfError := logTypeError

				switch {
				case errors.Is(err, ErrorUnknownCommand):
					typeOfError = logTypeCommandNotFound
				case errors.Is(err, ErrorCommandPanic):
					typeOfError = logTypePanic
				}
				m.logsChan <- log{typeOfError, command, err.Error(), err, time.Now()}
				m.logsChan <- log{logTypeCommandFailure, command, err.Error(), err, time.Now()}
			}
		}()

		return m, tea.Batch(ticker())
	}

	if m.isRunningTUI {
		return m.tuiUpdater(msg)
	}

	m.input, _ = m.input.Update(msg)

	return m, nil
}

func (m *model) handleMouseMsg(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.Button {
	case tea.MouseButtonWheelUp:
		m.logsViewport.ScrollUp(scrollLines)
	case tea.MouseButtonWheelDown:
		m.logsViewport.ScrollDown(scrollLines)
	default:
	}

	m.logsViewport, cmd = m.logsViewport.Update(msg)

	return m, cmd
}

func (m *model) onLogsChan(l log, msg tea.Msg) (tea.Model, tea.Cmd) {
	l.Message = wordwrap.String(l.Message, m.logsViewport.Width)
	m.logs.AddLog(l)
	m.updateLogsHeight()
	m.logsViewport.SetContent(m.logs.Render())
	m.logsViewport.GotoBottom()
	m.logsViewport, _ = m.logsViewport.Update(msg)

	if l.Type == logTypeCommandSuccess || l.Type == logTypeCommandFailure || l.Type == logTypeCommandNotFound {
		m.input.running = false
		m.input, _ = m.input.Update(msg)

		return m, tea.Batch(ticker())
	}

	return m, ticker()
}

func (m *model) onTUIChan(t TUIRequest, msg tea.Msg) (tea.Model, tea.Cmd) {
	m.isRunningTUI = true
	m.runningTUI = &t
	m.updateLogsHeight()
	m.logsViewport.GotoBottom()
	m.tuiViewport.Height = m.windowHeight
	m.tuiViewport.Width = m.windowWidth
	m.tuiViewport, _ = m.tuiViewport.Update(t)

	switch t.Type {
	case tuiTypeSelectOne:
		m.selectOne.SetParams(t.Payload.(TUISelectOneParams))
	case tuiTypeInputText:
		m.inputText.SetParams(t.Payload.(TUIInputTextParams))
	case tuiTypeInputInt:
		m.inputInt.SetParams(t.Payload.(TUIInputIntParams))
	case tuiTypeInputFile:
		m.inputFile.SetParams(t.Payload.(TUIInputFileParams))
	case tuiTypeConfirm:
		m.confirm.SetParams(t.Payload.(TUIConfirmParams), t.Response)
	}

	return m, nil
}

func (m *model) handleTickMsg(msg tea.Msg) (tea.Model, tea.Cmd) {
	select {
	case l := <-m.logsChan:
		return m.onLogsChan(l, msg)
	case t := <-m.tuiChan:
		return m.onTUIChan(t, msg)
	case <-m.tuiClose:
		m.isRunningTUI = false
		m.runningTUI = nil
		m.tuiViewport.SetContent("")
		m.updateLogsHeight()
		m.logsViewport.GotoBottom()
		m.tuiViewport, _ = m.tuiViewport.Update(msg)
		m.logsViewport, _ = m.logsViewport.Update(msg)

		return m, ticker()
	default:
		if m.runningCommand == "" {
			m.input.running = false
			m.input, _ = m.input.Update(msg)

			return m, nil
		}

		return m, ticker()
	}
}

func (m *model) updateViewport(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	m.logsViewport, cmd = m.logsViewport.Update(msg)

	return cmd
}

// Update - BubbleTea model method.
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.handleWindowSizeMsg(msg)
	case inputResizeMsg:
		m.logsViewport.Height -= msg.Delta
		m.logsViewport, _ = m.logsViewport.Update(msg)

		return m, nil
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)
	case tick:
		return m.handleTickMsg(msg)
	case tea.MouseMsg:
		return m.handleMouseMsg(msg)
	default:
		if m.app.Params.EnableInputBlinking {
			m.input, _ = m.input.Update(msg)

			return m, nil
		}

		return m, nil
	}
}
