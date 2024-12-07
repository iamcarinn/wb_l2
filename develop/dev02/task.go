package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"unicode"
)

/*
=== Задача на распаковку ===

Создать Go функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы / руны, например:
	- "a4bc2d5e" => "aaaabccddddde"
	- "abcd" => "abcd"
	- "45" => "" (некорректная строка)
	- "" => ""
Дополнительное задание: поддержка escape - последовательностей
	- qwe\4\5 => qwe45 (*)
	- qwe\45 => qwe44444 (*)
	- qwe\\5 => qwe\\\\\ (*)

В случае если была передана некорректная строка функция должна возвращать ошибку. Написать unit-тесты.

Функция должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func main() {
	inputS := "a4bc2d5e"

	// Распаковываем строку
	outputS, err := Unpack(inputS)

	if err != nil {
		os.Exit(1) // ошибка распаковки строки
	}

	fmt.Printf(outputS + "\n")
}

// Unpack распаковывает строку с учетом повторений и escape-последовательностей
func Unpack(input string) (string, error) {
	// Если строка пуста, возвращаем пустой результат
	if len(input) == 0 {
		return "", nil
	}

	result := []rune{}
	var prev rune    // предыдущий символ (для повторений)
	var escaped bool // флаг, показывающий, что обр-ся escape-последовательность

	//  Цикл по всем символам строки
	for i, char := range input {
		switch {
		case escaped:
			// если предыдущий символ '\', обрабатываем escape-последовательность
			if err := processEscape(&result, char); err != nil {
				return "", err
			}
			escaped = false
			prev = char
		case char == '\\':
			// если текущий символ — '\', включаем режим обработки escape
			escaped = true
		case unicode.IsDigit(char):
			// если текущий символ — цифра, повторяем предыдущий символ
			if err := func() error {
				var _ int = i
				return repeatChar(&result, prev, char)
			}(); err != nil {
				return "", err
			}
		case unicode.IsLetter(char):
			// если текущий символ — буква, добавляем его в результат
			prev = appendChar(&result, char)
		default:
			// если символ недопустим, возвращаем ошибку
			return "", errors.New("invalid input string: contains non-alphanumeric characters")
		}
	}

	// если строка оканчивается на '\', это ошибка
	if escaped {
		return "", errors.New("trailing escape character")
	}

	return string(result), nil
}

// processEscape обрабатывает символ после escape
func processEscape(result *[]rune, char rune) error {
	// разрешены только цифры и '\'
	if !unicode.IsDigit(char) && char != '\\' {
		return errors.New("invalid escape sequence")
	}
	*result = append(*result, char) // добавляем символ в результат
	return nil
}

// repeatChar повторяет предыдущий символ в зависимости от цифры
func repeatChar(result *[]rune, prev rune, char rune) error {
	// если нет предыдущего символа, это ошибка
	if prev == 0 {
		return errors.New("invalid input string")
	}
	count, _ := strconv.Atoi(string(char)) // конвертируем цифру в число
	// повторяем предыдущий символ нужное количество раз
	for j := 1; j < count; j++ {
		*result = append(*result, prev)
	}
	return nil
}

// appendChar добавляет символ в результат и возвращает его как предыдущий
func appendChar(result *[]rune, char rune) rune {
	*result = append(*result, char)
	return char
}
