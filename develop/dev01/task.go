package main

import (
	"fmt"
	"os"
    "strconv"
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
    input_s := "a4bc2d5e"

    output_s, err := Unpack(input_s)

    if err != nil {
        // ошибочка
        os.Exit(1)
    }
 
    fmt.Printf(output_s + "\n")
}

func Unpack(input_s string) (string, error) {
    if len(input_s) == 0 {
		return "", nil
	}

    var output_s []rune
    var rune_add rune
    var count int

    for _, char := range input_s {
        if char >= '0' && char <= '9' {
            count, _ = strconv.Atoi(string(char))
        } else {
            rune_add = char
            output_s = append(output_s, rune_add)
            count = 1
        }
        for i := 0; i < count - 1; i++ {
			output_s = append(output_s, rune_add)
		}
    }

    if len(output_s) == 0 {
		return "", fmt.Errorf("invalid input string")
	}

	return string(output_s), nil
}
