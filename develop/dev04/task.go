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

// FindAnagrams ищет множества анаграмм в словаре.
func FindAnagrams(words *[]string) *map[string][]string {
	anagramGroups := make(map[string][]string)
	seenWords := make(map[string]bool)

	for _, word := range *words {
		// Приводим слово к нижнему регистру
		lowerWord := strings.ToLower(word)

		// Пропускаем слово, если оно уже было добавлено
		if seenWords[lowerWord] {
			continue
		}

		// Создаем "ключ анаграмм" для текущего слова
		key := sortString(lowerWord)

		// Добавляем слово в соответствующую группу
		anagramGroups[key] = append(anagramGroups[key], lowerWord)

		// Помечаем слово как добавленное
		seenWords[lowerWord] = true
	}

	// Создаем результирующую мапу, исключая группы из одного элемента
	result := make(map[string][]string)
	for _, group := range anagramGroups {
		if len(group) > 1 {
			sort.Strings(group) // Сортируем группу
			result[group[0]] = group
		}
	}

	return &result
}

// sortString возвращает отсортированную строку для ключа анаграмм.
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
