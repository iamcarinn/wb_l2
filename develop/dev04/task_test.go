package main

import (
	"reflect"
	"testing"
)

func TestFindAnagrams(t *testing.T) {
	tests := []struct {
		name     string
		words    []string
		expected map[string][]string
	}{
		{
			name:  "Basic test with multiple groups",
			words: []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик"},
			expected: map[string][]string{
				"пятак":  {"пятак", "пятка", "тяпка"},
				"листок": {"листок", "слиток", "столик"},
			},
		},
		{
			name:  "Case sensitivity and duplicates",
			words: []string{"Пятак", "пятка", "Пятак", "пЯтка", "Тяпка"},
			expected: map[string][]string{
				"пятак": {"пятак", "пятка", "тяпка"},
			},
		},
		{
			name:  "No anagrams",
			words: []string{"мама", "папа", "дом"},
			expected: map[string][]string{
				// Пустой результат, так как нет групп больше одного элемента
			},
		},
		{
			name:  "Mixed anagrams and single words",
			words: []string{"пятак", "пятка", "тяпка", "дом", "слово", "ловос"},
			expected: map[string][]string{
				"пятак": {"пятак", "пятка", "тяпка"},
				"ловос": {"ловос", "слово"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := FindAnagrams(&test.words)

			// Сравниваем результат с ожидаемым
			if !reflect.DeepEqual(*result, test.expected) {
				t.Errorf("Expected %v, but got %v", test.expected, *result)
			}
		})
	}
}
