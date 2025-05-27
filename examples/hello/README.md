# `hello` - welcome message

You can use this code to create a welcome message, and it's easy to do so.:

## 1. Create a variable with the type `*replyme.App`:

```go
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
```

## 2. Run the 'replyme`:

```go
err := replyme.Run(app)
if err != nil {
    panic(err)
}
```

Ready! Now you can run your code using `go run`