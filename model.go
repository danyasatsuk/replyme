package replyme

import (
	"bytes"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

const standardWidth = 56
const standardHeight = 10

type model struct {
	app *App
	modelElements
	modelLogs
	modelTUI

	cachedMultiline bool
}

type modelElements struct {
	logsViewport viewport.Model
	input        terminalInput
	spinner      spinner.Model

	windowHeight int
	windowWidth  int
}

type modelLogs struct {
	logs                *logs
	history             []string
	selectedHistoryItem int
	stdout              *bytes.Buffer
	stderr              *bytes.Buffer
	runningCommand      string
	logsDirty           bool

	logsChan chan log
}

type modelTUI struct {
	tuiChan      chan TUIRequest
	isRunningTUI bool
	tuiClose     chan bool
	runningTUI   *TUIRequest
	tuiViewport  viewport.Model

	selectOne selectOne
	inputText inputText
	inputInt  inputInt
	inputFile inputFile
	confirm   confirm
}

func createViewport() viewport.Model {
	v := viewport.New(standardWidth, standardHeight)
	v.MouseWheelEnabled = true

	return v
}

func createSpinner() spinner.Model {
	return spinner.New(spinner.WithSpinner(spinner.Pulse))
}

func createTUIViewport() viewport.Model {
	return viewport.New(standardWidth, standardHeight)
}

func createInput() textinput.Model {
	t := textinput.New()
	t.Focus()
	t.Cursor = cursor.New()
	t.Prompt = styles.GrayStyle(">> ")

	return t
}

func createModel(app *App) *model {
	tuiClose := make(chan bool, 1)
	m := &model{
		app: app,
		modelTUI: modelTUI{
			tuiViewport: createTUIViewport(),
			tuiChan:     make(chan TUIRequest),
			selectOne:   selectOneNew(),
			inputText:   inputTextNew(),
			inputInt:    inputIntNew(),
			inputFile:   inputFileNew(),
			confirm:     confirmNew(tuiClose),
			tuiClose:    tuiClose,
		},
		modelElements: modelElements{
			logsViewport: createViewport(),
			input:        newTerminalInput(),
			spinner:      createSpinner(),
		},
		modelLogs: modelLogs{
			logs:                &logs{},
			history:             make([]string, 0),
			selectedHistoryItem: -1,
			logsChan:            make(chan log),
			stdout:              bytes.NewBuffer(nil),
			stderr:              bytes.NewBuffer(nil),
		},
	}
	m.tuiViewport.Style = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("32"))

	return m
}
