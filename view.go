package replyme

// View - method of the BubbleTea model
func (m *model) View() string {
	if m.isRunningTUI {
		switch m.runningTUI.Type {
		case tuiType_SelectOne:
			m.tuiViewport.SetContent(m.selectOne.View())
		case tuiType_SelectSeveral:
			//m.tuiViewport.SetContent(m.selectSeveral.View())
		case tuiType_InputText:
			m.tuiViewport.SetContent(m.inputText.View())
		case tuiType_InputInt:
			m.tuiViewport.SetContent(m.inputInt.View())
		case tuiType_InputFile:
			m.tuiViewport.SetContent(m.inputFile.View())
		case tuiType_Confirm:
			m.tuiViewport.SetContent(m.confirm.View())
		}
		return m.logsViewport.View() + "\n" + m.tuiViewport.View()
	}
	return m.logsViewport.View() + " \n" + m.input.View()
}
