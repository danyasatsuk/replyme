package replyme

import (
	"errors"
	"io"
	"slices"
	"time"
)

func createCommandFlow(app *App, ast *ASTNode) ([]*Command, error) {
	cmds := make([]*Command, 0)

	if ast.Subcommands == nil || len(ast.Subcommands) == 0 {
		cmd, err := app.Commands.getCommand(ast.Command)
		if err != nil {
			return nil, err
		}

		err = insertDataInCommand(cmd, ast, false)
		if err != nil {
			return nil, err
		}

		cmds = append(cmds, cmd)

		return cmds, nil
	}

	cmd, err := app.Commands.getCommand(ast.Command)
	if err != nil {
		return nil, err
	}

	err = insertDataInCommand(cmd, ast, true)
	if err != nil {
		return nil, err
	}

	cmds = append(cmds, cmd)

	for i, subcommand := range ast.Subcommands {
		cmd, err := app.Commands.getCommand(subcommand)
		if err != nil {
			return nil, err
		}

		err = insertDataInCommand(cmd, ast, i != len(ast.Subcommands)-1)
		if err != nil {
			return nil, err
		}

		cmds = append(cmds, cmd)
	}

	return cmds, nil
}

//nolint:cyclop
func runActions(command *Command, ctx *Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch t := r.(type) {
			case error:
				err = newErrorCommandPanic(t.Error())
			case string:
				err = newErrorCommandPanic(t)
			default:
				err = newErrorCommandPanic("unknown panic")
			}
		}
	}()

	if command == nil {
		return nil
	}

	run := true

	if command.Before != nil {
		run, err = command.Before(ctx)
		if err != nil {
			return err
		}
	}

	if run && command.Action != nil {
		err = command.Action(ctx)
		if err != nil {
			return err
		}
	}

	return err
}

func runEnd(command *Command, ctx *Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch t := r.(type) {
			case error:
				err = newErrorCommandPanic(t.Error())
			case string:
				err = newErrorCommandPanic(t)
			default:
				err = newErrorCommandPanic("unknown panic")
			}
		}
	}()

	if command.OnEnd == nil {
		return nil
	}

	return command.OnEnd(ctx)
}

func appRunCleaner(app *App) {
	commandCleaner(app.Commands)
}

func commandCleaner(commands []*Command) {
	for _, command := range commands {
		for _, flag := range command.Flags {
			flag.Clear()
		}

		for _, arg := range command.Arguments {
			arg.value = ""
		}

		if command.Subcommands != nil {
			commandCleaner(command.Subcommands)
		}
	}
}

type fullRunCommandParams struct {
	command    string
	app        *App
	logsChan   chan<- log
	stdout     io.Writer
	stderr     io.Writer
	emitLog    func(logMsg)
	emitTUI    func(TUIRequest)
	emitTUICLI func(TUIRequest, chan<- bool)
	isCLI      bool
}

//nolint:cyclop,funlen,lll
func fullRunCommand(p fullRunCommandParams) error {
	defer func() {
		appRunCleaner(p.app)
	}()

	ast, err := parseCommand(createCommandSchema(p.app.Commands),
		createFlagSchema(p.app.Commands), createArgsSchema(p.app.Commands), p.command)
	if err != nil {
		return err
	}

	flow, err := createCommandFlow(p.app, ast)
	if err != nil {
		return err
	}

	if len(flow) > 0 {
		flags := flow[len(flow)-1].Flags
		if flagI := slices.IndexFunc(flags, func(f Flag) bool {
			if f.GetName() == "help" {
				if d, err := f.ParsedValue(); err == nil && d.(bool) && d.(bool) == true {
					return true
				}
			}

			return false
		}); flagI != -1 {
			help, err := helpCommand(flow[len(flow)-1])
			if err != nil {
				p.logsChan <- log{
					logTypeError,
					p.command,
					"ERROR",
					err,
					time.Now(),
				}

				return nil
			}
			p.logsChan <- log{
				logTypeMessage,
				p.command,
				help,
				nil,
				time.Now(),
			}
			p.logsChan <- log{
				logTypeCommandSuccess,
				p.command,
				p.command,
				nil,
				time.Now(),
			}

			return nil
		}
	}

	for _, cmd := range flow {
		ctx := createPreContext(cmd, ast)
		ctx.isCLI = p.isCLI
		ctx.emitLog = p.emitLog
		ctx.stdout = p.stdout
		ctx.stderr = p.stderr
		if p.isCLI {
			ctx.emitTUICLI = p.emitTUICLI
		} else {
			ctx.emitTUI = p.emitTUI
		}

		err = runActions(cmd, ctx)
		if err != nil {
			if errors.Is(err, ErrorCommandPanic) {

			}

			return err
		}
	}

	slices.Reverse(flow)

	for _, cmd := range flow {
		ctx := createPreContext(cmd, ast)
		ctx.isCLI = p.isCLI
		ctx.emitLog = p.emitLog
		ctx.stdout = p.stdout
		ctx.stderr = p.stderr
		if p.isCLI {
			ctx.emitTUICLI = p.emitTUICLI
		} else {
			ctx.emitTUI = p.emitTUI
		}

		err = runEnd(cmd, ctx)
		if err != nil {
			return err
		}
	}
	p.logsChan <- log{
		logTypeCommandSuccess,
		p.command,
		p.command,
		nil,
		time.Now(),
	}

	return nil
}

func (m *model) runCommand(command string) error {
	err := fullRunCommand(fullRunCommandParams{
		command, m.app, m.logsChan, m.stdout, m.stderr,
		m.emitLog, m.emitTUI, nil, false,
	})
	m.runningCommand = ""
	m.input.running = false

	return err
}

func runCommand(app *App, ctx *Context, command string) error {
	ast, err := parseCommand(createCommandSchema(app.Commands),
		createFlagSchema(app.Commands), createArgsSchema(app.Commands), command)
	if err != nil {
		return err
	}

	flow, err := createCommandFlow(app, ast)
	if err != nil {
		return err
	}

	for _, cmd := range flow {
		err = runActions(cmd, ctx)
		if err != nil {
			return err
		}
	}

	slices.Reverse(flow)

	for _, cmd := range flow {
		err = runEnd(cmd, ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
