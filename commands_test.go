package replyme

import (
	"github.com/go-faker/faker/v4"
	"testing"
)

func TestCommands_GetCommand(t *testing.T) {
	testCmd := faker.Word()
	cmd := Commands{
		{
			Name:  faker.Word(),
			Usage: faker.Sentence(),
		},
		{
			Name:  testCmd,
			Usage: faker.Sentence(),
		},
		{
			Name:  faker.Word(),
			Usage: faker.Sentence(),
		},
	}
	command, err := cmd.getCommand(testCmd)
	if err != nil {
		t.Fatal(err)
	}
	if command.Name != testCmd {
		t.Errorf("got command name %q, want %q", command.Name, testCmd)
	}
}

func TestSubber(t *testing.T) {
	testCmd := faker.Word()
	cmd := &Command{
		Name:  faker.Word(),
		Usage: faker.Sentence(),
		Subcommands: Commands{
			{
				Name:  faker.Word(),
				Usage: faker.Sentence(),
				Subcommands: Commands{
					{
						Name:  faker.Word(),
						Usage: faker.Sentence(),
					},
					{
						Name:  faker.Word(),
						Usage: faker.Sentence(),
						Subcommands: Commands{
							{
								Name:  faker.Word(),
								Usage: faker.Sentence(),
							},
						},
					},
				},
			},
			{
				Name:  testCmd,
				Usage: faker.Sentence(),
			},
			{
				Name:  faker.Word(),
				Usage: faker.Sentence(),
			},
		},
	}

	arr := subber(cmd)
	if len(arr) != 6 {
		t.Errorf("got subber length %d, want 6", len(arr))
	}
}

func TestCommands_GetCommandsArray(t *testing.T) {
	testCmd := faker.Word()
	cmd := Commands{
		&Command{
			Name:  faker.Word(),
			Usage: faker.Sentence(),
			Subcommands: Commands{
				{
					Name:  faker.Word(),
					Usage: faker.Sentence(),
					Subcommands: Commands{
						{
							Name:  faker.Word(),
							Usage: faker.Sentence(),
						},
						{
							Name:  faker.Word(),
							Usage: faker.Sentence(),
							Subcommands: Commands{
								{
									Name:  faker.Word(),
									Usage: faker.Sentence(),
								},
							},
						},
					},
				},
				{
					Name:  testCmd,
					Usage: faker.Sentence(),
				},
				{
					Name:  faker.Word(),
					Usage: faker.Sentence(),
				},
			},
		},
		&Command{
			Name:  faker.Word(),
			Usage: faker.Sentence(),
			Subcommands: Commands{
				{
					Name:  faker.Word(),
					Usage: faker.Sentence(),
					Subcommands: Commands{
						{
							Name:  faker.Word(),
							Usage: faker.Sentence(),
						},
						{
							Name:  faker.Word(),
							Usage: faker.Sentence(),
							Subcommands: Commands{
								{
									Name:  faker.Word(),
									Usage: faker.Sentence(),
								},
							},
						},
					},
				},
				{
					Name:  testCmd,
					Usage: faker.Sentence(),
				},
				{
					Name:  faker.Word(),
					Usage: faker.Sentence(),
				},
			},
		},
		&Command{
			Name:  faker.Word(),
			Usage: faker.Sentence(),
			Subcommands: Commands{
				{
					Name:  faker.Word(),
					Usage: faker.Sentence(),
					Subcommands: Commands{
						{
							Name:  faker.Word(),
							Usage: faker.Sentence(),
						},
						{
							Name:  faker.Word(),
							Usage: faker.Sentence(),
							Subcommands: Commands{
								{
									Name:  faker.Word(),
									Usage: faker.Sentence(),
								},
							},
						},
					},
				},
				{
					Name:  testCmd,
					Usage: faker.Sentence(),
				},
				{
					Name:  faker.Word(),
					Usage: faker.Sentence(),
				},
			},
		},
	}

	arr := cmd.getCommandsArray()
	if len(arr) != 21 {
		t.Errorf("got subber length %d, want 21", len(arr))
	}
}

func TestCommands_MustGetCommand(t *testing.T) {
	testCmd := faker.Word()
	cmd := Commands{
		{
			Name:  faker.Word(),
			Usage: faker.Sentence(),
		},
		{
			Name:  testCmd,
			Usage: faker.Sentence(),
		},
		{
			Name:  faker.Word(),
			Usage: faker.Sentence(),
		},
	}
	command := cmd.mustGetCommand(testCmd)
	if command.Name != testCmd {
		t.Errorf("got command name %q, want %q", command.Name, testCmd)
	}
}
