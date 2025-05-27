package replyme

import (
	"errors"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/wordwrap"
)

func (m *Model) tuiUpdater(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.runningTUI.Type {
	case TUIType_SelectOne:
		m.selectOne, cmd = m.selectOne.Update(msg)
	case TUIType_SelectSeveral:
		// TODO(medium): Implement SelectSeveral and add
	case TUIType_InputText:
		m.inputText, cmd = m.inputText.Update(msg)
	case TUIType_InputInt:
		m.inputInt, cmd = m.inputInt.Update(msg)
	case TUIType_InputFile:
		m.inputFile, cmd = m.inputFile.Update(msg)
	case TUIType_Confirm:
		m.confirm, cmd = m.confirm.Update(msg)
	}
	return m, tea.Batch(cmd, Ticker())
}

func (m *Model) updateLogsHeight() {
	if m.isRunningTUI {
		m.logsViewport.Height = 0
		m.tuiViewport.Height = m.windowHeight
	} else {
		m.tuiViewport.Height = 0
		m.logsViewport.Height = m.windowHeight - m.input.GetLines()
	}
}

func (m *Model) handleWindowSizeMsg(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
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

func (m *Model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "enter":
		command := m.input.Value()
		if command != "" {
			if command == "exit" {
				return m, tea.Quit
			}
			m.logs.Add(LogTypeCommandRunning, command)
			m.logsViewport.SetContent(wordwrap.String(m.logs.Render(), m.logsViewport.Width))
			m.logsViewport, _ = m.logsViewport.Update(msg)
			m.runningCommand = command
			m.input.running = true
			m.input, _ = m.input.Update(msg)

			go func() {
				err := m.runCommand(command)
				if err != nil {
					typeOfError := LogTypeError
					switch {
					case errors.Is(err, ErrorUnknownCommand):
						typeOfError = LogTypeCommandNotFound
					case errors.Is(err, ErrorCommandPanic):
						typeOfError = LogTypePanic
					}
					m.logsChan <- Log{typeOfError, command, err.Error(), err, time.Now()}
					m.logsChan <- Log{LogTypeCommandFailure, command, err.Error(), err, time.Now()}
				}
			}()

			return m, tea.Batch(Ticker())
		}
	}
	if m.isRunningTUI {
		return m.tuiUpdater(msg)
	}
	m.input, _ = m.input.Update(msg)
	return m, nil
}

func (m *Model) handleTickMsg(msg tea.Msg) (tea.Model, tea.Cmd) {
	select {
	case l := <-m.logsChan:
		l.Message = wordwrap.String(l.Message, m.logsViewport.Width)
		m.logs.AddLog(l)
		m.updateLogsHeight()
		m.logsViewport.SetContent(m.logs.Render())
		m.logsViewport.GotoBottom()
		m.logsViewport, _ = m.logsViewport.Update(msg)
		if l.Type == LogTypeCommandSuccess || l.Type == LogTypeCommandFailure || l.Type == LogTypeCommandNotFound {
			m.input.running = false
			m.input, _ = m.input.Update(msg)
			return m, tea.Batch(Ticker())
		}
		return m, Ticker()

	case t := <-m.tuiChan:
		m.isRunningTUI = true
		m.runningTUI = &t
		m.updateLogsHeight()
		m.logsViewport.GotoBottom()
		m.tuiViewport.Height = m.windowHeight
		m.tuiViewport.Width = m.windowWidth
		m.tuiViewport, _ = m.tuiViewport.Update(t)
		switch t.Type {
		case TUIType_SelectOne:
			m.selectOne.SetParams(t.Payload.(TUISelectOneParams))
		case TUIType_InputText:
			m.inputText.SetParams(t.Payload.(TUIInputTextParams))
		case TUIType_InputInt:
			m.inputInt.SetParams(t.Payload.(TUIInputIntParams))
		case TUIType_InputFile:
			m.inputFile.SetParams(t.Payload.(TUIInputFileParams))
		case TUIType_Confirm:
			m.confirm.SetParams(t.Payload.(TUIConfirmParams), t.Response)
		}
		return m, nil

	case <-m.tuiClose:
		m.isRunningTUI = false
		m.runningTUI = nil
		m.tuiViewport.SetContent("")
		m.updateLogsHeight()
		m.logsViewport.GotoBottom()
		m.tuiViewport, _ = m.tuiViewport.Update(msg)
		m.logsViewport, _ = m.logsViewport.Update(msg)
		return m, Ticker()

	default:
		if m.runningCommand == "" {
			m.input.running = false
			m.input, _ = m.input.Update(msg)
			return m, nil
		}
		return m, Ticker()
	}
}

func (m *Model) updateViewport(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	m.logsViewport, cmd = m.logsViewport.Update(msg)
	return cmd
}

// Update - BubbleTea model method
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.handleWindowSizeMsg(msg)
	case InputResizeMsg:
		m.logsViewport.Height -= msg.Delta
		m.logsViewport, _ = m.logsViewport.Update(msg)
		return m, nil
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)
	case Tick:
		return m.handleTickMsg(msg)
	default:
		if m.app.Params.EnableInputBlinking {
			m.input, _ = m.input.Update(msg)
			return m, nil
		}
		return m, nil
	}
}
