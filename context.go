package replyme

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/google/uuid"
	"io"
	"os/exec"
	"time"
)

type ctxInterface interface { //nolint:interfacebloat
	GetName() string
	GetCommandNameTree() []string
	GetFlagInt(name string, defaultValue int) int
	GetFlagString(name string, defaultValue string) string
	GetFlagIntArray(name string) []int
	GetFlagStringArray(name string) []string
	GetFlagBool(name string) bool
	Print(data ...interface{})
	Printf(format string, data ...interface{})
	PrintMarkdown(markdown string, data ...interface{})
	Warn(data ...interface{})
	Warnf(format string, data ...interface{})
	Error(data ...interface{})
	Errorf(format string, data ...interface{})
	StartTime() time.Time
	Elapsed() time.Duration
	Command() string
	Stdout() io.Writer
	Stderr() io.Writer
	Ctx() context.Context
	Done() <-chan struct{}
	IsCancelled() bool
	Set(key string, value interface{})
	Delete(key string)
	Get(key string) interface{}
	MustGetString(key string) string
	MustGetInt(key string) int
	Exec(cmd string, args ...string) (string, string, error)
	ExecLive(cmd string, args ...string) error
	ExecSilent(cmd string, args ...string) error
	streamOutput(r io.Reader, status logMsgStatus)
	SelectOne(p *TUISelectOneParams) (TUISelectOneResult, error)
	InputText(p *TUIInputTextParams) (string, error)
	InputInt(p *TUIInputIntParams) (int, error)
	InputFile(p *TUIInputFileParams) (TUIInputFileResult, error)
	Confirm(p *TUIConfirmParams) (bool, error)
}

// LogMsgStatus is an enum for describing the status of a log message.
type logMsgStatus uint16

const (
	logMsgStatusPrint logMsgStatus = iota
	logMsgStatusPrintf
	logMsgStatusPrintMarkdown
	logMsgStatusWarn
	logMsgStatusWarnf
	logMsgStatusError
	logMsgStatusErrorf
)

type logMsg struct {
	Status  logMsgStatus
	Content string
	Time    time.Time
	Data    []interface{}
}

func createPreContext(command *Command, ast *ASTNode) *Context {
	ctx, cancel := context.WithCancel(context.Background())

	return &Context{
		ctx:       ctx,
		cancel:    cancel,
		command:   command,
		ast:       ast,
		memory:    &memory,
		startTime: time.Now(),
	}
}

// Context - the structure that is passed when executing the functions of the command.
type Context struct {
	ctx        context.Context
	cancel     context.CancelFunc
	command    *Command
	ast        *ASTNode
	memory     *map[string]interface{}
	emitLog    func(logMsg)
	emitTUI    func(TUIRequest)
	emitTUICLI func(TUIRequest, chan<- bool)
	stdout     io.Writer
	stderr     io.Writer
	startTime  time.Time
	isCLI      bool
}

// GetName - returns the name of the command.
func (c *Context) GetName() string {
	return c.command.Name
}

// GetCommandNameTree is a method for getting a command tree.
func (c *Context) GetCommandNameTree() []string {
	return c.ast.CommandTree
}

// GetFlagInt is a method for getting a flag int value.
func (c *Context) GetFlagInt(name string, defaultValue int) int {
	return c.command.Flags.GetFlagInt(name, defaultValue)
}

// GetFlagString is a method for getting a flag string value.
func (c *Context) GetFlagString(name string, defaultValue string) string {
	return c.command.Flags.GetFlagString(name, defaultValue)
}

// GetFlagIntArray is a method for getting a flag int array value.
func (c *Context) GetFlagIntArray(name string) []int {
	return c.command.Flags.GetFlagIntArray(name)
}

// GetFlagStringArray is a method for getting a flag string array value.
func (c *Context) GetFlagStringArray(name string) []string {
	return c.command.Flags.GetFlagStringArray(name)
}

// GetFlagBool is a method for getting a flag bool value.
func (c *Context) GetFlagBool(name string) bool {
	return c.command.Flags.GetFlagBool(name)
}

// Print is a method for printing a message.
func (c *Context) Print(data ...interface{}) {
	c.emitLog(logMsg{
		Status:  logMsgStatusPrint,
		Content: fmt.Sprint(data...),
		Time:    time.Now(),
	})
}

