package main

import "github.com/danyasatsuk/replyme"

func main() {
	app := &replyme.App{
		Name:  "flow",
		Usage: "Flow example",
		Commands: []*replyme.Command{
			{
				Name:  "mainCommand",
				Usage: "Main command",
				Flags: []replyme.Flag{
					&replyme.FlagValue[string]{
						Name:  "testFlag",
						Usage: "Test flag",
					},
				},
				Before: func(ctx *replyme.Context) (bool, error) {
					ctx.Printf("This function will be launched very first. --testFlag %s", ctx.GetFlagString("testFlag", "unknown"))

					return true, nil
				},
				Action: func(ctx *replyme.Context) error {
					ctx.Print("This function will be started if no subcommands are specified.")

					return nil
				},
				OnEnd: func(ctx *replyme.Context) error {
					ctx.Print("This function will be launched most recently.")

					return nil
				},
				Subcommands: []*replyme.Command{
					{
						Name:  "subCommand",
						Usage: "Sub command",
						Before: func(ctx *replyme.Context) (bool, error) {
							ctx.Print("This function will be launched next after the main command.")

							return true, nil
						},
						Action: func(ctx *replyme.Context) error {
							ctx.Print("This function will be started if no subcommands are specified.")

							return nil
						},
						OnEnd: func(ctx *replyme.Context) error {
							ctx.Print("This function will be launched after the subcommand has completed all its actions.")

							return nil
						},
						Flags: []replyme.Flag{
							&replyme.FlagValue[string]{
								Name:  "testSubFlag",
								Usage: "Test sub flag",
							},
						},
						Subcommands: []*replyme.Command{
							{
								Name:  "subSubCommand",
								Usage: "Sub sub command",
								Before: func(ctx *replyme.Context) (bool, error) {
									ctx.Print("This function will start before executing the Action.")

									return true, nil
								},
								Action: func(ctx *replyme.Context) error {
									ctx.Print("This function runs the main code of the subcommand.")

									return nil
								},
								OnEnd: func(ctx *replyme.Context) error {
									ctx.Print("This function will start after executing the Action.")

									return nil
								},
								Flags: []replyme.Flag{
									&replyme.FlagValue[string]{
										Name:  "testSubSubFlag",
										Usage: "Test sub/sub flag",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	err := replyme.Run(app)
	if err != nil {
		panic(err)
	}
}
