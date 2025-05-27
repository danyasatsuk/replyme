package replyme

import (
	"errors"
	"fmt"
)

var ErrorUnknownCommand = errors.New("unknown command")

func NewErrorUnknownCommand(cmd string) error {
	return fmt.Errorf("%w: %s", ErrorUnknownCommand, cmd)
}

var ErrorArgumentNotFound = errors.New("argument not found for command")

func NewErrorArgumentNotFound(cmd string) error {
	return fmt.Errorf("%w: %s", ErrorArgumentNotFound, cmd)
}

var ErrorCommandPanic = errors.New("cmdpanic")

func NewErrorCommandPanic(cmd string) error {
	return fmt.Errorf("%w: %s", ErrorCommandPanic, cmd)
}

var ErrorUnknownFlagType = errors.New("unknown flag type")

func NewErrorUnknownFlagType(t string) error {
	return fmt.Errorf("%w: %s", ErrorUnknownFlagType, t)
}

var ErrorCommandEmpty = errors.New("command empty")

var ErrorSubcommandUnknown = errors.New("unknown subcommand")

func NewErrorSubcommandUnknown(cmd string) error {
	return fmt.Errorf("%w: %s", ErrorSubcommandUnknown, cmd)
}

var ErrorCommandUnclosedQuotes = errors.New("unclosed quotes")

var ErrorIncompleteEscapeSequence = errors.New("incomplete escape sequence")
