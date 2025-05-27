package replyme

import (
	"errors"
	"slices"
	"time"
)

func createCommandFlow(app *App, ast *ASTNode) ([]*Command, error) {
	cmds := make([]*Command, 0)
	if ast.Subcommands == nil || len(ast.Subcommands) == 0 {
		cmd, err := app.Commands.GetCommand(ast.Command)
		if err != nil {
			return nil, err
		}
		err = insertDataInCommand(cmd, ast, false)
		if err != nil {
			return nil, err
		}
		cmds = append(cmds, cmd)
	} else {
		cmd, err := app.Commands.GetCommand(ast.Command)
		if err != nil {
			return nil, err
		}
		err = insertDataInCommand(cmd, ast, true)
		if err != nil {
			return nil, err
		}
		cmds = append(cmds, cmd)
		for i, subcommand := range ast.Subcommands {
			cmd, err := app.Commands.GetCommand(subcommand)
			if err != nil {
				return nil, err
			}
			err = insertDataInCommand(cmd, ast, i != len(ast.Subcommands)-1)
			if err != nil {
				return nil, err
			}
			cmds = append(cmds, cmd)
		}
	}
	return cmds, nil
}

func runActions(command *Command, ctx *Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch t := r.(type) {
			case error:
				err = NewErrorCommandPanic(t.Error())
			case string:
				err = NewErrorCommandPanic(t)
			default:
				err = NewErrorCommandPanic("unknown panic")
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
	if run {
		if command.Action != nil {
			err = command.Action(ctx)
			if err != nil {
				return err
			}
		}
	}
	return err
}

func runEnd(command *Command, ctx *Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch t := r.(type) {
			case error:
				err = NewErrorCommandPanic(t.Error())
			case string:
				err = NewErrorCommandPanic(t)
			default:
				err = NewErrorCommandPanic("unknown panic")
			}
		}
	}()
	if command.OnEnd == nil {
		return nil
	}
	return command.OnEnd(ctx)
}

func (m *Model) runCommand(command string) error {
	ast, err := ParseCommand(createCommandSchema(m.app.Commands), createFlagSchema(m.app.Commands), createArgsSchema(m.app.Commands), command)
	if err != nil {
		return err
	}
	flow, err := createCommandFlow(m.app, ast)
	if err != nil {
		return err
	}
	for _, cmd := range flow {
		ctx := createPreContext(cmd, ast)
		ctx.emitLog = m.emitLog
		ctx.stdout = m.stdout
		ctx.stderr = m.stderr
		ctx.emitTUI = m.emitTUI

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
		ctx.emitLog = m.emitLog
		ctx.stdout = m.stdout
		ctx.stderr = m.stderr
		ctx.emitTUI = m.emitTUI
		err = runEnd(cmd, ctx)
		if err != nil {
			return err
		}
	}
	m.logsChan <- Log{
		LogTypeCommandSuccess,
		command,
		command,
		nil,
		time.Now(),
	}
	m.runningCommand = ""
	m.input.running = false
	return nil
}

func runCommand(app *App, ctx *Context, command string) error {
	ast, err := ParseCommand(createCommandSchema(app.Commands), createFlagSchema(app.Commands), createArgsSchema(app.Commands), command)
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
