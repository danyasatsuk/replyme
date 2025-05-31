package replyme

import (
	"fmt"
	"github.com/charmbracelet/glamour"
	"slices"
	"strings"
	"time"
)

func (m *model) renderMarkdown(content string) string {
	renderer, err := glamour.NewTermRenderer(glamour.WithAutoStyle(), glamour.WithWordWrap(m.windowHeight))
	if err != nil {
		m.logsChan <- log{logTypeError, m.runningCommand, fmt.Sprintf("error creating renderer: %v", err), nil, time.Now()}
	}

	render, err := renderer.Render(content)
	if err != nil {
		m.logsChan <- log{logTypeError, m.runningCommand, fmt.Sprintf("error render: %v", err), nil, time.Now()}
	}

	return render
}

func (m *model) emitLog(l logMsg) {
	if l.Data == nil {
		l.Data = []interface{}{}
	}

	switch l.Status {
	case logMsgStatusPrint:
		m.logsChan <- log{logTypeLog, m.runningCommand, l.Content, nil, time.Now()}
	case logMsgStatusPrintf:
		m.logsChan <- log{logTypeLog, m.runningCommand, fmt.Sprintf(l.Content, l.Data...), nil, time.Now()}
	case logMsgStatusPrintMarkdown:
		m.logsChan <- log{logTypeLog, m.runningCommand, m.renderMarkdown(l.Content), nil, time.Now()}
	case logMsgStatusWarn:
		m.logsChan <- log{logTypeWarn, m.runningCommand, l.Content, nil, time.Now()}
	case logMsgStatusWarnf:
		m.logsChan <- log{logTypeWarn, m.runningCommand, fmt.Sprintf(l.Content, l.Data...), nil, time.Now()}
	case logMsgStatusError:
		m.logsChan <- log{logTypeError, m.runningCommand, l.Content, nil, time.Now()}
	case logMsgStatusErrorf:
		m.logsChan <- log{logTypeError, m.runningCommand, fmt.Sprintf(l.Content, l.Data...), nil, time.Now()}
	}
}

// LogType is the type of log.
type logType uint16

const (
	logTypeCommandRunning logType = iota
	logTypeCommandSuccess
	logTypeCommandFailure
	logTypeCommandNotFound
	logTypeCommandNotEnoughArguments

	logTypeMessage
	logTypePanic
	logTypeLog
	logTypeDebug
	logTypeWarn
	logTypeError
)

type log struct {
	Type    logType
	Command string
	Message string
	Error   error
	Time    time.Time
}

func renderRunning(s string) string {
	return fmt.Sprintf("⏳ %s %s", styles.GrayStyle(">>"), s)
}

func renderSuccess(s string) string {
	return fmt.Sprintf("%s %s %s", greenIcon.Render("✔"), styles.GrayStyle(">>"), s)
}

func renderFailure(s string) string {
	return fmt.Sprintf("%s %s %s", redIcon.Render("✖"), styles.GrayStyle(">>"), s)
}

func renderPanic(s string) string {
	return fmt.Sprintf("%s: %s", styles.ErrorHeaderStyle("[PANIC]"), styles.ErrorTextStyle(s))
}

func renderLog(s string) string {
	return fmt.Sprintf("%s: %s", styles.LogStyle("[LOG]"), s)
}

func renderDebug(s string) string {
	return fmt.Sprintf("%s: %s", styles.DebugStyle("[DEBUG]"), s)
}

func renderWarn(s string) string {
	return fmt.Sprintf("%s: %s", styles.WarnStyle("[WARN]"), s)
}

func renderError(s string) string {
	return fmt.Sprintf("%s: %s", styles.ErrorHeaderStyle("[ERROR]"), styles.ErrorTextStyle(s))
}

func renderCommandError(cmd string, s string) string {
	return fmt.Sprintf("%s %s %s %s", redIcon.Render("✖"), styles.GrayStyle(">>"), cmd, styles.GrayStyle("("+s+")"))
}

// Render - renders the log.
//
//nolint:cyclop
func (l log) Render() string {
	switch l.Type {
	case logTypeCommandRunning:
		return renderRunning(l.Command)
	case logTypeCommandSuccess:
		return renderSuccess(l.Command)
	case logTypeCommandFailure:
		return renderFailure(l.Command)
	case logTypeCommandNotFound, logTypeCommandNotEnoughArguments:
		return renderCommandError(l.Command, l.Message)
	case logTypeMessage:
		return l.Message
	case logTypeLog:
		return renderLog(l.Message)
	case logTypeDebug:
		return renderDebug(l.Message)
	case logTypeWarn:
		return renderWarn(l.Message)
	case logTypeError:
		return renderError(l.Message)
	case logTypePanic:
		return renderPanic(l.Message)
	default:
		return l.Message
	}
}

type logs []log

// Render - renders the logs.
func (l *logs) Render() string {
	var b strings.Builder
	for _, log := range *l {
		b.WriteString(log.Render() + "\n")
	}

	return b.String()
}

// RenderLimit - renders the logs with line limit.
func (l *logs) RenderLimit(limit int) string {
	var b strings.Builder

	if len(*l) < limit {
		limit = len(*l)
	}

	for _, log := range (*l)[len(*l)-1-limit:] {
		b.WriteString(log.Render() + "\n")
	}

	return b.String()
}

// RenderLimitFrom - renders the logs with line limit from a given index.
func (l *logs) RenderLimitFrom(from, limit int) string {
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

func (l *logs) Add(t logType, message string) {
	*l = append(*l, log{
		Type:    t,
		Command: message,
		Message: message,
		Time:    time.Now(),
	})
}

// AddLog - adds a log.
func (l *logs) AddLog(lg log) {
	if lg.Type == logTypeCommandSuccess || lg.Type == logTypeCommandFailure ||
		lg.Type == logTypeCommandNotFound || lg.Type == logTypeCommandNotEnoughArguments {
		i := slices.IndexFunc(*l, func(l log) bool {
			return lg.Command == l.Command && l.Type == logTypeCommandRunning
		})
		if i == -1 {
			return
		}

		(*l)[i] = lg

		return
	}

	*l = append(*l, lg)
}
