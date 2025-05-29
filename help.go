package replyme

import (
	"bytes"
	"github.com/charmbracelet/lipgloss"
	"strings"
	"text/template"
)

type HelpStruct struct {
	Name        string
	Authors     []string
	License     string
	Usage       string
	Flags       []HelpFlagsStruct
	Arguments   []HelpArgumentsStruct
	Subcommands []HelpSubcommandsStruct
	I18n        HelpI18nStruct
}

type HelpFlagsStruct struct {
	Name  string
	Usage string
	Type  string
	Alias string
}

type HelpArgumentsStruct struct {
	Name  string
	Usage string
}

type HelpSubcommandsStruct struct {
	Name  string
	Usage string
}

type HelpI18nStruct struct {
	Authors     string
	Subcommands string
	Flags       string
	Arguments   string
	License     string
}

var HelpCommandTemplate = `{{ Bold .Name }} - {{ .Usage }}

{{ if .Authors }}{{ Bold .I18n.Authors }}:
  {{ StringsJoin .Authors ", " }}
{{ end }}{{ if .License }}{{ Bold .I18n.License }}:
  {{ .License }}
{{ end }}{{ if .Subcommands }}{{ Bold .I18n.Subcommands }}:
{{ range .Subcommands }}  {{ Green .Name }} - {{ .Usage }}
{{ end }}{{ end }}{{ if .Flags }}{{ Bold .I18n.Flags}}:
{{ range .Flags }}  --{{ Blue .Name }}{{ if .Alias }}(-{{ Gray .Alias }}){{ end }}{{ Cyan .Type}} - {{ .Usage }}
{{ end }}{{ end }}{{ if .Arguments }}{{ Bold .I18n.Arguments }}:
{{ range .Arguments }}  {{ Purple .Name }} - {{ .Usage }}
{{ end }}{{ end }}`

func buildHelpFlags(command *Command) []HelpFlagsStruct {
	if command.Flags == nil {
		return nil
	}
	flags := make([]HelpFlagsStruct, len(command.Flags))
	for i, flag := range command.Flags {
		flags[i] = HelpFlagsStruct{
			Name:  flag.GetName(),
			Usage: flag.GetUsage(),
			Alias: flag.GetAlias(),
		}
		switch flag.ValueType() {
		case "string":
			flags[i].Type = "=" + L(i18n_help_flag_type_string)
		case "bool":
			flags[i].Type = ""
		case "int":
			flags[i].Type = "=" + L(i18n_help_flag_type_int)
		case "[]string":
			flags[i].Type = "=" + L(i18n_help_flag_type_string_array)
		case "[]int":
			flags[i].Type = "=" + L(i18n_help_flag_type_int_array)
		}
	}
	return flags
}

func buildHelpArguments(command *Command) []HelpArgumentsStruct {
	if command.Arguments == nil {
		return nil
	}
	args := make([]HelpArgumentsStruct, len(command.Arguments))
	for i, arg := range command.Arguments {
		args[i] = HelpArgumentsStruct{
			Name:  arg.Name,
			Usage: arg.Usage,
		}
	}
	return args
}

func buildHelpSubcommands(command *Command) []HelpSubcommandsStruct {
	if command.Subcommands == nil {
		return nil
	}
	subcommands := make([]HelpSubcommandsStruct, len(command.Subcommands))
	for i, subcommand := range command.Subcommands {
		subcommands[i] = HelpSubcommandsStruct{
			Name:  subcommand.Name,
			Usage: subcommand.Usage,
		}
	}
	return subcommands
}

func buildHelpCommands(app *App) []HelpSubcommandsStruct {
	if app.Commands == nil {

	}
	commands := make([]HelpSubcommandsStruct, len(app.Commands))
	for i, subcommand := range app.Commands {
		commands[i] = HelpSubcommandsStruct{
			Name:  subcommand.Name,
			Usage: subcommand.Usage,
		}
	}
	return commands
}

func buildHelpI18n() HelpI18nStruct {
	return HelpI18nStruct{
		Authors:     L(i18n_help_authors),
		Subcommands: L(i18n_help_subcommands),
		Flags:       L(i18n_help_flags),
		Arguments:   L(i18n_help_arguments),
		License:     L(i18n_help_license),
	}
}

var tmpl *template.Template

func createTemplate() {
	if tmpl == nil {
		tmpl = template.New("help")
		tmpl = tmpl.Funcs(template.FuncMap{
			"StringsJoin": strings.Join,
			"Bold":        lipgloss.NewStyle().Bold(true).Render,
			"Blue":        lipgloss.NewStyle().Foreground(lipgloss.Color("4")).Render,
			"Green":       lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Render,
			"Purple":      lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Render,
			"Gray":        lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render,
			"Cyan":        lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Render,
		})
	}
}

func HelpCommand(command *Command) (string, error) {
	createTemplate()
	t := HelpStruct{
		Name:        command.Name,
		Usage:       command.Usage,
		Flags:       buildHelpFlags(command),
		Arguments:   buildHelpArguments(command),
		Subcommands: buildHelpSubcommands(command),
		I18n:        buildHelpI18n(),
	}

	buf := &bytes.Buffer{}
	parse, err := tmpl.Parse(HelpCommandTemplate)
	if err != nil {
		return "", err
	}
	err = parse.Execute(buf, t)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func HelpApp(app *App) (string, error) {
	createTemplate()
	t := HelpStruct{
		Name:        app.Name,
		Usage:       app.Usage,
		Authors:     app.Authors,
		License:     app.License,
		Subcommands: buildHelpCommands(app),
		I18n:        buildHelpI18n(),
	}

	buf := &bytes.Buffer{}
	parse, err := tmpl.Parse(HelpCommandTemplate)
	if err != nil {
		return "", err
	}
	err = parse.Execute(buf, t)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
