package main

/*
=== Утилита telnet ===

Реализовать примитивный telnet клиент:
Примеры вызовов:
go-telnet --timeout=10s host port go-telnet mysite.ru 8080 go-telnet --timeout=3s 1.1.1.1 123

Программа должна подключаться к указанному хосту (ip или доменное имя) и порту по протоколу TCP.
После подключения STDIN программы должен записываться в сокет, а данные полученные и сокета должны выводиться в STDOUT
Опционально в программу можно передать таймаут на подключение к серверу (через аргумент --timeout, по умолчанию 10s).

При нажатии Ctrl+D программа должна закрывать сокет и завершаться. Если сокет закрывается со стороны сервера, программа должна также завершаться.
При подключении к несуществующему сервер, программа должна завершаться через timeout.
*/
import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"syscall"
	"time"
)

var timeout string

func init() {
	flag.StringVar(&timeout, "timeout", "10s", "connection time limits")
}

// Функция валидирует таймаут
func validateTimeout(timeout string) time.Duration {
	re := regexp.MustCompile(`^\d+s$`)
	if !re.MatchString(timeout) {
		log.Panicf("Invalid timeout format: %s. Expected format is a number followed by 's' (e.g., '10s')", timeout)
	}

	timeToInt, err := strconv.Atoi(timeout[:len(timeout)-1]) // приводим к int
	if err != nil {
		log.Panicf("Failed to parse timeout: %v", err)
	}
	return time.Duration(timeToInt) * time.Second // приводим к секундам
}

// Функция устанавливает соединение с сервером с учетом таймаута
func connectToServer(host, port string, timeout time.Duration) net.Conn {
	address := fmt.Sprintf("%s:%s", host, port)
	// уст-ем соединение с исп-ем таймаута
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		log.Fatalf("Failed to connect to %s: %v", address, err)
	}
	log.Printf("Connected to %s", address)
	return conn
}

// Функция обрабатывает ввод и вывод между сокетом и stdin/stdout
func handleIO(conn net.Conn) {
	// передача данных из stdin в сокет
	go func() {
		_, err := io.Copy(conn, os.Stdin)
		if err != nil {
			log.Printf("Error writing to socket: %v", err)
		}
	}()

	// передача данных из сокета в stdout
	go func() {
		_, err := io.Copy(os.Stdout, conn)
		if err != nil {
			log.Printf("Error reading from socket: %v", err)
		}
	}()
}

// Функция ожидает завершения программы через сигнал или таймаут
func waitForExit(timeout time.Duration, conn net.Conn) {
	signalChan := make(chan os.Signal, 1)	// канал для обработки системных сигналов
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM) // подписка канала на полуения сигналов

	select {
	case <-signalChan:			// отправлен сигнал о прерывании (например Ctrl + C)
		log.Println("Interrupt signal received, closing connection...")
	case <-time.After(timeout):	// время истекло
		log.Println("Timeout exceeded, closing connection...")
	}

	conn.Close()
	log.Println("Connection closed.")
}

func runApp(host, port string, timeout time.Duration) {
	// подключение к серверу
	conn := connectToServer(host, port, timeout)
	defer conn.Close()

	// обработка ввода/вывода
	handleIO(conn)

	// ожидание завершения
	waitForExit(timeout, conn)
}

func main() {
	flag.Parse()
	
	if len(flag.Args()) < 2 {
		log.Fatal("usage: --timeout=1s host port")
	}

	host := flag.Arg(0)
	port := flag.Arg(1)
	
	timeoutDuration := validateTimeout(timeout)

	runApp(host, port, timeoutDuration)
}


