package main

import (
	"os"
	"testing"
)

func TestProcessSort(t *testing.T) {
	// Задаем входной и выходной файл
	inputFile := "input.txt"
	outputFile := "output.txt"

	// Сортируем
	err := ProcessSort(inputFile, outputFile)
	if err != nil {
		t.Fatalf("Error in ProcessSort: %v", err)
	}

	// Чтение содержимого output.txt
	expectedOutput := "10 fig\n100 apple\n100 apple\n25 cherry\n50 banana\n" // Здесь ты можешь указать ожидаемый результат
	actualOutput, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Error reading output file: %v", err)
	}

	// Проверка, что содержимое output.txt соответствует ожидаемому
	if string(actualOutput) != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, string(actualOutput))
	}
}
