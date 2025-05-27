<img src="assets/replyme_dark.png" alt="Логотип Replyme" style="width: 200px;display: block; margin-right: auto; margin-left: auto;"/>
<h1 style="text-align: center">Replyme</h1>
<p style="text-align: center">Инструмент для создания REPL сессий в Golang.</p>
<p style="text-align: center"><b><a href="README.ru.md">[English]</a> [Русский]</b></p>

> [!CAUTION]
> Этот проект может содержать ошибки и недоработки, так как его разработка только начата. Подробнее об этом в разделе "Баги и Roadmap"

### Использование

Для того чтобы использовать Replyme, сначала установите его в ваш проект с помощью команды:

```bash
go get github.com/danyasatsuk/replyme
```

После этого приступайте к созданию своего первого REPL!

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

Найти больше информации можно в директории [`examples`](/examples/README.ru.md)

### Баги и Roadmap

Будьте осторожнее, на данный момент этот проект еще в стадии самой начальной разработки, и в нем может быть очень много недоработок. На данный момент будущие правки внесены в файл [TODO.md (русский)](./TODO.md). Если вы нашли ошибку, то, пожалуйста, опишите ее, создав новый Issue, я буду очень вам благодарен!