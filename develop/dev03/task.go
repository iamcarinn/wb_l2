package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

/*
=== Утилита sort ===

Отсортировать строки в файле по аналогии с консольной утилитой sort (man sort — смотрим описание и основные параметры): на входе подается файл из несортированными строками, на выходе — файл с отсортированными.

Реализовать поддержку утилитой следующих ключей:

-k — указание колонки для сортировки (слова в строке могут выступать в качестве колонок, по умолчанию разделитель — пробел)
-n — сортировать по числовому значению
-r — сортировать в обратном порядке
-u — не выводить повторяющиеся строки

Дополнительно

Реализовать поддержку утилитой следующих ключей:

-M — сортировать по названию месяца
-b — игнорировать хвостовые пробелы
-c — проверять отсортированы ли данные
-h — сортировать по числовому значению с учетом суффиксов

*/

// Flags - структура, хранящая значения флагов
type Flags struct {
	k int64
	n bool
	r bool
	u bool
}

var fl Flags

func init() {
    flag.Int64Var(&fl.k, "k", 0, "enter column (0 means the whole line)")
    flag.BoolVar(&fl.n, "n", false, "sort by num")
    flag.BoolVar(&fl.r, "r", false, "reverse sort")
    flag.BoolVar(&fl.u, "u", false, "do not duplicate lines")
}

func main() {
	// Используем стандартные файлы input.txt и output.txt
	err := ProcessSort("input.txt", "output.txt")
	if err != nil {
		fmt.Println(err)
	}
}

// ProcessSort выполняет сортировку с учетом флагов.
func ProcessSort(inputFile, outputFile string) error {
	flag.Parse()

	// Чтение строк из файла
	lines, err := readFile(inputFile)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	// Процесс сортировки
	lines = sortProcess(lines)

	// Запись отсортированных строк в файл
	err = writeFile(outputFile, lines)
	if err != nil {
		return fmt.Errorf("error writing file: %v", err)
	}

	fmt.Println("Sorting completed. Check", outputFile)
	return nil
}

// Функция сортировки с обработкой флагов
func sortProcess(lines []string) []string{
	// Флаг -u — не выводить повторяющиеся строки
	if fl.u {
		lines = removeDuplicates(lines)
	}

	// Флаг -k — указание колонки для сортировки
	if fl.k > 0 {
		sortByColumn(lines, fl.k, fl.n)

	// Флаг -n — сортировать по числовому значению
	} else if fl.n {
		sortByNumber(lines, fl.n)
	} else {
		// сортировка по умолчанию
		sort.Strings(lines)
	}

	// Флаг -r — сортировать в обратном порядке
	if fl.r {
		for i, j := 0, len(lines)-1; i < j; i, j = i+1, j-1 {
			lines[i], lines[j] = lines[j], lines[i]
		}
	}

	return lines
}

// Функция чтения файла
func readFile(filepath string) ([]string, error) {
	file, err := os.Open(filepath)	// Открываем почитать
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)	// будет считывать файл строка за строкой
	for scanner.Scan() {	// возвращает true пока удается считать строку
		lines = append(lines, scanner.Text())	// scanner.Text() возвращвет текущую строку
	}
	return lines, scanner.Err()
}

// Функция создание файлв им записи в него
func writeFile(filepath string, lines []string) error {
	file, err := os.Create(filepath)	// Создаем/пересоздаем файл
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)	// будет записывать данные в файл с буферизацией
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")	// записываем строку в буфер
		if err != nil {
			return err
		}
	}
	return writer.Flush()	// записываем содержимое буфера в файл, возвращает ошибку
}

// Функция обработки флага -u - удаление дубликатов
func removeDuplicates(lines []string) []string {
	var result []string

	seen := make(map[string]bool) // хранятся все стреченные строки

	for _, line := range lines {
		if !seen[line] {	// если строки еще нет в seem, то добавляем ее туда
			seen[line] = true
			result = append(result, line)	// добавляем строку  в результат
		}
	}

	return result
}

// Функция сортировки строки по указанной колонке.
func sortByColumn(lines []string, column int64, numeric bool) {
	sort.Slice(lines, func(i, j int) bool {
		// Разделяем строки на слова
		columnsI := strings.Fields(lines[i])
		columnsJ := strings.Fields(lines[j])

		// проверяем наличие колонки
		valI := ""
		if int(column) <= len(columnsI) {
			valI = columnsI[column-1]
		}
		valJ := ""
		if int(column) <= len(columnsJ) {
			valJ = columnsJ[column-1]
		}

		if numeric {
			// Флаг -n — сортировать по числовому значению
			numI, _ := strconv.Atoi(valI)
			numJ, _ := strconv.Atoi(valJ)
			return numI < numJ
		}

		// Лексикографическая сортировка
		return valI < valJ
	})
}

// Функция числовой сортировки
func sortByNumber(lines []string, numeric bool) {
	sort.Slice(lines, func(i, j int) bool {
		// конвертируем строки в числа
		numI, errI := strconv.Atoi(lines[i])
		numJ, errJ := strconv.Atoi(lines[j])

		// при ошибки конвертации, сравниваем строки лексикографически.
		if errI != nil || errJ != nil {
			return lines[i] < lines[j]
		}

		// числовая сортировка
		return numI < numJ
	})
}
