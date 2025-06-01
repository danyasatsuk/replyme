package replyme

import "github.com/charmbracelet/lipgloss"

var redIcon = lipgloss.NewStyle().
	Foreground(lipgloss.Color("196"))

var greenIcon = lipgloss.NewStyle().
	Foreground(lipgloss.Color("40"))

var inputTextContainer = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1)

type stylesStruct struct {
	GrayStyle, LogStyle, DebugStyle, WarnStyle, ErrorHeaderStyle,
	ErrorTextStyle, CMDCommandStyle, CMDFlagStyle, CMDFlagValueStyle,
	CMDArgValueStyle, CMDStringStyle, InputTextTitle, InputTextDescription func(strs ...string) string
}

var styles = stylesStruct{
	GrayStyle:            lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Render,
	LogStyle:             lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true).Render,
	DebugStyle:           lipgloss.NewStyle().Background(lipgloss.Color("12")).Bold(true).Render,
	WarnStyle:            lipgloss.NewStyle().Background(lipgloss.Color("11")).Bold(true).Render,
	ErrorHeaderStyle:     lipgloss.NewStyle().Background(lipgloss.Color("9")).Bold(true).Render,
	ErrorTextStyle:       lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render,
	CMDCommandStyle:      lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render,
	CMDFlagStyle:         lipgloss.NewStyle().Foreground(lipgloss.Color("4")).Render,
	CMDFlagValueStyle:    lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Render,
	CMDArgValueStyle:     lipgloss.NewStyle().Foreground(lipgloss.Color("13")).Render,
	CMDStringStyle:       lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Render,
	InputTextTitle:       lipgloss.NewStyle().Background(lipgloss.Color("4")).Padding(0, 2).Bold(true).Render,
	InputTextDescription: lipgloss.NewStyle().Foreground(lipgloss.Color("7")).Render,
}
