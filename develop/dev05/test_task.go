package main

import (
	"strings"
	"testing"
)

// Тест для флага -A (печатает строки после совпадения)
func TestGrepWithAfterFlag(t *testing.T) {
    options := grepOptions{
        after:   2,
        pattern: "example",
    }

    lines := []string{
        "This is a test file.",
        "For example, some lines may contain the word \"test\".",
        "Others may not contain the word at all.",
        "Case sensitivity is also something we want to test.",
        "End of the test file.",
    }

    matches, totalMatches := grep(lines, options)

    if totalMatches != 1 {
        t.Errorf("Expected 1 match, got %d", totalMatches)
    }

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

// Тест для флага -B (печатает строки до совпадения)
func TestGrepWithBeforeFlag(t *testing.T) {
    options := grepOptions{
        before:  2,
        pattern: "example",
    }

    lines := []string{
        "This is a test file.",
        "For example, some lines may contain the word \"test\".",
        "Others may not contain the word at all.",
        "Case sensitivity is also something we want to test.",
        "End of the test file.",
    }

    matches, totalMatches := grep(lines, options)

    if totalMatches != 1 {
        t.Errorf("Expected 1 match, got %d", totalMatches)
    }

    expectedMatches := []string{
        "This is a test file.",
        "For example, some lines may contain the word \"test\".",
    }

    for i, match := range matches {
        if !strings.Contains(match, expectedMatches[i]) {
            t.Errorf("Expected match: %v, got: %v", expectedMatches[i], match)
        }
    }
}

// Тест для флага -C (печатает строки вокруг совпадения)
func TestGrepWithContextFlag(t *testing.T) {
    options := grepOptions{
        context: 2,
        pattern: "example",
    }

    lines := []string{
        "This is a test file.",
        "For example, some lines may contain the word \"test\".",
        "Others may not contain the word at all.",
        "Case sensitivity is also something we want to test.",
        "End of the test file.",
    }

    matches, totalMatches := grep(lines, options)

    if totalMatches != 1 {
        t.Errorf("Expected 1 match, got %d", totalMatches)
    }

    expectedMatches := []string{
        "This is a test file.",
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

// Тест для флага -i (игнорирование регистра)
func TestGrepWithIgnoreCaseFlag(t *testing.T) {
    options := grepOptions{
        ignoreCase: true,
        pattern:    "test",
    }

    lines := []string{
        "This is a test file.",
        "For example, some lines may contain the word \"Test\".",
        "TEST should match test if the -i flag is used.",
        "End of the test file.",
    }

    matches, totalMatches := grep(lines, options)

    if totalMatches != 4 {
        t.Errorf("Expected 4 matches, got %d", totalMatches)
    }

    expectedMatches := []string{
        "This is a test file.",
        "For example, some lines may contain the word \"Test\".",
        "TEST should match test if the -i flag is used.",
		"End of the test file.",
    }

    for i, match := range matches {
        if !strings.Contains(match, expectedMatches[i]) {
            t.Errorf("Expected match: %v, got: %v", expectedMatches[i], match)
        }
    }
}

// Тест для флага -c (выводим только количество совпадений)
func TestGrepWithCountFlag(t *testing.T) {
    options := grepOptions{
        count:   true,
        pattern: "test",
    }

    lines := []string{
        "This is a test file.",
        "For example, some lines may contain the word \"test\".",
        "Others may not contain the word at all.",
        "Case sensitivity is also something we want to test.",
    }

    _, totalMatches := grep(lines, options)

    if totalMatches != 3 {
        t.Errorf("Expected 3 matches, got %d", totalMatches)
    }
}

// Тест для флага -v (инвертируем результат)
func TestGrepWithInvertFlag(t *testing.T) {
    options := grepOptions{
        invert:  true,
        pattern: "test",
    }

    lines := []string{
        "This is a test file.",
        "For example, some lines may contain the word \"test\".",
        "Others may not contain the word at all.",
        "End of the test file.",
    }

    matches, totalMatches := grep(lines, options)

    if totalMatches != 1 {
        t.Errorf("Expected 1 matches, got %d", totalMatches)
    }

    expectedMatches := []string{
        "Others may not contain the word at all.",
    }

    for i, match := range matches {
        if !strings.Contains(match, expectedMatches[i]) {
            t.Errorf("Expected match: %v, got: %v", expectedMatches[i], match)
        }
    }
}
