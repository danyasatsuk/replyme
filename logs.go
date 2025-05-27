package replyme

import (
	"fmt"
	"github.com/charmbracelet/glamour"
	"slices"
	"strings"
	"time"
)

func (m *Model) emitLog(log LogMsg) {
	if log.Data == nil {
		log.Data = []interface{}{}
	}
	switch log.Status {
	case LogMsgStatus_Print:
		m.logsChan <- Log{LogTypeLog, m.runningCommand, log.Content, nil, time.Now()}
	case LogMsgStatus_Printf:
		m.logsChan <- Log{LogTypeLog, m.runningCommand, fmt.Sprintf(log.Content, log.Data...), nil, time.Now()}
	case LogMsgStatus_PrintMarkdown:
		renderer, err := glamour.NewTermRenderer(glamour.WithAutoStyle(), glamour.WithWordWrap(m.windowHeight))
		if err != nil {
			m.logsChan <- Log{LogTypeError, m.runningCommand, fmt.Sprintf("error creating renderer: %v", err), nil, time.Now()}
		}
		render, err := renderer.Render(log.Content)
		if err != nil {
			m.logsChan <- Log{LogTypeError, m.runningCommand, fmt.Sprintf("error render: %v", err), nil, time.Now()}
		}
		m.logsChan <- Log{LogTypeLog, m.runningCommand, render, nil, time.Now()}
	case LogMsgStatus_Warn:
		m.logsChan <- Log{LogTypeWarn, m.runningCommand, log.Content, nil, time.Now()}
	case LogMsgStatus_Warnf:
		m.logsChan <- Log{LogTypeWarn, m.runningCommand, fmt.Sprintf(log.Content, log.Data...), nil, time.Now()}
	case LogMsgStatus_Error:
		m.logsChan <- Log{LogTypeError, m.runningCommand, log.Content, nil, time.Now()}
	case LogMsgStatus_Errorf:
		m.logsChan <- Log{LogTypeError, m.runningCommand, fmt.Sprintf(log.Content, log.Data...), nil, time.Now()}
	}
}

// LogType is the type of log
type LogType uint16

const (
	LogTypeCommandRunning LogType = iota
	LogTypeCommandSuccess
	LogTypeCommandFailure
	LogTypeCommandNotFound
	LogTypeCommandNotEnoughArguments

	LogTypeMessage
	LogTypePanic
	LogTypeLog
	LogTypeDebug
	LogTypeWarn
	LogTypeError
)

// Log - a structure for transferring and storing logs inside the REPL
type Log struct {
	Type    LogType
	Command string
	Message string
	Error   error
	Time    time.Time
}

func renderRunning(s string) string {
	return fmt.Sprintf("⏳ %s %s", GrayStyle(">>"), s)
}

func renderSuccess(s string) string {
	return fmt.Sprintf("%s %s %s", GreenIcon.Render("✔"), GrayStyle(">>"), s)
}

func renderFailure(s string) string {
	return fmt.Sprintf("%s %s %s", RedIcon.Render("✖"), GrayStyle(">>"), s)
}

func renderPanic(s string) string {
	return fmt.Sprintf("%s: %s", ErrorHeaderStyle("[PANIC]"), ErrorTextStyle(s))
}

func renderLog(s string) string {
	return fmt.Sprintf("%s: %s", LogStyle("[LOG]"), s)
}

func renderDebug(s string) string {
	return fmt.Sprintf("%s: %s", DebugStyle("[DEBUG]"), s)
}

func renderWarn(s string) string {
	return fmt.Sprintf("%s: %s", WarnStyle("[WARN]"), s)
}

func renderError(s string) string {
	return fmt.Sprintf("%s: %s", ErrorHeaderStyle("[ERROR]"), ErrorTextStyle(s))
}

func renderCommandError(cmd string, s string) string {
	return fmt.Sprintf("%s %s %s %s", RedIcon.Render("✖"), GrayStyle(">>"), cmd, GrayStyle("("+s+")"))
}

// Render - renders the log
func (l Log) Render() string {
	switch l.Type {
	case LogTypeCommandRunning:
		return renderRunning(l.Command)
	case LogTypeCommandSuccess:
		return renderSuccess(l.Command)
	case LogTypeCommandFailure:
		return renderFailure(l.Command)
	case LogTypeCommandNotFound, LogTypeCommandNotEnoughArguments:
		return renderCommandError(l.Command, l.Message)
	case LogTypeMessage:
		return l.Message
	case LogTypeLog:
		return renderLog(l.Message)
	case LogTypeDebug:
		return renderDebug(l.Message)
	case LogTypeWarn:
		return renderWarn(l.Message)
	case LogTypeError:
		return renderError(l.Message)
	case LogTypePanic:
		return renderPanic(l.Message)
	default:
		return l.Message
	}
}

// Logs - a slice of logs
type Logs []Log

// Render - renders the logs
func (l *Logs) Render() string {
	var b strings.Builder
	for _, log := range *l {
		b.WriteString(log.Render() + "\n")
	}
	return b.String()
}

// RenderLimit - renders the logs with line limit
func (l *Logs) RenderLimit(limit int) string {
	var b strings.Builder
	if len(*l) < limit {
		limit = len(*l)
	}
	for _, log := range (*l)[len(*l)-1-limit:] {
		b.WriteString(log.Render() + "\n")
	}
	return b.String()
}

// RenderLimitFrom - renders the logs with line limit from a given index
func (l *Logs) RenderLimitFrom(from, limit int) string {
	var b strings.Builder
	if from >= len(*l) {
		return ""
	}
	end := from + limit
	if end > len(*l) {
		end = len(*l)
	}
	for _, log := range (*l)[from:end] {
		b.WriteString(log.Render() + "\n")
	}
	return b.String()
}

func (l *Logs) Add(t LogType, message string) {
	*l = append(*l, Log{
		Type:    t,
		Command: message,
		Message: message,
		Time:    time.Now(),
	})
}

// AddLog - adds a log
func (l *Logs) AddLog(log Log) {
	if log.Type == LogTypeCommandSuccess || log.Type == LogTypeCommandFailure || log.Type == LogTypeCommandNotFound || log.Type == LogTypeCommandNotEnoughArguments {
		i := slices.IndexFunc(*l, func(l Log) bool {
			return log.Command == l.Command && l.Type == LogTypeCommandRunning
		})
		if i == -1 {
			return
		}
		(*l)[i] = log
		return
	}
	*l = append(*l, log)
}
