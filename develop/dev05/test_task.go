package main

import (
	//"os"
	"strings"
	"testing"
)

// Тест для флага -A (печатает строки после совпадения)
func TestGrepWithAfterFlag(t *testing.T) {
	// Подготовим флаги
	options := grepOptions{
		after: 2,    // Печатаем 2 строки после совпадения
		pattern: "example",
	}

	// Читаем данные из файла input.txt
	lines, err := readInput()
	if err != nil {
		t.Fatalf("Error reading input: %v", err)
	}

	// Применяем grep
	matches, totalMatches := grep(lines, options)

	println(totalMatches)
	// Ожидаем, что будет найдено 1 совпадение
	if totalMatches != 1 {
		t.Errorf("Expected 1 matches, got %d", totalMatches)
	}

	// Проверяем, что вывод содержит соответствующие строки
	expectedMatches := []string{
		"For example, some lines may contain the word \"test\".",
		"Others may not contain the word at all.",
		"Case sensitivity is also something we want to test.",
	}

	for i, match := range matches {
		if !strings.Contains(match, expectedMatches[i]) {
			t.Errorf("Expected match: %v, got: %v", expectedMatches[i], match)
		}
	}
}

// // Тест для флага -B (печатает строки до совпадения)
// func TestGrepWithBeforeFlag(t *testing.T) {
// 	options := grepOptions{
// 		before: 2,  // Печатаем 2 строки до совпадения
// 		pattern: "test",
// 	}

// 	lines, err := readInput()
// 	if err != nil {
// 		t.Fatalf("Error reading input: %v", err)
// 	}

// 	matches, totalMatches := grep(lines, options)

// 	// Ожидаем 4 совпадения
// 	if totalMatches != 4 {
// 		t.Errorf("Expected 4 matches, got %d", totalMatches)
// 	}

// 	// Проверка совпадений
// 	expectedMatches := []string{
// 		"This is a simple test file.",
// 		"It contains several lines of text.",
// 		"Each line is designed to help test various cases.",
// 		"For example, some lines may contain the word \"test\".",
// 		"TEST should match test if the -i flag is used.",
// 	}

// 	for i, match := range matches {
// 		if !strings.Contains(match, expectedMatches[i]) {
// 			t.Errorf("Expected match: %v, got: %v", expectedMatches[i], match)
// 		}
// 	}
// }

// // Тест для флага -C (печатает строки вокруг совпадения)
// func TestGrepWithContextFlag(t *testing.T) {
// 	options := grepOptions{
// 		context: 2, // Печатаем 2 строки до и 2 строки после совпадения
// 		pattern: "test",
// 	}

// 	lines, err := readInput()
// 	if err != nil {
// 		t.Fatalf("Error reading input: %v", err)
// 	}

// 	matches, totalMatches := grep(lines, options)

// 	// Ожидаем 4 совпадения
// 	if totalMatches != 4 {
// 		t.Errorf("Expected 4 matches, got %d", totalMatches)
// 	}

// 	// Проверка совпадений
// 	expectedMatches := []string{
// 		"This is a simple test file.",
// 		"It contains several lines of text.",
// 		"Each line is designed to help test various cases.",
// 		"For example, some lines may contain the word \"test\".",
// 		"TEST should match test if the -i flag is used.",
// 	}

// 	for i, match := range matches {
// 		if !strings.Contains(match, expectedMatches[i]) {
// 			t.Errorf("Expected match: %v, got: %v", expectedMatches[i], match)
// 		}
// 	}
// }

// // Тест для флага -i (игнорирование регистра)
// func TestGrepWithIgnoreCaseFlag(t *testing.T) {
// 	options := grepOptions{
// 		ignoreCase: true, // Игнорируем регистр
// 		pattern:    "test",
// 	}

// 	lines, err := readInput()
// 	if err != nil {
// 		t.Fatalf("Error reading input: %v", err)
// 	}

// 	matches, totalMatches := grep(lines, options)

// 	// Ожидаем 4 совпадения
// 	if totalMatches != 4 {
// 		t.Errorf("Expected 4 matches, got %d", totalMatches)
// 	}

// 	// Проверка совпадений
// 	expectedMatches := []string{
// 		"This is a simple test file.",
// 		"TEST should match test if the -i flag is used.",
// 	}

// 	for i, match := range matches {
// 		if !strings.Contains(match, expectedMatches[i]) {
// 			t.Errorf("Expected match: %v, got: %v", expectedMatches[i], match)
// 		}
// 	}
// }

// // Тест для флага -c (выводим только количество совпадений)
// func TestGrepWithCountFlag(t *testing.T) {
// 	options := grepOptions{
// 		count:   true, // Выводим количество совпадений
// 		pattern: "test",
// 	}

// 	lines, err := readInput()
// 	if err != nil {
// 		t.Fatalf("Error reading input: %v", err)
// 	}

// 	matches, totalMatches := grep(lines, options)

// 	// Ожидаем 4 совпадения
// 	if totalMatches != 4 {
// 		t.Errorf("Expected 4 matches, got %d", totalMatches)
// 	}

// 	// Ожидаем, что выводится только количество совпадений
// 	if len(matches) != 0 {
// 		t.Errorf("Expected no matches in output, got %v", matches)
// 	}
// }

// // Тест для флага -v (инвертируем результат)
// func TestGrepWithInvertFlag(t *testing.T) {
// 	options := grepOptions{
// 		invert:  true, // Инвертируем совпадение
// 		pattern: "test",
// 	}

// 	lines, err := readInput()
// 	if err != nil {
// 		t.Fatalf("Error reading input: %v", err)
// 	}

// 	matches, totalMatches := grep(lines, options)

// 	// Ожидаем 2 совпадения (все строки, кроме содержащих "test")
// 	if totalMatches != 2 {
// 		t.Errorf("Expected 2 matches, got %d", totalMatches)
// 	}

// 	// Проверка совпадений
// 	expectedMatches := []string{
// 		"It contains several lines of text.",
// 		"Others may not contain the word at all.",
// 	}

// 	for i, match := range matches {
// 		if !strings.Contains(match, expectedMatches[i]) {
// 			t.Errorf("Expected match: %v, got: %v", expectedMatches[i], match)
// 		}
// 	}
// }
