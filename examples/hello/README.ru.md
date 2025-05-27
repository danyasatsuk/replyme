# `hello` - приветственное сообщение

С помощью этого кода можно создать приветственное сообщение, делается это легко:

## 1. Создайте переменную с типом `*replyme.App`:

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

## 2. Запустите `replyme`:

```go
err := replyme.Run(app)
if err != nil {
    panic(err)
}
```

Готово! Теперь можно запускать ваш код с помощью `go run`