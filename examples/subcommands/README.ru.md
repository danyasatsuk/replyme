# `subcommands` - создание подкоманд

<img src="subcommands.png" alt="Картинка с флагами" style="max-width: 600px;"/>

Для того чтобы создать подкоманду для вашей основной команды, используйте поле `Subcommands` внутри `[]*replyme.Command`:

```go
Subcommands: []&replyme.Command{
    {
        Name: "mySubCommand",
        Usage: "testMySubCommand",
        Action: func (ctx *replyme.Context) error {
            // ваш код
			return nil
        }
    }
}
```

Таким образом вы можете добавлять дополнительные подкоманды в ваш REPL. 

### Управление флагами

Флаги можно использовать как и у основных команд, так и у дополнительных. 

Например: если у основной команды `create` были описаны флаги, и после вы добавили подкоманду `data`, то тогда такая команда будет полноценно работать:

```plain text
create --myCreateFlag="test" data --myDataFlag="test"
```

Чтобы ознакомится с примером работы, посмотрите файл [subcommands/main.go](./main.go).

### Поток действий

Для того чтобы работать с командами и подкомандами, есть определенный поток выполнения, о нем можно посмотреть внутри примера [Flow](../flow/README.ru.md)