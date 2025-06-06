package replyme

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
	"testing"
)

func TestConfirmTUI(t *testing.T) {
	err := i18nInit()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Yes", func(t *testing.T) {
		m := confirmNew(make(chan bool), true)

		msg := make(chan TUIResponse)

		m = m.SetParams(TUIConfirmParams{
			Name:        "Confirm",
			Description: "Are you sure you want to proceed?",
		}, msg)

		tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(30, 10))

		tm.Send(tea.KeyMsg{
			Type: tea.KeyEnter,
		})

		res := <-msg

		if res.Err != nil {
			t.Fatal(res.Err)
		}

		if res.Value != true {
			t.Fatal("Expected true, got false")
		}
	})

	t.Run("No", func(t *testing.T) {
		m := confirmNew(make(chan bool), true)

		msg := make(chan TUIResponse)

		m = m.SetParams(TUIConfirmParams{
			Name:        "Confirm",
			Description: "Are you sure you want to proceed?",
		}, msg)

		tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(30, 10))

		tm.Send(tea.KeyMsg{
			Type: tea.KeyRight,
		})

		tm.Send(tea.KeyMsg{
			Type: tea.KeyEnter,
		})

		res := <-msg

		if res.Err != nil {
			t.Fatal(res.Err)
		}

		if res.Value != false {
			t.Fatal("Expected false, got true")
		}
	})
}
