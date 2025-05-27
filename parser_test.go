package replyme

import (
	"errors"
	"reflect"
	"testing"
)

func TestParseCommand_BasicCommand(t *testing.T) {
	input := `deploy`
	commands := CommandsSchema{{Name: "deploy"}}
	schema := FlagSchema{}
	argsSchema := ArgsSchema{}

	ast, err := ParseCommand(commands, schema, argsSchema, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ast.Command != "deploy" {
		t.Errorf("expected command 'deploy', got '%s'", ast.Command)
	}
}

func TestParseCommand_WithFlagAndArgs(t *testing.T) {
	input := `build --optimize=true main.go`
	commands := CommandsSchema{{
		Name: "build",
	}}
	schema := FlagSchema{
		"build": {
			"optimize": FlagTypeBool,
		},
	}
	argsSchema := ArgsSchema{
		"build": {
			{Name: "input"},
		},
	}

	ast, err := ParseCommand(commands, schema, argsSchema, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ast.Flags["build"]["optimize"][0].Value != "true" {
		t.Errorf("expected flag value 'true', got '%v'", ast.Flags["build"]["optimize"][0].Value)
	}
	if len(ast.Arguments) != 1 || ast.Arguments[0].Value != "main.go" {
		t.Errorf("expected argument 'main.go', got '%v'", ast.Arguments)
	}
}

func TestParseCommand_WithSubcommands(t *testing.T) {
	input := `db insert users.json`
	commands := CommandsSchema{{
		Name: "db",
		Subcommands: []CommandSchema{{
			Name: "insert",
		}},
	}}
	schema := FlagSchema{}
	argsSchema := ArgsSchema{
		"insert": {
			{Name: "file"},
		},
	}

	ast, err := ParseCommand(commands, schema, argsSchema, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(ast.CommandTree, []string{"db", "insert"}) {
		t.Errorf("expected command tree ['db', 'insert'], got %v", ast.CommandTree)
	}
	if len(ast.Arguments) != 1 || ast.Arguments[0].Value != "users.json" {
		t.Errorf("expected argument 'users.json', got '%v'", ast.Arguments)
	}
}

func TestParseCommand_UnknownCommand(t *testing.T) {
	_, err := ParseCommand(
		CommandsSchema{{Name: "known"}},
		FlagSchema{},
		ArgsSchema{},
		"unknowncmd",
	)
	if err == nil || err.Error() != "unknown subcommand: unknowncmd" && !errors.Is(err, ErrorUnknownCommand) {
		t.Errorf("expected unknown command error, got %v", err)
	}
}

func TestTokenize(t *testing.T) {
	input := `cmd --flag="some string" --int=42 value1 value2`
	tokens, err := tokenize(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []string{"cmd", "--flag=some string", "--int=42", "value1", "value2"}
	if !reflect.DeepEqual(tokens, expected) {
		t.Errorf("expected tokens %v, got %v", expected, tokens)
	}
}
