package replyme

import (
	"github.com/go-faker/faker/v4"
	"reflect"
	"slices"
	"testing"
)

//nolint:cyclop
func TestParseFlagSchema(t *testing.T) {
	commands := make(Commands, 200)
	for i := range commands {
		flags := make(Flags, 200)
		for i := range flags {
			var flag Flag

			n, err := faker.RandomInt(1, 5)
			if err != nil {
				t.Fatal(err)
			}

			switch n[0] {
			case 1:
				flag = &FlagValue[string]{
					Name: faker.UUIDHyphenated(),
				}
			case 2:
				flag = &FlagValue[int]{
					Name: faker.UUIDHyphenated(),
				}
			case 3:
				flag = &FlagValue[bool]{
					Name: faker.UUIDHyphenated(),
				}
			case 4:
				flag = &FlagValue[[]string]{
					Name: faker.UUIDHyphenated(),
				}
			case 5:
				flag = &FlagValue[[]int]{
					Name: faker.UUIDHyphenated(),
				}
			}

			flags[i] = flag
		}

		command := &Command{
			Name:  faker.UUIDHyphenated(),
			Usage: faker.Sentence(),
			Flags: flags,
		}
		commands[i] = command
	}

	schema := parseFlagSchema(commands)
	cmdLen := 0
	flagsLen := 0

	for _, v := range schema {
		cmdLen++

		for range v {
			flagsLen++
		}
	}

	if cmdLen != 200 || flagsLen != 40000 {
		t.Fatal("parse flag schema failed")
	}
}

func TestParseFlagSchemaSingle(t *testing.T) {

	command := &Command{
		Name:  "test",
		Usage: faker.Sentence(),
		Flags: []Flag{
			&FlagValue[string]{
				Name: "a",
			},
			&FlagValue[int]{
				Name: "b",
			},
			&FlagValue[bool]{
				Name: "c",
			},
			&FlagValue[[]int]{
				Name: "d",
			},
			&FlagValue[[]string]{
				Name: "e",
			},
		},
	}
	value := flagSchema{
		"test": {
			"a": FlagTypeString,
			"b": FlagTypeInt,
			"c": FlagTypeBool,
			"d": FlagTypeIntArray,
			"e": FlagTypeStringArray,
		},
	}

	schema := parseFlagSchemaSingle(command)

	if !reflect.DeepEqual(value, schema) {
		t.Fatal("parse flag schema failed", value, schema)
	}
}

func TestApp_ParseFlagSchema(t *testing.T) {
	app := &App{
		Commands: Commands{
			&Command{
				Name:  "test",
				Usage: "test command",
				Flags: Flags{
					&FlagValue[string]{
						Name: "a",
					},
					&FlagValue[int]{
						Name: "b",
					},
				},
				Subcommands: Commands{
					&Command{
						Name:  "sub",
						Usage: "sub command",
						Flags: Flags{
							&FlagValue[bool]{
								Name: "c",
							},
						},
					},
				},
			},
		},
	}

	schema := app.getFlagSchema()

	expected := flagSchema{
		"test": {
			"a": FlagTypeString,
			"b": FlagTypeInt,
		},
		"sub": {
			"c": FlagTypeBool,
		},
	}

	if !reflect.DeepEqual(schema, expected) {
		t.Fatal("parse flag schema failed", schema, expected)
	}
}

func TestApp_SetHelpFlags(t *testing.T) {
	i18nInit()

	app := &App{
		Commands: Commands{
			&Command{
				Name:  "test",
				Usage: "test command",
				Flags: Flags{
					&FlagValue[string]{
						Name: "a",
					},
					&FlagValue[int]{
						Name: "b",
					},
				},
				Subcommands: Commands{
					&Command{
						Name:  "sub",
						Usage: "sub command",
						Flags: Flags{
							&FlagValue[bool]{
								Name: "c",
							},
						},
					},
				},
			},
		},
	}

	app.setHelpFlags()

	setHelpFlagsChecker(t, app.Commands)
}

func setHelpFlagsChecker(t *testing.T, commands Commands) {
	for _, cmd := range commands {
		if slices.IndexFunc(cmd.Flags, func(flag Flag) bool {
			return flag.GetName() == "help"
		}) == -1 {
			t.Fatal("help flag not set")
		}
		if cmd.Subcommands != nil {
			setHelpFlagsChecker(t, cmd.Subcommands)
		}
	}
}
