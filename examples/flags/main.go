package main

import "github.com/danyasatsuk/replyme"

func main() {
	app := &replyme.App{
		Name:  "flags",
		Usage: "Ouch, flags!",
		Commands: []*replyme.Command{
			{
				Name: "yourName",
				Flags: []replyme.Flag{
					&replyme.FlagValue[string]{
						Name:  "name",
						Usage: "Your name",
					},
				},
				Action: func(ctx *replyme.Context) error {
					ctx.Printf("Hello, %s!\n", ctx.GetFlagString("name", "Unknown"))
					return nil
				},
			},
			{
				Name: "register",
				Flags: []replyme.Flag{
					&replyme.FlagValue[string]{
						Name:  "login",
						Usage: "Your login",
					},
					&replyme.FlagValue[string]{
						Name:  "password",
						Usage: "Your password",
					},
					&replyme.FlagValue[int]{
						Name:  "count",
						Usage: "How many apples do you want?",
					},
				},
				Action: func(ctx *replyme.Context) error {
					ctx.Printf("Hello, %s!", ctx.GetFlagString("login", "UnknownLogin"))
					ctx.Printf("You have %d apples!", ctx.GetFlagInt("count", 0))
					ctx.Printf("Your password is %s", ctx.GetFlagString("password", "'stop, i don't want to see it'"))
					return nil
				},
			},
		},
	}

	err := replyme.Run(app)
	if err != nil {
		panic(err)
	}
}
