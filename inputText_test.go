package replyme

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
	"testing"
)

func TestInputTextTUI(t *testing.T) {
	err := i18nInit()
	if err != nil {
		t.Fatal(err)
	}

	m := inputTextNew(make(chan bool), true)

	msg := make(chan TUIResponse)

	m = m.SetParams(TUIInputTextParams{
		Name:        "InputText",
		Description: "Please enter your input:",
		Placeholder: "Type here...",
	}, msg)

	tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(30, 10))

	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("Test input"),
	})

	tm.Send(tea.KeyMsg{
		Type: tea.KeyEnter,
	})

	res := <-msg

	if res.Err != nil {
		t.Fatal(res.Err)
	}

	if res.Value != "Test input" {
		t.Fatalf("Expected 'Test input', got '%s'", res.Value)
	}

}
