package replyme

import (
	"errors"
	"fmt"
	"golang.org/x/exp/slices"
	"strconv"
	"strings"
	"unicode"
)

// FlagType - a specific type of flags
type FlagType uint16

const (
	FlagTypeInt FlagType = iota
	FlagTypeString
	FlagTypeIntArray
	FlagTypeStringArray
	FlagTypeBool
)

// FlagSchema - a schema of flags
type FlagSchema map[string]map[string]FlagType

// ArgsSchema - a schema of arguments
type ArgsSchema map[string][]*Argument

// CommandsSchema - a schema of commands
type CommandsSchema []CommandSchema

// CommandSchema - a schema of command
type CommandSchema struct {
	Name        string
	Subcommands []CommandSchema
}

// ParseCommand is a function for parsing commands in ASTNode
func ParseCommand(
	commands CommandsSchema,
	schema FlagSchema,
	argsSchema ArgsSchema,
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
	var currentCmdSchema *CommandSchema

	// 🔍 Поиск начальной команды
	if len(tokens) == 0 {
		return nil, errors.New("пустая команда")
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
		if strings.HasPrefix(token, "-") {
			var name, value string
			if strings.Contains(token, "=") {
				parts := strings.SplitN(token, "=", 2)
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
			// Команда
			// Если у команды нет подкоманд — всё, что дальше, это аргументы
			if len(currentCmdSchema.Subcommands) == 0 {
				posArgs = append(posArgs, token)
				continue
			}

			// Иначе проверяем как подкоманду
			found := false
			for j := range currentCmdSchema.Subcommands {
				if currentCmdSchema.Subcommands[j].Name == token {
					currentCmdSchema = &currentCmdSchema.Subcommands[j]
					found = true
					break
				}
			}
			if !found {
				return nil, fmt.Errorf("неизвестная подкоманда: %s", token)
			}
			lastCmd = token
			ast.CommandTree = append(ast.CommandTree, token)
			ast.Subcommands = append(ast.Subcommands, token)
		}
	}

	// Позиционные аргументы
	expected := argsSchema[lastCmd]
	if len(posArgs) < len(expected) {
		return nil, NewErrorArgumentNotFound(lastCmd)
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
		return nil, fmt.Errorf("незакрытая кавычка")
	}
	if escape {
		return nil, fmt.Errorf("незавершённый escape-последовательность")
	}
	if current.Len() > 0 {
		result = append(result, current.String())
	}

	return result, nil
}

func createFlagSchema(commands Commands) FlagSchema {
	schema := make(FlagSchema)
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

func createArgsSchema(commands Commands) ArgsSchema {
	schema := make(ArgsSchema)
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

func createCommandSchema(commands Commands) CommandsSchema {
	schema := make(CommandsSchema, len(commands))
	if commands == nil || len(commands) == 0 {
		return schema
	}
	for _, command := range commands {
		schema = append(schema, CommandSchema{
			Name:        command.Name,
			Subcommands: createCommandSchema(command.Subcommands),
		})
	}
	return schema
}

func insertDataInCommand(cmd *Command, ast *ASTNode, subcommand bool) error {
	if flags, ok := ast.Flags[cmd.Name]; ok {
		for _, cmdFlag := range cmd.Flags {
			if flag, ok := flags[cmdFlag.GetName()]; ok {
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
					return NewErrorArgumentNotFound(argument.Name)
				}
				cmd.Arguments[i].value = ast.Arguments[i].Value
			}
		}
	}
	return nil
}

// ColorCommand - a function for creating a color command (not used yet, to be added later)
// TODO(unimportant): Add color command
func ColorCommand(input string) string {
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
			result.WriteString(CMDStringStyle(current.String()))
			current.Reset()
			continue
		}

		if inQuote {
			current.WriteByte(ch)
			continue
		}

		if ch == ' ' {
			token := current.String()
			if token != "" {
				if isFirstToken {
					result.WriteString(CMDCommandStyle(token))
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
			result.WriteString(CMDCommandStyle(token))
		} else {
			result.WriteString(styleToken(token))
		}
	}

	return result.String()
}

func styleToken(token string) string {
	switch {
	case strings.HasPrefix(token, "--") || strings.HasPrefix(token, "-"):
		return CMDFlagStyle(token)
	case isQuoted(token):
		return CMDStringStyle(token)
	case isNumber(token):
		return CMDFlagValueStyle(token)
	default:
		return CMDArgValueStyle(token)
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
