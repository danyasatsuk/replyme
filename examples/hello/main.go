package main

import "github.com/danyasatsuk/replyme"

func main() {
	app := &replyme.App{
		Name:  "hello",
		Usage: "Hello World",
		Commands: []*replyme.Command{
			{
				Name:  "hello",
				Usage: "Print Hello World",
				Action: func(ctx *replyme.Context) error {
					ctx.Print("Hello, World!")
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
