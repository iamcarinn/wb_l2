## Запуск сервера
```bash
go run server/server.go
## Компиляция программы
```bash
go build -o go-telnet task.go
## Запуск программы
```bash
./go-telnet --timeout=5s localhost 8080

## Запуск тестов
```bash
go test -v
