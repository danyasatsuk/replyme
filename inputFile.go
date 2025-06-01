package replyme

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const kilobyte = 1024

var (
	okStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
)

type inputFile struct {
	input       textinput.Model
	IsValidated bool
	IsExit      bool
	Value       TUIInputFileResult
	params      TUIInputFileParams
	statusLine  string
	statusStyle lipgloss.Style
	isCLI       bool
	c           chan TUIResponse
	close       chan bool
}

func inputFileNew(c chan bool, isCLI ...bool) inputFile {
	var cli bool
	if len(isCLI) > 0 && isCLI[0] {
		cli = true
	}

	m := inputFile{
		input: textinput.New(),
		isCLI: cli,
		close: c,
	}

	return m
}

func (m inputFile) SetParams(p TUIInputFileParams, c chan TUIResponse) inputFile {
	m.params = p
	m.input.Placeholder = L(i18n_inputfile_placeholder)
	m.input.Focus()
	m.c = c

	return m
}

func (m inputFile) Init() tea.Cmd {
	return nil
}

func (m inputFile) Focus() {
	m.input.Focus()
}

func (m inputFile) Blur() {
	m.input.Blur()
}

func (m inputFile) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyInsert && len(msg.String()) > 1 {
			drop := strings.Trim(msg.String(), "\"'")
			m.input.SetValue(drop)
			m.updateStatus()

			return m, nil
		}

		switch msg.String() {
		case "ctrl+c", "esc":
			m.IsExit = true

			if m.isCLI {
				return m, tea.Quit
			}

			return m, nil

		case "enter":
			return m.onEnter()
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)

	m.updateStatus()

	return m, cmd
}

func (m inputFile) View() string {
	return fmt.Sprintf(`%s

%s

%s

%s`, m.params.Name, m.params.Description, m.input.View(), m.statusStyle.Render(m.statusLine))
}

func (m inputFile) checkPath() (os.FileInfo, string, error) {
	path := m.input.Value()

	abs, err := filepath.Abs(path)
	if err != nil {
		m.setStatus(L(i18n_inputfile_fullpath_error), errorStyle)

		return nil, "", err
	}

	stat, err := os.Stat(abs)
	if err != nil || stat.IsDir() {
		m.setStatus(L(i18n_inputfile_file_notfound), errorStyle)

		return nil, "", err
	}

	return stat, abs, nil
}

func (m inputFile) onEnter() (tea.Model, tea.Cmd) {
	stat, abs, err := m.checkPath()
	if err != nil {
		return m, nil
	}

	if len(m.params.Extensions) > 0 {
		ok := false
		ext := strings.ToLower(filepath.Ext(abs))

		for _, allowed := range m.params.Extensions {
			if strings.ToLower(allowed) == ext {
				ok = true

				break
			}
		}

		if !ok {
			m.setStatus(L(i18n_inputfile_extension_error), errorStyle)

			return m, nil
		}
	}

	if m.params.MaxFileSize > 0 {
		sizeKB := stat.Size() / kilobyte
		if int(sizeKB) > m.params.MaxFileSize {
			m.setStatus(fmt.Sprintf(L(i18n_inputfile_size_error), sizeKB, m.params.MaxFileSize), errorStyle)

			return m, nil
		}
	}

	var contents []byte

	if !m.params.DoNotOutput {
		data, err := os.ReadFile(abs)
		if err != nil {
			m.setStatus(L(i18n_inputfile_read_error), errorStyle)

			return m, nil
		}

		contents = data
	}

	m.Value = TUIInputFileResult{
		Path: abs,
		File: contents,
	}
	m.IsValidated = true
	m.input.Reset()
	m.setStatus(L(i18n_inputfile_success), okStyle)

	m.c <- TUIResponse{
		Value: m.Value,
		Err:   nil,
	}

	if m.isCLI {
		return m, tea.Quit
	}

	m.close <- true

	return m, nil
}

func (m inputFile) setStatus(text string, style lipgloss.Style) {
	m.statusLine = text
	m.statusStyle = style
}

func (m inputFile) updateStatus() {
	path := m.input.Value()

	absPath, err := filepath.Abs(path)
	if err != nil {
		m.setStatus(L(i18n_inputfile_path_error), errorStyle)

		return
	}

	if stat, err := os.Stat(absPath); err == nil && !stat.IsDir() {
		m.setStatus(absPath, okStyle)
	} else {
		m.setStatus(absPath, errorStyle)
	}
}
