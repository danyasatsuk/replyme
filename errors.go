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
