package main

import (
	"bytes"
	"io"
	"net"
	"os"
	"strings"
	"testing"
	"time"
	"fmt"
)
func TestValidateTimeout(t *testing.T) {
	t.Run("Valid timeout", func(t *testing.T) {
		got := validateTimeout("5s")
		want := 5 * time.Second
		if got != want {
			t.Errorf("Expected %v, got %v", want, got)
		}
	})

	t.Run("Invalid timeout", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected a panic for invalid timeout")
			} else {
				if !strings.Contains(fmt.Sprint(r), "Invalid timeout format") {
					t.Errorf("Unexpected panic message: %v", r)
				}
			}
		}()
		validateTimeout("invalid")
	})
}


func TestRunApp(t *testing.T) {
	// Создаем mock-сервер
	listener, err := net.Listen("tcp", "localhost:8082")
	if err != nil {
		t.Fatalf("Failed to start mock server: %v", err)
	}
	defer listener.Close()

	// Обработка подключений на сервере
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			t.Logf("Connection error: %v", err)
			return
		}
		defer conn.Close()
		io.Copy(conn, conn) // Эхо-сервер
	}()

	// Создаем pipe для подмены os.Stdin и os.Stdout
	oldStdin := os.Stdin
	oldStdout := os.Stdout
	defer func() {
		os.Stdin = oldStdin
		os.Stdout = oldStdout
	}()

	stdinReader, stdinWriter, _ := os.Pipe()
	stdoutReader, stdoutWriter, _ := os.Pipe()

	// Подменяем stdin и stdout
	os.Stdin = stdinReader
	os.Stdout = stdoutWriter

	// Ввод данных через stdin
	input := "Hello, runApp!"
	go func() {
		_, _ = stdinWriter.Write([]byte(input))
		stdinWriter.Close()
	}()

	// Чтение вывода из stdout
	var output bytes.Buffer
	go func() {
		io.Copy(&output, stdoutReader)
	}()

	// Вызываем runApp с корректными параметрами
	runApp("localhost", "8082", 5*time.Second)

	// Даем время на выполнение горутин
	time.Sleep(100 * time.Millisecond)

	// Проверяем результат
	if !strings.Contains(output.String(), input) {
		t.Errorf("Expected %s in output, got %s", input, output.String())
	}
}
