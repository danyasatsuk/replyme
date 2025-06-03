package replyme

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/danyasatsuk/replyme/internal/filepicker"
	Log "log"
	"os"
	"strings"
	"time"
)

var bottomPadding = 10

type inputFile struct {
	picker      filepicker.Model
	IsValidated bool
	IsExit      bool
	Value       TUIInputFileResult
	params      TUIInputFileParams
	isCLI       bool
	c           chan TUIResponse
	close       chan bool
	width       int
	height      int
	err         string
}

type clearErrorMsg struct{}

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

func inputFileNew(c chan bool, isCLI ...bool) inputFile {
	var cli bool
	if len(isCLI) > 0 && isCLI[0] {
		cli = true
	}

	p := filepicker.New()
	p.CurrentDirectory, _ = os.Getwd()

	m := inputFile{
		picker: p,
		isCLI:  cli,
		close:  c,
	}

	return m
}

func (m inputFile) SetParams(p TUIInputFileParams, c chan TUIResponse) inputFile {
	m.params = p

	if !m.isCLI {
		m.picker = m.picker.SetHeight(m.height - bottomPadding)
	}

	if p.Extensions != nil && len(p.Extensions) > 0 {
		m.picker.AllowedTypes = p.Extensions
	}

	m.IsValidated = false
	m.IsExit = false
	m.c = c

	return m
}

func (m inputFile) Init() tea.Cmd {
	m.picker.CurrentDirectory, _ = os.Getwd()
	m.picker.ShowHidden = true

	return m.picker.Init()
}

func (m inputFile) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		var cmd tea.Cmd

		m.width = msg.Width
		m.height = msg.Height
		m.picker = m.picker.SetHeight(m.height - bottomPadding)

		return m, cmd
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.IsExit = true
			if m.isCLI {
				return m, tea.Quit
			}
			m.close <- true

			return m, nil
		case "y":
			Log.Print("test")
		}
	case clearErrorMsg:
		m.err = ""
	}

	var cmd tea.Cmd

	m.picker, cmd = m.picker.Update(msg)

	if didSelect, path := m.picker.DidSelectFile(msg); didSelect {
		var file []byte

		if !m.params.DoNotOutput {
			f, err := os.ReadFile(path)
			if err != nil {
				m.c <- TUIResponse{Err: err}
			}

			file = f
		}

		m.c <- TUIResponse{
			Value: TUIInputFileResult{
				Path: path,
				File: file,
			},
		}

		if m.isCLI {
			return m, tea.Quit
		}

		m.close <- true
	}

	if didSelect, _ := m.picker.DidSelectDisabledFile(msg); didSelect {
		m.err = L(i18n_tui_inputFile_err)

		return m, tea.Batch(cmd, clearErrorAfter(2*time.Second))
	}

	return m, cmd
}

func (m inputFile) View() string {
	if m.err != "" {
		return inputContainer.Width(m.width - 2).Height(m.height - 2).Render(fmt.Sprintf(`%s

%s
%s

%s`, styles.InputTitle(m.params.Name), styles.ErrorTextStyle(m.err), styles.InputDescription(m.params.Description), m.picker.View()))
	}

	if m.params.Extensions != nil && len(m.params.Extensions) > 0 {
		return inputContainer.Width(m.width - 2).Height(m.height - 2).Render(fmt.Sprintf(`%s

%s
%s

%s`, styles.InputTitle(m.params.Name), styles.InputDescription(m.params.Description),
			styles.InputDescription("("+strings.Join(m.params.Extensions, ", ")+")"), m.picker.View()))
	}

	return inputContainer.Width(m.width - 2).Height(m.height - 2).Render(fmt.Sprintf(`%s

%s


%s`, styles.InputTitle(m.params.Name),
		styles.InputDescription(m.params.Description), m.picker.View()))
}
