package replyme

import (
	"errors"
	"golang.org/x/exp/slices"
	"reflect"
	"strconv"
	"strings"
)

// Flag is an interface for getting information about flags and parsing it.
type Flag interface {
	// GetName returns the name of the flag.
	GetName() string
	// GetAlias returns the alias of the flag.
	GetAlias() string
	// ValueType returns the type of the value.
	ValueType() string
	// Value returns the value of the flag.
	Value() string
	// Parse parses the flag.
	Parse(flag string) (interface{}, error)
	// ParsedValue returns the parsed value of the flag.
	ParsedValue() (interface{}, error)
	// GetUsage returns the usage of the flag.
	GetUsage() string
	// Clear clears the flag.
	Clear()
}

// FlagValue is a structure for passing information about flags to a command.
type FlagValue[T any] struct {
	// Flag name
	Name string
	// Flag usage
	Usage string
	// Flag alias
	Alias string
	// Flag parser
	Parser         func(s string) (T, error)
	preParsedValue string
	value          T
	hasValue       bool
}

// GetUsage returns the usage of the flag.
func (f *FlagValue[T]) GetUsage() string {
	return f.Usage
}

// ParsedValue returns the parsed value of the flag.
func (f *FlagValue[T]) ParsedValue() (interface{}, error) {
	if !f.hasValue {
		return nil, errors.New("value is nil")
	}

	return f.value, nil
}

// GetName returns the name of the flag.
func (f *FlagValue[T]) GetName() string {
	return f.Name
}

// GetAlias returns the alias of the flag.
func (f *FlagValue[T]) GetAlias() string {
	return f.Alias
}

// ValueType returns the type of the value.
func (f *FlagValue[T]) ValueType() string {
	var zero T

	return reflect.TypeOf(zero).String()
}

// Value returns the value of the flag.
func (f *FlagValue[T]) Value() string {
	return f.preParsedValue
}

// Parse parses the flag.
func (f *FlagValue[T]) Parse(flag string) (interface{}, error) {
	var parsed T

	var err error

	if f.Parser != nil {
		parsed, err = f.Parser(flag)

		return parsed, err
	}

	switch any(parsed).(type) {
	case int:
		var d interface{}
		d, err = f.parseInt(flag)
		parsed = d.(T)
	case string:
		var d interface{}
		d = f.parseString(flag)
		parsed = d.(T)
	case []int:
		var d interface{}
		d, err = f.parseIntArray(flag)
		parsed = d.(T)
	case []string:
		var d interface{}
		d = f.parseStringArray(flag)
		parsed = d.(T)
	case bool:
		var d interface{}
		d, err = f.parseBool(flag)
		parsed = d.(T)
	default:
		return nil, newErrorUnknownFlagType(reflect.TypeOf(parsed).String())
	}

	f.value = parsed
	f.hasValue = true

	return parsed, err
}

func (f *FlagValue[T]) parseInt(flag string) (v interface{}, err error) {
	v, err = strconv.Atoi(flag)
	if err != nil {
		return 0, err
	}

	return
}

func (f *FlagValue[T]) parseString(flag string) (v interface{}) {
	return flag
}

func (f *FlagValue[T]) parseIntArray(flag string) (interface{}, error) {
	arr := []int{}

	parts := strings.Split(flag, ",")

	for _, part := range parts {
		n, convErr := strconv.Atoi(strings.TrimSpace(part))
		if convErr != nil {
			return nil, convErr
		}

		arr = append(arr, n)
	}

	return arr, nil
}

func (f *FlagValue[T]) parseStringArray(flag string) interface{} {
	var arr []string

	arr = strings.Split(flag, ",")
	for i := range arr {
		arr[i] = strings.TrimSpace(arr[i])
	}

	return arr
}

func (f *FlagValue[T]) parseBool(flag string) (v interface{}, err error) {
	return flag == "true", nil
}

// Clear clears the flag.
func (f *FlagValue[T]) Clear() {
	f.hasValue = false
	f.value = *new(T)
	f.preParsedValue = ""
}

// Flags is an abbreviation for the type `[]Flag`, which adds additional methods for convenient management.
type Flags []Flag

// GetFlagInt returns the value of the flag with the specified name as an int.
func (f Flags) GetFlagInt(name string, defaultValue int) int {
	i := slices.IndexFunc(f, func(flag Flag) bool {
		return flag.GetName() == name && flag.ValueType() == "int"
	})
	if i == -1 {
		return defaultValue
	}

	p, err := f[i].ParsedValue()
	if err != nil {
		return defaultValue
	}

	return p.(int)
}

// GetFlagString returns the value of the flag with the specified name as a string.
func (f Flags) GetFlagString(name string, defaultValue string) string {
	i := slices.IndexFunc(f, func(flag Flag) bool {
		return flag.GetName() == name && flag.ValueType() == "string"
	})
	if i == -1 {
		return defaultValue
	}

	p, err := f[i].ParsedValue()
	if err != nil {
		return defaultValue
	}

	return p.(string)
}

// GetFlagIntArray returns the value of the flag with the specified name as an array of ints.
func (f Flags) GetFlagIntArray(name string) []int {
	i := slices.IndexFunc(f, func(flag Flag) bool {
		return flag.GetName() == name && flag.ValueType() == "[]int"
	})
	if i == -1 {
		return []int{}
	}

	p, err := f[i].ParsedValue()
	if err != nil {
		return []int{}
	}

	return p.([]int)
}

// GetFlagStringArray returns the value of the flag with the specified name as an array of strings.
func (f Flags) GetFlagStringArray(name string) []string {
	i := slices.IndexFunc(f, func(flag Flag) bool {
		return flag.GetName() == name && flag.ValueType() == "[]string"
	})
	if i == -1 {
		return []string{}
	}

	p, err := f[i].ParsedValue()
	if err != nil {
		return []string{}
	}

	return p.([]string)
}

// GetFlagBool returns the value of the flag with the specified name as a bool.
func (f Flags) GetFlagBool(name string) bool {
	i := slices.IndexFunc(f, func(flag Flag) bool {
		return flag.GetName() == name && flag.ValueType() == "bool"
	})
	if i == -1 {
		return false
	}

	p, err := f[i].ParsedValue()
	if err != nil {
		return false
	}

	return p.(bool)
}
