package main

import (
	"bytes"
	"io"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"path/filepath"
)

// Тестируем команду cd и pwd
func TestCdAndPwd(t *testing.T) {
	shell := NewShell("")

	// Текущая директория
	initialDir, _ := os.Getwd()

	// Создаем временную директорию
	tempDir := os.TempDir()

	// Переходим в временную директорию
	shell.cmdCD([]string{"cd", tempDir})

	// Получаем текущую директорию
	cwd, _ := os.Getwd()

	// Убираем префикс "/private" для macOS
	cwd = strings.TrimPrefix(cwd, "/private")

	// Нормализуем пути перед сравнением
	if filepath.Clean(cwd) != filepath.Clean(tempDir) {
		t.Errorf("Expected directory: %s, got: %s", tempDir, cwd)
	}

	// Возвращаемся в начальную директорию
	shell.cmdCD([]string{"cd", initialDir})
}


// Тестируем команду echo
func TestEcho(t *testing.T) {
	shell := NewShell("")
	output := captureOutput(func() {
		shell.cmdEcho([]string{"echo", "Hello,", "world!"})
	})

	expected := "Hello, world!\n"
	if output != expected {
		t.Errorf("Expected output: %q, got: %q", expected, output)
	}
}

// Тестируем команду kill
func TestKill(t *testing.T) {
	// Создаем процесс, который будет завершен
	cmd := exec.Command("sleep", "10")
	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start test process: %v", err)
	}

	// Убиваем процесс
	shell := NewShell("")
	shell.cmdKill([]string{"kill", strconv.Itoa(cmd.Process.Pid)})

	// Проверяем, что процесс завершен
	if err := cmd.Wait(); err == nil {
		t.Error("Expected process to be killed")
	}
}

// Тестируем команду ps
func TestPs(t *testing.T) {
	output := captureOutput(func() {
		NewShell("").cmdPS()
	})

	// Проверяем, что в выводе присутствует текущий процесс
	if !strings.Contains(output, "ps") {
		t.Error("Expected 'ps' command in output")
	}
}

// Тестируем команду nc (netcat)
func TestNetcat(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to start test server: %v", err)
	}
	defer listener.Close()

	go func() {
		conn, _ := listener.Accept()
		defer conn.Close()

		// Читаем данные, отправленные клиентом
		buf := make([]byte, 256)
		n, _ := conn.Read(buf)
		if string(buf[:n]) != "Hello, netcat!" {
			t.Errorf("Expected 'Hello, netcat!', got: %s", buf[:n])
		}
	}()

	shell := NewShell("")
	stdinReader, stdinWriter, _ := os.Pipe() // Создаем пайп для stdin
	defer stdinReader.Close()
	defer stdinWriter.Close()

	originalStdin := os.Stdin
	os.Stdin = stdinReader // Перенаправляем stdin
	defer func() { os.Stdin = originalStdin }()

	go func() {
		stdinWriter.WriteString("Hello, netcat!")
		stdinWriter.Close()
	}()

	// Запускаем команду nc
	shell.cmdNetcat([]string{"nc", "127.0.0.1", strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)})
}

// Вспомогательная функция для захвата вывода
func captureOutput(f func()) string {
	stdoutReader, stdoutWriter, _ := os.Pipe() // Создаем пайп для stdout
	defer stdoutReader.Close()

	originalStdout := os.Stdout
	os.Stdout = stdoutWriter // Перенаправляем stdout
	defer func() { os.Stdout = originalStdout }()

	output := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, stdoutReader)
		output <- buf.String()
	}()

	f()
	stdoutWriter.Close()
	return <-output
}
