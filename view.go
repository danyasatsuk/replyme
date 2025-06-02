package replyme

// View - method of the BubbleTea model.
func (m *model) View() string {
	if m.isRunningTUI {
		switch m.runningTUI.Type {
		case tuiTypeSelectOne:
			m.tuiViewport.SetContent(m.selectOne.View())
		case tuiTypeInputText:
			m.tuiViewport.SetContent(m.inputText.View())
		case tuiTypeInputInt:
			m.tuiViewport.SetContent(m.inputInt.View())
		case tuiTypeInputFile:
			m.tuiViewport.SetContent(m.inputFile.View())
		case tuiTypeConfirm:
			m.tuiViewport.SetContent(m.confirm.View())
		}

		return m.logsViewport.View() + "\n" + m.tuiViewport.View()
	}

	return m.logsViewport.View() + " \n" + m.input.View()
}
