package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

/*
=== Утилита grep ===

Реализовать утилиту фильтрации (man grep)

Поддержать флаги:
-A - "after" печатать +N строк после совпадения
-B - "before" печатать +N строк до совпадения
-C - "context" (A+B) печатать ±N строк вокруг совпадения
-c - "count" (количество строк)
-i - "ignore-case" (игнорировать регистр)
-v - "invert" (вместо совпадения, исключать)
-F - "fixed", точное совпадение со строкой, не паттерн
-n - "line num", печатать номер строки

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type grepOptions struct {
	after		int	// -A
	before		int	// -B
	context		int	// -C
	count		bool	// -c
	ignoreCase	bool 	// -i
	invert 		bool	// -v
	fixed		bool	// -F
	lineNum		bool	// -n
	pattern		string	// паттерн
}

var options grepOptions

func init() {
	flag.IntVar(&options.after, "A", 0, "Print +N lines after match")
	flag.IntVar(&options.before, "B", 0, "Print +N lines before match")
	flag.IntVar(&options.context, "C", 0, "Print ±N lines around match")
	flag.BoolVar(&options.count, "c", false, "Print count of matching lines")
	flag.BoolVar(&options.ignoreCase, "i", false, "Ignore case")
	flag.BoolVar(&options.invert, "v", false, "Invert match")
	flag.BoolVar(&options.fixed, "F", false, "Exact match, not a pattern")
	flag.BoolVar(&options.lineNum, "n", false, "Print line numbers")
}

func parseFlags() {
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "Error: pattern not provided")
		os.Exit(1)
	}
	options.pattern = flag.Arg(0)
}

// Функция открытия файла и чтения из него или из stdin
func readInput() ([]string, error) {
	var lines []string

	// Если указаны файлы
	if flag.NArg() > 1 {
		// открываем файл
		for _, fileName := range flag.Args()[1:] {
			file, err := os.Open(fileName)
			if err != nil {
				return nil, fmt.Errorf("could not open file %s: %w", fileName, err)
			}
			defer file.Close()
			
			// читаем из файла
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				lines = append(lines, scanner.Text())
			}
			if err := scanner.Err(); err != nil {
				return nil, err
			}
		}
		return lines, nil
	}

	// Если файл не указан, читаем из stdin
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

// Функция применяет фильтры к строкам
func grep(lines []string, options grepOptions) ([]string, int) {
	var matches []string	// строки для вывода
	totalMatches := 0	// число совпадений
	// Мама для отслеживания добавленных строк
	seen := map[int]struct{}{}

	// Нормализуем паттерн, если включен -i
	pattern := options.pattern
	if options.ignoreCase {
		pattern = strings.ToLower(pattern)	// приводим к нижнему регистру
	}

	// Проходим по каждой строке
	for i, line := range lines {
		// Флаг -i 
		if options.ignoreCase {
			line = strings.ToLower(line)	// приводим строку к нижнему регистру
		}

		var match bool	// было ли совпадение
		// Флаг -F
		if options.fixed {
			match = line == pattern	// если строка полностью совпадает с паттерном
		} else {
			match = strings.Contains(line, pattern)	// содержится ли паттерн в строке
		}

		// Флаг -v
		if options.invert {
			match = !match
		}

		// Если было совпадение
		if match {
			totalMatches++

			// Определяем диапазон строк, которые нужно добавить в результат
			start := max(0, i-options.before-options.context)
			end := min(len(lines)-1, i+options.after+options.context)

			// Добавляем строки перед, после или с контекстом
			for j := start; j <= end; j++ {
				if _, exists := seen[j]; !exists { // чтобы не добавлять одинаковые
					// Флаг -n
					if options.lineNum {
						matches = append(matches, fmt.Sprintf("%d:%s", j+1, lines[j]))
					} else {
						matches = append(matches, lines[j])
					}
					// Добавляем строку в множество, чтобы избежать дублирования
					seen[j] = struct{}{}
				}
			}
		}
	}
	// Возвращаем найденные строки и общее количество совпадений
	return matches, totalMatches
}

// Функция возвращает максимальное значение
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Функция возвращает минимальное значение
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}


func main() {
	parseFlags()

	// Читаем входные данные
	lines, err := readInput()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input:", err)
		os.Exit(1)
	}

	// Обрабатываем флаги
	matches, totalMatches := grep(lines, options)

	// Выводим результаты
	if options.count {
		fmt.Println(totalMatches)
	} else {
		for _, match := range matches {
			fmt.Println(match)
		}
	}
}
