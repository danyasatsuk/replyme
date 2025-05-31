package replyme

import (
	"fmt"
	"golang.org/x/exp/slices"
	"strconv"
	"strings"
	"unicode"
)

type FlagType uint16

const (
	FlagTypeInt FlagType = iota
	FlagTypeString
	FlagTypeIntArray
	FlagTypeStringArray
	FlagTypeBool
)

type flagSchema map[string]map[string]FlagType

type argsSchema map[string][]*Argument

type commandsSchema []commandSchema

type commandSchema struct {
	Name        string
	Subcommands []commandSchema
}

//nolint:gocognit,cyclop,funlen
func parseCommand(
	commands commandsSchema,
	schema flagSchema,
	argsSchema argsSchema,
	input string,
) (*ASTNode, error) {
	ast := &ASTNode{
		Flags:       map[string]map[string][]ASTFlag{},
		CommandTree: []string{},
		Subcommands: []string{},
	}

	tokens, err := tokenize(input)
	if err != nil {
		return nil, err
	}

	skip := -1

	var lastCmd string

	inArgs := false

	var posArgs []string

	var currentCmdSchema *commandSchema

	if len(tokens) == 0 {
		return nil, ErrorCommandEmpty
	}

	first := tokens[0]
	for i := range commands {
		if commands[i].Name == first {
			currentCmdSchema = &commands[i]

			break
		}
	}

	if currentCmdSchema == nil {
		return nil, fmt.Errorf("%w: %s", ErrorUnknownCommand, first)
	}

	ast.Command = first
	ast.CommandTree = append(ast.CommandTree, first)
	lastCmd = first

	for i := 1; i < len(tokens); i++ {
		token := tokens[i]

		if skip == i {
			continue
		}

		if token == "--" {
			inArgs = true

			continue
		}

		if inArgs {
			posArgs = append(posArgs, token)

			continue
		}

		if strings.HasPrefix(token, "-") { //nolint:nestif
			var name, value string

			if strings.Contains(token, "=") {
				parts := strings.SplitN(token, "=", 2) //nolint:mnd
				name = strings.TrimLeft(parts[0], "-")
				value = parts[1]
			} else {
				name = strings.TrimLeft(token, "-")
				flagType := schema[lastCmd][name]

				if flagType == 0 {
					flagType = schema["global"][name]
				}

				if flagType == FlagTypeBool || i+1 >= len(tokens) || strings.HasPrefix(tokens[i+1], "-") {
					value = "true"
				} else {
					value = tokens[i+1]
					skip = i + 1
				}
			}

			if ast.Flags[lastCmd] == nil {
				ast.Flags[lastCmd] = map[string][]ASTFlag{}
			}

			flagType := schema[lastCmd][name]
			if flagType == 0 {
				flagType = schema["global"][name]
			}

			ast.Flags[lastCmd][name] = append(ast.Flags[lastCmd][name], ASTFlag{Type: flagType, Value: value})
		} else {
			if len(currentCmdSchema.Subcommands) == 0 {
				posArgs = append(posArgs, token)

				continue
			}

			found := false

			for j := range currentCmdSchema.Subcommands {
				if currentCmdSchema.Subcommands[j].Name == token {
					currentCmdSchema = &currentCmdSchema.Subcommands[j]
					found = true

					break
				}
			}

			if !found {
				return nil, newErrorSubcommandUnknown(token)
			}

			lastCmd = token
			ast.CommandTree = append(ast.CommandTree, token)
			ast.Subcommands = append(ast.Subcommands, token)
		}
	}

	expected := argsSchema[lastCmd]
	if len(posArgs) < len(expected) {
		return nil, newErrorArgumentNotFound(lastCmd)
	}

	for i, def := range expected {
		if i < len(posArgs) {
			ast.Arguments = append(ast.Arguments, ASTArgument{Name: def.Name, Value: posArgs[i]})
		}
	}

	ast.Args = posArgs
	ast.FullCommand = input

	return ast, nil
}

//nolint:cyclop
func tokenize(input string) ([]string, error) {
	var result []string

	var current strings.Builder

	var inQuote bool

	var quoteChar rune

	var escape bool

	for _, r := range input {
		switch {
		case escape:
			current.WriteRune(r)

			escape = false

		case r == '\\':
			escape = true

		case r == '"' || r == '\'':
			if inQuote {
				if r == quoteChar {
					inQuote = false
				} else {
					current.WriteRune(r)
				}
			} else {
				inQuote = true
				quoteChar = r
			}

		case unicode.IsSpace(r):
			if inQuote {
				current.WriteRune(r)
			} else if current.Len() > 0 {
				result = append(result, current.String())
				current.Reset()
			}

		default:
			current.WriteRune(r)
		}
	}

	if inQuote {
		return nil, ErrorCommandUnclosedQuotes
	}

	if escape {
		return nil, ErrorIncompleteEscapeSequence
	}

	if current.Len() > 0 {
		result = append(result, current.String())
	}

	return result, nil
}

