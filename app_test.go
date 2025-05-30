package replyme

import (
	"github.com/go-faker/faker/v4"
	"reflect"
	"testing"
)

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
	//t.Log(cmdLen, flagsLen)
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
