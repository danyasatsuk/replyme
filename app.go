package replyme

import "golang.org/x/exp/slices"

// AppParams - the structure of the application parameters
type AppParams struct {
	// Allows you to turn on the cursor blinking inside the input. It's not working yet, added to the Roadmap.
	// TODO(unimportant): Add cursor blinking
	EnableInputBlinking bool
}

// App - the structure of the application
type App struct {
	// The name of your application
	Name string
	// Description of the application, why it is needed
	Usage string
	// The authors of the application
	Authors []string
	// Your copyright in the form of "YEAR-YEAR author"
	Copyright string
	// The license under which you distribute the code
	License string

	// A list of all your commands
	Commands Commands

	// Allows you to enable Debug mode (with it, all Debug messages are output to the console)
	Debug bool
	// Allows you to disable the color output
	NoColor bool
	// Application parameters. For more information, see AppParams.
	Params AppParams
}

// GetFlagSchema - allows you to get a diagram of all flags in the FlagSchema type
func (a *App) GetFlagSchema() FlagSchema {
	return parseFlagSchema(a.Commands)
}

func parseFlagSchema(commands Commands) FlagSchema {
	schema := FlagSchema{}
	for _, command := range commands {
		newSchema := parseFlagSchemaSingle(command)
		for k, v := range newSchema {
			schema[k] = v
		}
	}
	return schema
}

func parseFlagSchemaSingle(command *Command) FlagSchema {
	allFlags := make(map[string]map[string]FlagType)
	f := make(map[string]FlagType)
	for _, flag := range command.Flags {
		switch flag.ValueType() {
		case "bool":
			f[flag.GetName()] = FlagTypeBool
		case "string":
			f[flag.GetName()] = FlagTypeString
		case "int":
			f[flag.GetName()] = FlagTypeInt
		case "[]string":
			f[flag.GetName()] = FlagTypeStringArray
		case "[]int":
			f[flag.GetName()] = FlagTypeIntArray
		}
	}
	if len(command.Subcommands) > 0 {
		schema := parseFlagSchema(command.Subcommands)
		for k, v := range schema {
			allFlags[k] = v
		}
	}
	allFlags[command.Name] = f
	return allFlags
}

func (a *App) setHelpFlags() {
	a.Commands = setHelpFlag(a.Commands)
}

func setHelpFlag(commands []*Command) []*Command {
	for i := range commands {
		if slices.IndexFunc(commands[i].Flags, func(flag Flag) bool {
			return flag.GetName() == "help"
		}) == -1 {
			commands[i].Flags = append(commands[i].Flags, &FlagValue[bool]{
				Name:  "help",
				Alias: "h",
				Usage: L(i18n_app_help_usage),
			})
		}
		if commands[i].Subcommands != nil && len(commands[i].Subcommands) > 0 {
			commands[i].Subcommands = setHelpFlag(commands[i].Subcommands)
		}
	}
	return commands
}
