package replyme

// Argument is a structure that describes the arguments for your command.
type Argument struct {
	Name  string
	Usage string
	value string
}

// GetValue - method for getting the value of an argument.
func (arg *Argument) GetValue() string {
	return arg.value
}

// SetValue - method for setting the value of an argument.
func (arg *Argument) SetValue(v string) {
	arg.value = v
}
