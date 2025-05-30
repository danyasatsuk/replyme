package replyme

import (
	"bytes"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

// Model - the BubbleTea model
type Model struct {
	app *App
	ModelElements
	ModelLogs
	ModelTUI

	cachedMultiline bool
}

// ModelElements - the elements of the BubbleTea model
type ModelElements struct {
	logsViewport viewport.Model
	input        TerminalInput
	spinner      spinner.Model

	windowHeight int
	windowWidth  int
}

// ModelLogs - the logs of the BubbleTea model
type ModelLogs struct {
	logs                *Logs
	history             []string
	selectedHistoryItem int
	stdout              *bytes.Buffer
	stderr              *bytes.Buffer
	runningCommand      string
	logsDirty           bool

	logsChan chan Log
}

// ModelTUI - the TUI of the BubbleTea model
type ModelTUI struct {
	tuiChan      chan TUIRequest
	isRunningTUI bool
	tuiClose     chan bool
	runningTUI   *TUIRequest
	tuiViewport  viewport.Model

	selectOne *SelectOne
	inputText *InputText
	inputInt  *InputInt
	inputFile *InputFile
	confirm   *Confirm
}

func createViewport() viewport.Model {
	v := viewport.New(56, 10)
	v.MouseWheelEnabled = true
	return v
}

func createSpinner() spinner.Model {
	return spinner.New(spinner.WithSpinner(spinner.Pulse))
}

func createTUIViewport() viewport.Model {
	return viewport.New(56, 10)
}

func createInput() textinput.Model {
	t := textinput.New()
	t.Focus()
	t.Cursor = cursor.New()
	t.Prompt = GrayStyle(">> ")
	return t
}

// CreateModel - create a new BubbleTea model
func CreateModel(app *App) *Model {
	tuiClose := make(chan bool, 1)
	m := &Model{
		app: app,
		ModelTUI: ModelTUI{
			tuiViewport: createTUIViewport(),
			tuiChan:     make(chan TUIRequest),
			selectOne:   SelectOneNew(),
			inputText:   InputTextNew(),
			inputInt:    InputIntNew(),
			inputFile:   InputFileNew(),
			confirm:     ConfirmNew(tuiClose),
			tuiClose:    tuiClose,
		},
		ModelElements: ModelElements{
			logsViewport: createViewport(),
			input:        NewTerminalInput(),
			spinner:      createSpinner(),
		},
		ModelLogs: ModelLogs{
			logs:                &Logs{},
			history:             make([]string, 0),
			selectedHistoryItem: -1,
			logsChan:            make(chan Log),
			stdout:              bytes.NewBuffer(nil),
			stderr:              bytes.NewBuffer(nil),
		},
	}
	m.tuiViewport.Style = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("32"))
	return m
}
