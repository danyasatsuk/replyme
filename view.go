package replyme

// View - method of the BubbleTea model
func (m *Model) View() string {
	if m.isRunningTUI {
		switch m.runningTUI.Type {
		case TUIType_SelectOne:
			m.tuiViewport.SetContent(m.selectOne.View())
		case TUIType_SelectSeveral:
			//m.tuiViewport.SetContent(m.selectSeveral.View())
		case TUIType_InputText:
			m.tuiViewport.SetContent(m.inputText.View())
		case TUIType_InputInt:
			m.tuiViewport.SetContent(m.inputInt.View())
		case TUIType_InputFile:
			m.tuiViewport.SetContent(m.inputFile.View())
		case TUIType_Confirm:
			m.tuiViewport.SetContent(m.confirm.View())
		}
		return m.logsViewport.View() + "\n" + m.tuiViewport.View()
	}
	return m.logsViewport.View() + " \n" + m.input.View()
}
