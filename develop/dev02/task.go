package main

import (
	"fmt"
	"os"
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
    string input_s = "a4bc2d5e"

    output_s, err := Unpack(input_s)

    if err != nil {
        // ошибочка
        os.Exit(1)
    }

    fmt.Printf(output_s)
}

Unpack(string input_s) string, error {
    output_s := make([]rune, len(input_s))

    for i, val := []rune(input_s) {
        if val
    }

    return string(output_s), err
}