//nolint:cyclop
func createFlagSchema(commands Commands) flagSchema {
	schema := make(flagSchema)

	for _, command := range commands {
		for _, flag := range command.Flags {
			if schema[command.Name] == nil {
				schema[command.Name] = make(map[string]FlagType)
			}

			switch flag.ValueType() {
			case "string":
				schema[command.Name][flag.GetName()] = FlagTypeString
			case "int":
				schema[command.Name][flag.GetName()] = FlagTypeInt
			case "[]string":
				schema[command.Name][flag.GetName()] = FlagTypeStringArray
			case "[]int":
				schema[command.Name][flag.GetName()] = FlagTypeIntArray
			case "bool":
				schema[command.Name][flag.GetName()] = FlagTypeBool
			}
		}

		if command.Subcommands != nil && len(command.Subcommands) > 0 {
			newSchema := createFlagSchema(command.Subcommands)
			for k, v := range newSchema {
				schema[k] = v
			}
		}
	}

	return schema
}

func createArgsSchema(commands Commands) argsSchema {
	schema := make(argsSchema)
	for _, command := range commands {
		schema[command.Name] = command.Arguments

		if command.Subcommands != nil && len(command.Subcommands) > 0 {
			newSchema := createArgsSchema(command.Subcommands)
			for k, v := range newSchema {
				schema[k] = v
			}
		}
	}

	return schema
}

func createCommandSchema(commands Commands) commandsSchema {
	schema := make(commandsSchema, len(commands))
	if commands == nil || len(commands) == 0 {
		return schema
	}

	for _, command := range commands {
		schema = append(schema, commandSchema{
			Name:        command.Name,
			Subcommands: createCommandSchema(command.Subcommands),
		})
	}

	return schema
}

//nolint:cyclop
func insertDataInCommand(cmd *Command, ast *ASTNode, subcommand bool) error {
	if flags, ok := ast.Flags[cmd.Name]; ok { //nolint:nestif
		for _, cmdFlag := range cmd.Flags {
			if flag, ok := flags[cmdFlag.GetName()]; ok {
				_, err := cmdFlag.Parse(flag[0].Value)
				if err != nil {
					return err
				}
			}

			if flag, ok := flags[cmdFlag.GetAlias()]; ok {
				_, err := cmdFlag.Parse(flag[0].Value)
				if err != nil {
					return err
				}
			}
		}
	}

	if !subcommand {
		if cmd.Arguments != nil && len(cmd.Arguments) > 0 {
			for _, argument := range cmd.Arguments {
				i := slices.IndexFunc(ast.Arguments, func(a ASTArgument) bool {
					return a.Name == argument.Name
				})
				if i == -1 {
					return newErrorArgumentNotFound(argument.Name)
				}

				cmd.Arguments[i].value = ast.Arguments[i].Value
			}
		}
	}

	return nil
}

//nolint:cyclop,funlen
func colorCommand(input string) string {
	var result strings.Builder

	inQuote := false

	quoteChar := byte(0)

	current := strings.Builder{}

	isFirstToken := true

	for i := 0; i < len(input); i++ {
		ch := input[i]

		if (ch == '"' || ch == '\'') && !inQuote {
			inQuote = true
			quoteChar = ch
			current.WriteByte(ch)

			continue
		} else if inQuote && ch == quoteChar {
			inQuote = false

			current.WriteByte(ch)
			result.WriteString(styles.CMDStringStyle(current.String()))
			current.Reset()

			continue
		}

		if inQuote {
			current.WriteByte(ch)

			continue
		}

		if ch == ' ' { //nolint:nestif
			token := current.String()
			if token != "" {
				if isFirstToken {
					result.WriteString(styles.CMDCommandStyle(token))

					isFirstToken = false
				} else {
					result.WriteString(styleToken(token))
				}

				result.WriteByte(' ')
				current.Reset()
			} else {
				result.WriteByte(' ')
			}

			continue
		}

		current.WriteByte(ch)
	}

	if current.Len() > 0 {
		token := current.String()
		if isFirstToken {
			result.WriteString(styles.CMDCommandStyle(token))
		} else {
			result.WriteString(styleToken(token))
		}
	}

	return result.String()
}

func styleToken(token string) string {
	switch {
	case strings.HasPrefix(token, "--") || strings.HasPrefix(token, "-"):
		return styles.CMDFlagStyle(token)
	case isQuoted(token):
		return styles.CMDStringStyle(token)
	case isNumber(token):
		return styles.CMDFlagValueStyle(token)
	default:
		return styles.CMDArgValueStyle(token)
	}
}

func isQuoted(s string) bool {
	return (strings.HasPrefix(s, "\"") && strings.HasSuffix(s, "\"")) ||
		(strings.HasPrefix(s, "'") && strings.HasSuffix(s, "'"))
}

func isNumber(s string) bool {
	_, err := strconv.Atoi(s)

	return err == nil
}
