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

var (
	okStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("10")) // зелёный
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))  // красный
)

type InputFile struct {
	input       textinput.Model
	IsValidated bool
	IsExit      bool
	Value       TUIInputFileResult
	params      TUIInputFileParams
	statusLine  string
	statusStyle lipgloss.Style
}

func InputFileNew() *InputFile {
	m := &InputFile{
		input: textinput.New(),
	}
	return m
}

func (m *InputFile) SetParams(p TUIInputFileParams) {
	m.params = p
	m.input.Placeholder = "Путь до файла"
	m.input.Focus()
}

func (m *InputFile) Init() tea.Cmd {
	return nil
}

func (m *InputFile) Focus() {
	m.input.Focus()
}

func (m *InputFile) Blur() {
	m.input.Blur()
}

func (m *InputFile) Update(msg tea.Msg) (*InputFile, tea.Cmd) {
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
			return m, nil

		case "enter":
			path := m.input.Value()
			abs, err := filepath.Abs(path)
			if err != nil {
				m.setStatus("Невозможно определить абсолютный путь", errorStyle)
				return m, nil
			}

			stat, err := os.Stat(abs)
			if err != nil || stat.IsDir() {
				m.setStatus("Файл не найден или это папка", errorStyle)
				return m, nil
			}

			// Проверка на расширение
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
					m.setStatus("Недопустимое расширение файла", errorStyle)
					return m, nil
				}
			}

			// Проверка на размер
			if m.params.MaxFileSize > 0 {
				sizeKB := stat.Size() / 1024
				if int(sizeKB) > m.params.MaxFileSize {
					m.setStatus(fmt.Sprintf("Файл слишком большой: %d KB > %d KB", sizeKB, m.params.MaxFileSize), errorStyle)
					return m, nil
				}
			}

			// Всё прошло — читаем файл
			var contents []byte
			if !m.params.DoNotOutput {
				data, err := os.ReadFile(abs)
				if err != nil {
					m.setStatus("Ошибка при чтении файла", errorStyle)
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
			m.setStatus("Файл успешно выбран", okStyle)
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)

	m.updateStatus()

	return m, cmd
}

func (m *InputFile) View() string {
	return fmt.Sprintf(`%s

%s

%s

%s`, m.params.Name, m.params.Description, m.input.View(), m.statusStyle.Render(m.statusLine))
}

func (m *InputFile) setStatus(text string, style lipgloss.Style) {
	m.statusLine = text
	m.statusStyle = style
}

func (m *InputFile) updateStatus() {
	path := m.input.Value()
	absPath, err := filepath.Abs(path)
	if err != nil {
		m.setStatus("Невозможно определить путь", errorStyle)
		return
	}
	if stat, err := os.Stat(absPath); err == nil && !stat.IsDir() {
		m.setStatus(absPath, okStyle)
	} else {
		m.setStatus(absPath, errorStyle)
	}
}
