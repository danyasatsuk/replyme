package fmt

import (
	"fmt"
	"io"
	"os"
)

var stdout io.Writer = os.Stdout
var stderr io.Writer = os.Stderr

func SetStdout(out io.Writer) {
	stdout = out
}

func SetStderr(out io.Writer) {
	stderr = out
}

func Print(a ...interface{}) (n int, err error) {
	return fmt.Fprint(stdout, a...)
}

func Printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(stdout, format, a...)
}

func Println(a ...interface{}) (n int, err error) {
	return fmt.Fprintln(stdout, a...)
}

func Append(b []byte, a ...any) []byte {
	return fmt.Append(b, a...)
}

func Appendf(b []byte, format string, a ...any) []byte {
	return fmt.Appendf(b, format, a)
}

func Appendln(b []byte, a ...any) []byte {
	return fmt.Appendln(b, a...)
}

func Errorf(format string, a ...any) error {
	return fmt.Errorf(format, a...)
}

func FormatString(state fmt.State, verb rune) string {
	return fmt.FormatString(state, verb)
}

func Sprint(a ...any) string {
	return fmt.Sprint(a...)
}

func Sprintf(format string, a ...any) string {
	return fmt.Sprintf(format, a...)
}

func Sprintln(a ...any) string {
	return fmt.Sprintln(a...)
}
