<img src="assets/replyme_dark.png" alt="Replyme Logo" style="width: 200px;display: block; margin-right: auto; margin-left: auto;"/>
<h1 style="text-align: center">Replyme</h1>
<p style="text-align: center">A tool for creating REPL sessions in Golang.</p>
<p style="text-align: center"><b>[English] <a href="README.ru.md">[Русский]</a></b></p>

> [!CAUTION]
> This project may contain errors and flaws, as its development has just begun. Read more about this in the Bugs and Roadmap section

### Usage

To use Replyme, first install it in your project using the command:

```bash
go get github.com/danyasatsuk/replyme
```

After that, start creating your first REPL!

```go
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

```

You can find out more about how it works in the [`examples`](/examples/README.md) directory.

### Bugs and Roadmap

Be careful, at the moment this project is still in the very initial development stage, and there may be a lot of flaws in it. At the moment, future edits have been made to the file. [TODO.md (Russian)](./TODO.md). If you find a mistake, please describe it by creating a new Issue, I will be very grateful to you!

