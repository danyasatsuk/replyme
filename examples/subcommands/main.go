package main

import "github.com/danyasatsuk/replyme"

func main() {
	app := &replyme.App{
		Name:  "subcommands",
		Usage: "My subcommands app",
		Commands: []*replyme.Command{
			{
				Name:  "auth",
				Usage: "Authenticates the user",
				Flags: []replyme.Flag{
					&replyme.FlagValue[string]{
						Name:  "server",
						Usage: "The server to authenticate to",
					},
				},
				Subcommands: []*replyme.Command{
					{
						Name:  "login",
						Usage: "Please login",
						Flags: []replyme.Flag{
							&replyme.FlagValue[string]{
								Name:  "username",
								Usage: "The username to login with",
							},
							&replyme.FlagValue[string]{
								Name:  "password",
								Usage: "The password to login with",
							},
						},
						Action: func(ctx *replyme.Context) error {
							ctx.Printf("Hello, %s!\n", ctx.GetFlagString("username", "anonymous"))
							return nil
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