// Printf is a method for printing a formatted message.
func (c *Context) Printf(format string, data ...interface{}) {
	c.emitLog(logMsg{
		Status:  logMsgStatusPrintf,
		Content: format,
		Time:    time.Now(),
		Data:    data,
	})
}

// PrintMarkdown is a method for printing a markdown message.
func (c *Context) PrintMarkdown(markdown string, data ...interface{}) {
	c.emitLog(logMsg{
		Status:  logMsgStatusPrintMarkdown,
		Content: markdown,
		Time:    time.Now(),
		Data:    data,
	})
}

// Warn is a method for printing a warning message.
func (c *Context) Warn(data ...interface{}) {
	c.emitLog(logMsg{
		Status:  logMsgStatusWarn,
		Content: fmt.Sprint(data...),
		Time:    time.Now(),
	})
}

// Warnf is a method for printing a formatted warning message.
func (c *Context) Warnf(format string, data ...interface{}) {
	c.emitLog(logMsg{
		Status:  logMsgStatusWarnf,
		Content: format,
		Time:    time.Now(),
		Data:    data,
	})
}

// Error is a method for printing an error message.
func (c *Context) Error(data ...interface{}) {
	c.emitLog(logMsg{
		Status:  logMsgStatusError,
		Content: fmt.Sprint(data...),
		Time:    time.Now(),
	})
}

// Errorf is a method for printing a formatted error message.
func (c *Context) Errorf(format string, data ...interface{}) {
	c.emitLog(logMsg{
		Status:  logMsgStatusErrorf,
		Content: format,
		Time:    time.Now(),
		Data:    data,
	})
}

// StartTime is a method for getting the start time of the command.
func (c *Context) StartTime() time.Time {
	return c.startTime
}

// Elapsed is a method for getting the elapsed time of the command.
func (c *Context) Elapsed() time.Duration {
	return time.Since(c.startTime)
}

// Command is a method for getting the command string.
func (c *Context) Command() string {
	return c.ast.FullCommand
}

// Stdout is a method for getting the stdout writer.
func (c *Context) Stdout() io.Writer {
	return c.stdout
}

// Stderr is a method for getting the stderr writer.
func (c *Context) Stderr() io.Writer {
	return c.stderr
}

// Ctx is a method for getting the context.Context.
func (c *Context) Ctx() context.Context {
	return c.ctx
}

// Done is a method for getting the done channel.
func (c *Context) Done() <-chan struct{} {
	return c.ctx.Done()
}

// IsCancelled is a method for checking if the context is cancelled.
func (c *Context) IsCancelled() bool {
	select {
	case <-c.ctx.Done():
		return true
	default:
		return false
	}
}

// Set is a method for setting a value in the memory.
func (c *Context) Set(key string, value interface{}) {
	if c.memory == nil {
		newMem := make(map[string]interface{})
		c.memory = &newMem
	}

	(*c.memory)[key] = value
}

// Delete is a method for deleting a value from the memory.
func (c *Context) Delete(key string) {
	if c.memory != nil {
		delete(*c.memory, key)
	}
}

// Get is a method for getting a value from the memory.
func (c *Context) Get(key string) interface{} {
	if c.memory != nil {
		return (*c.memory)[key]
	}

	return nil
}

// MustGetString is a method for getting a value from the memory and converting it to a string.
func (c *Context) MustGetString(key string) string {
	if c.memory != nil {
		if value, ok := (*c.memory)[key]; ok {
			if d, ok := value.(string); ok {
				return d
			}
		}
	}

	return ""
}

// MustGetInt is a method for getting a value from the memory and converting it to an int.
func (c *Context) MustGetInt(key string) int {
	if c.memory != nil {
		if value, ok := (*c.memory)[key]; ok {
			if d, ok := value.(int); ok {
				return d
			}
		}
	}

	return 0
}

// Exec is a method for executing a shell command.
func (c *Context) Exec(cmd string, args ...string) (string, string, error) {
	command := exec.CommandContext(c.ctx, cmd, args...)

	var stdout bytes.Buffer

	var stderr bytes.Buffer

	command.Stdout = &stdout
	command.Stderr = &stderr

	err := command.Run()

	return stdout.String(), stderr.String(), err
}

