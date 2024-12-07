package main

/*
=== Взаимодействие с ОС ===

Необходимо реализовать собственный шелл

встроенные команды: cd/pwd/echo/kill/ps
поддержать fork/exec команды
конвеер на пайпах

Реализовать утилиту netcat (nc) клиент
принимать данные из stdin и отправлять в соединение (tcp/udp)
Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Shell - структура с приглашением к вводу
type Shell struct {
	prompt string
}

// NewShell создает новую структуру Shell
func NewShell(prompt string) *Shell {
	return &Shell{prompt: prompt}
}

// Run циклически обрабатывает команды
func (s *Shell) Run() {
	reader := bufio.NewReader(os.Stdin) // для чтения из stdin

	for {
		fmt.Print(s.prompt)	// приглашаем к вводу
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input) //удаляем пробелы и \n

		if input == "" {
			continue
		}

		args := strings.Fields(input) // строка разбивается на аргументы

		//  обработка встроенных команд
		if s.executeBuiltin(args) {
			continue
		}

		// обработка внешних команд
		s.runExternalCommand(args)
	}
}

// Функция обработки встроенных команды
func (s *Shell) executeBuiltin(args []string) bool {
	switch args[0] {
	case "cd": 			// смена директории
		s.cmdCD(args)
	case "pwd": 		// текущая директория
		s.cmdPWD()
	case "echo": 		// вывод текста
		s.cmdEcho(args)
	case "kill": 		// убить процесс по PID
		s.cmdKill(args)
	case "ps": 			// список процессов
		s.cmdPS()
	case "nc": 			// netcat (TCP/UDP клиент)
		s.cmdNetcat(args)
	default:
		return false 	// если команда не встроенная
	}
	return true
}

// Функция обрабатывает команду cd
func (s *Shell) cmdCD(args []string) {
	if len(args) < 2 { // проверка, указан ли путь
		fmt.Println("Usage: cd <directory>")
		return
	}
	if err := os.Chdir(args[1]); err != nil { // меняем текущую директорию
		fmt.Println("Error:", err)
	}
}

// Функция обработки команды pwd
func (s *Shell) cmdPWD() {
	dir, err := os.Getwd() // получаем текущую директорию
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(dir) // выводим путь до текущей директории
}

// Функция обработки команды echo
func (s *Shell) cmdEcho(args []string) {
	fmt.Println(strings.Join(args[1:], " ")) // выводим аргументы без "echo"
}

// Функция обработки команды kill
func (s *Shell) cmdKill(args []string) {
	// проверка, указан ли PID
	if len(args) < 2 {
		fmt.Println("Usage: kill <pid>")
		return
	}

	pid := parsePID(args[1]) // преобразуем PID в целое число

	// получаем процесс по PID
	process, err := os.FindProcess(pid)
	if err != nil {
		fmt.Println("Error finding process:", err)
		return
	}

	// убиваем процесс
	if err := process.Kill(); err != nil {
		fmt.Println("Error killing process:", err)
	}

}

// Функция обработки команды ps
func (s *Shell) cmdPS() {
	cmd := exec.Command("ps") // запускаем системную команду ps
	cmd.Stdout = os.Stdout    // Перенаправляем вывод в консоль
	cmd.Stderr = os.Stderr    // Перенаправляем ошибки в консоль
	if err := cmd.Run(); err != nil {
		fmt.Println("Error:", err)
	}
}

// Функция обработки команды nc (подключениt к указанному хосту через TCP)
func (s *Shell) cmdNetcat(args []string) {
	if len(args) < 3 { // проверяем, указаны ли хост и порт
		fmt.Println("Usage: nc <host> <port>")
		return
	}
	host, port := args[1], args[2]
	conn, err := net.Dial("tcp", net.JoinHostPort(host, port)) // устанавливаем TCP-соединение
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer conn.Close() // закрываем соединение

	// поток для чтения данных с сервера
	go func() {
		_, _ = io.Copy(os.Stdout, conn)
	}()
	// отправляем данные из stdin на сервер
	_, _ = io.Copy(conn, os.Stdin)
}

// Функция запускает внешние команды (ls, cat и др).
func (s *Shell) runExternalCommand(args []string) {
	cmd := exec.Command(args[0], args[1:]...) // создаем команду с аргументами
	cmd.Stdout = os.Stdout                   // перенаправляем вывод в консоль
	cmd.Stderr = os.Stderr                   // перенаправляем ошибки в консоль
	if err := cmd.Run(); err != nil {        // выполняем команду
		fmt.Println("Error:", err)
	}
}

// Функция преобразует строку в PID (целое число).
func parsePID(pid string) int {
	p, _ := strconv.Atoi(pid)
	return p
}

func main() {
	shell := NewShell("shell> ") // структура с приглашением к вводу
	shell.Run()                 // цикл обработки команд
}
