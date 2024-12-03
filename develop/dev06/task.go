package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

/*
=== Утилита cut ===

Принимает STDIN, разбивает по разделителю (TAB) на колонки, выводит запрошенные

Поддержать флаги:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type cutOptions struct {
	fields    string // -f
	delimiter string // -d
	separated bool   // -s
}

var options cutOptions

func init() {
	flag.StringVar(&options.fields, "f", "", "Select fields (e.g., 1,2-3)")
	flag.StringVar(&options.delimiter, "d", "\t", "Column delimiter (default is TAB)")
	flag.BoolVar(&options.separated, "s", false, "Only lines with the delimiter")
}

func main() {
	flag.Parse()

	// Флаг -f
	if options.fields == "" {
		fmt.Fprintln(os.Stderr, "Error: -f flag is required")
		os.Exit(1)
	}

	// Читаем строки из STDIN
	lines := readInput()

	// Обрабатываем каждую строку
	for _, line := range lines {
		processLine(line)
	}
}

// Функция чтения строк из STDIN
func readInput() []string {
	var lines []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	return lines
}

// Функция обработки строки
func processLine(line string) {
	// Флаг -s (пропускаем строки без разделителя)
	if options.separated && !strings.Contains(line, options.delimiter) {
		return
	}

	columns := strings.Split(line, options.delimiter)	// разбиваем строку на колонки
	selected := selectFields(columns, options.fields)	// получаем выбранные колонки

	// Выводим результат
	if len(selected) > 0 {
		fmt.Println(strings.Join(selected, options.delimiter))
	}
}

// Функция выбора указанных колонок
func selectFields(columns []string, fields string) []string {
	var result []string
	ranges := strings.Split(fields, ",")

	for _, r := range ranges {
		if strings.Contains(r, "-") {
			// Обрабатываем диапазон
			parts := strings.Split(r, "-")
			start, _ := strconv.Atoi(parts[0])
			end, _ := strconv.Atoi(parts[1])
			for i := start; i <= end && i <= len(columns); i++ {
				if i > 0 {
					result = append(result, columns[i-1])
				}
			}
		} else {
			// Обрабатываем одиночный номер
			index, _ := strconv.Atoi(r)
			if index > 0 && index <= len(columns) {
				result = append(result, columns[index-1])
			}
		}
	}

	return result
}
