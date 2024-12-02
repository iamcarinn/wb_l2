package main

import (
	"sort"
	"strings"
)

/*
=== Поиск анаграмм по словарю ===

Напишите функцию поиска всех множеств анаграмм по словарю.
Например:
'пятак', 'пятка' и 'тяпка' - принадлежат одному множеству,
'листок', 'слиток' и 'столик' - другому.

Входные данные для функции: ссылка на массив - каждый элемент которого - слово на русском языке в кодировке utf8.
Выходные данные: Ссылка на мапу множеств анаграмм.
Ключ - первое встретившееся в словаре слово из множества
Значение - ссылка на массив, каждый элемент которого, слово из множества. Массив должен быть отсортирован по возрастанию.
Множества из одного элемента не должны попасть в результат.
Все слова должны быть приведены к нижнему регистру.
В результате каждое слово должно встречаться только один раз.

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

// Функция поиска всех множеств анаграмм по словарю
func FindAnagrams(words *[]string) *map[string][]string {
	anagramGroups := make(map[string][]string)	// временная мапа для групп анаграмм, где ключ - сортированное слово
	seenWords := make(map[string]bool)	// мапа слов, кот. уже были добавлены

	for _, word := range *words {
		// приводим к нижнемсу регистру
		lowerWord := strings.ToLower(word)

		// пропуск добавленного слова
		if seenWords[lowerWord] {
			continue
		}

		// создаем ключ
		key := sortString(lowerWord)

		// добавление слова в группу
		anagramGroups[key] = append(anagramGroups[key], lowerWord)

		// помечаем слово как добавленное
		seenWords[lowerWord] = true
	}

	// итоговая мапа без групп из одного эл-та
	result := make(map[string][]string)
	for _, group := range anagramGroups {
		if len(group) > 1 {
			sort.Strings(group) // сортируем группу
			result[group[0]] = group
		}
	}

	return &result
}

// Функция возвращает отсортированную строку для ключа анаграмм
func sortString(s string) string {
	runes := []rune(s)
	sort.Slice(runes, func(i, j int) bool {
		return runes[i] < runes[j]
	})
	return string(runes)
}

func main() {
	words := []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "аптек", "пятак", "Пятка"}
	result := FindAnagrams(&words)

	// Вывод результата
	for key, group := range *result {
		println(key + ": " + strings.Join(group, ", "))
	}
}
