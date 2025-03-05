Что выведет программа? Объяснить вывод программы.

```go
package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	{
		// do something
	}
	return nil
}

func main() {
	var err error
	err = test()
	if err != nil {
		println("error")
		return
	}
	println("ok")
}
```

Ответ:
```
error
```
**Почему же `err != nil`, хотя в функции `test()` происходит `return nil`?**
- Функция `test()` возвращает `*customError`, а когда мы делаем `return nil`, Go выполняет автоматическое преобразование к интерфейсу `error`.

- В результате в `err` попадает интерфейс со:
	- **Статическм типом** `error`
	- **Динамическим типом** `*customError`
	- **Динамическим значением** `nil` (сам указатель `nil`)

Из-за наличия динамического типа интерфейс не `nil`, поэтому `err != nil`.

> **Интерфейс считается `nil`**, если его динамический тип и динамическое значение оба равны `nil`.
