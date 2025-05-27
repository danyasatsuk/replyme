package replyme

// ASTFlag - is a structure for storing flags inside the AST.
type ASTFlag struct {
	// Flag type
	Type FlagType
	// Flag value
	Value string
}

// ASTArgument - is a structure for storing arguments inside the AST.
type ASTArgument struct {
	// Argument name
	Name string
	// Argument value
	Value string
}

// ASTNode - is a structure that defines the storage of a command after its successful parsing.
type ASTNode struct {
	Command      string
	FullCommand  string
	CommandTree  []string
	Subcommands  []string
	Arguments    []ASTArgument
	Flags        map[string]map[string][]ASTFlag
	Args         []string
	ColorCommand string
}
