package replyme

import (
	"golang.org/x/exp/slices"
)

// Command - the structure for creating your command
type Command struct {
	// Command name
	Name string
	// What is this command for?
	Usage string
	// Allows you to specify abbreviations for the command
	Aliases []string
	// The subcommands of your main command
	Subcommands Commands
	// Flags for executing your command
	Flags Flags
	// Arguments for executing your command
	Arguments []*Argument
	// The function that is executed before executing the main function Action
	Before func(ctx *Context) (bool, error)
	// The main function of the command
	Action func(ctx *Context) error
	// A function that runs after the main Action function has successfully completed its action
	OnEnd func(ctx *Context) error
}

// Commands is an abbreviation for the type `[]*replyme.Command`
type Commands []*Command

func (c Commands) getCommand(name string) (*Command, error) {
	cmdsArr := c.getCommandsArray()
	i := slices.IndexFunc(cmdsArr, func(command *Command) bool {
		return command.Name == name
	})
	if i == -1 {
		return nil, newErrorUnknownCommand(name)
	}
	return cmdsArr[i], nil
}

func subber(commands *Command) []*Command {
	s := make([]*Command, 0)
	if commands.Subcommands != nil {
		for _, subcommand := range commands.Subcommands {
			if subcommand.Subcommands != nil {
				s = append(s, subber(subcommand)...)
			}
			s = append(s, subcommand)
		}
	}
	return s
}

func (c Commands) getCommandsArray() []*Command {
	cmds := make([]*Command, 0)
	cmds = append(cmds, c...)
	for _, command := range c {
		cmds = append(cmds, subber(command)...)
	}
	return cmds
}

func (c Commands) mustGetCommand(name string) *Command {
	i := slices.IndexFunc(c, func(command *Command) bool {
		return command.Name == name
	})
	if i == -1 {
		panic("unknown command: " + name)
	}
	return c[i]
}