// ExecLive is a method for executing a shell command in a live environment.
func (c *Context) ExecLive(cmd string, args ...string) error {
	command := exec.CommandContext(c.ctx, cmd, args...)

	stdoutPipe, err := command.StdoutPipe()
	if err != nil {
		return err
	}

	stderrPipe, err := command.StderrPipe()
	if err != nil {
		return err
	}

	if err := command.Start(); err != nil {
		return err
	}

	go c.streamOutput(stdoutPipe, logMsgStatusPrint)
	go c.streamOutput(stderrPipe, logMsgStatusError)

	return command.Wait()
}

// ExecSilent is a method for executing a shell command silently.
func (c *Context) ExecSilent(cmd string, args ...string) error {
	command := exec.CommandContext(c.ctx, cmd, args...)

	return command.Run()
}

// SelectOne is a method that triggers TUI to receive one item from the list from the user.
func (c *Context) SelectOne(p *TUISelectOneParams) (TUISelectOneResult, error) {
	req := TUIRequest{
		ID:       uuid.NewString(),
		Type:     tuiTypeSelectOne,
		Payload:  *p,
		Response: make(chan TUIResponse),
	}

	if c.isCLI {
		close := make(chan bool)

		go c.emitTUICLI(req, close)

		defer func() {
			// Wait for the CLI TUI to finish
			<-close
		}()
	} else {
		go c.emitTUI(req)
	}

	res := <-req.Response
	if res.Err != nil {
		return TUISelectOneResult{}, res.Err
	}

	return res.Value.(TUISelectOneResult), nil
}

// InputText is a method that triggers TUI to receive text from the user.
func (c *Context) InputText(p *TUIInputTextParams) (string, error) {
	req := TUIRequest{
		ID:       uuid.NewString(),
		Type:     tuiTypeInputText,
		Payload:  *p,
		Response: make(chan TUIResponse),
	}

	if c.isCLI {
		close := make(chan bool)

		go c.emitTUICLI(req, close)

		defer func() {
			// Wait for the CLI TUI to finish
			<-close
		}()
	} else {
		go c.emitTUI(req)
	}

	res := <-req.Response
	if res.Err != nil {
		return "", res.Err
	}

	return res.Value.(string), nil
}

// InputInt is a method that triggers TUI to receive integer from the user.
func (c *Context) InputInt(p *TUIInputIntParams) (int, error) {
	req := TUIRequest{
		ID:       uuid.NewString(),
		Type:     tuiTypeInputInt,
		Payload:  *p,
		Response: make(chan TUIResponse),
	}

	if c.isCLI {
		close := make(chan bool)

		go c.emitTUICLI(req, close)

		defer func() {
			// Wait for the CLI TUI to finish
			<-close
		}()
	} else {
		go c.emitTUI(req)
	}

	res := <-req.Response
	if res.Err != nil {
		return 0, res.Err
	}

	return res.Value.(int), nil
}

// InputFile is a method that triggers TUI to receive file from the user.
func (c *Context) InputFile(p *TUIInputFileParams) (TUIInputFileResult, error) {
	req := TUIRequest{
		ID:       uuid.NewString(),
		Type:     tuiTypeInputFile,
		Payload:  *p,
		Response: make(chan TUIResponse),
	}

	if c.isCLI {
		close := make(chan bool)

		go c.emitTUICLI(req, close)

		defer func() {
			// Wait for the CLI TUI to finish
			<-close
		}()
	} else {
		go c.emitTUI(req)
	}

	res := <-req.Response
	if res.Err != nil {
		return TUIInputFileResult{}, res.Err
	}

	return res.Value.(TUIInputFileResult), nil
}

// Confirm is a method that triggers TUI to receive confirmation from the user.
func (c *Context) Confirm(p *TUIConfirmParams) (bool, error) {
	req := TUIRequest{
		ID:       uuid.NewString(),
		Type:     tuiTypeConfirm,
		Payload:  *p,
		Response: make(chan TUIResponse),
	}

	if c.isCLI {
		close := make(chan bool)

		go c.emitTUICLI(req, close)

		defer func() {
			// Wait for the CLI TUI to finish
			<-close
		}()
	} else {
		go c.emitTUI(req)
	}

	res := <-req.Response
	if res.Err != nil {
		return false, res.Err
	}

	return res.Value.(bool), nil
}

func (c *Context) streamOutput(r io.Reader, status logMsgStatus) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		if c.emitLog != nil {
			c.emitLog(logMsg{
				Status:  status,
				Content: scanner.Text(),
				Time:    time.Now(),
			})
		}
	}
}
