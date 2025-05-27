package replyme

import "github.com/charmbracelet/lipgloss"

var HelloHeader = lipgloss.NewStyle().
	Foreground(lipgloss.Color("237")).
	Background(lipgloss.Color("255")).
	Bold(true).
	Align(lipgloss.Center)

var HelloText = lipgloss.NewStyle().
	Foreground(lipgloss.Color("243")).
	Align(lipgloss.Center)

var RedIcon = lipgloss.NewStyle().
	Foreground(lipgloss.Color("196"))

var GreenIcon = lipgloss.NewStyle().
	Foreground(lipgloss.Color("40"))

var SelectedIcon = lipgloss.NewStyle().
	Background(lipgloss.Color("7")).Render

var SpinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("45"))

var GrayStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Render

var LogStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true).Render

var DebugStyle = lipgloss.NewStyle().Background(lipgloss.Color("12")).Bold(true).Render

var WarnStyle = lipgloss.NewStyle().Background(lipgloss.Color("11")).Bold(true).Render

var ErrorHeaderStyle = lipgloss.NewStyle().Background(lipgloss.Color("9")).Bold(true).Render

var ErrorTextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render

var CMDCommandStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render

var CMDFlagStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("4")).Render

var CMDFlagValueStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Render

var CMDArgValueStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("13")).Render

var CMDStringStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Render
