package main

/*
	Реализовать паттерн «стратегия».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Strategy_pattern
*/

/*
Паттерн "Стратегия" (Strategy) — это паттерн проектирования, который позволяет выбирать алгоритм поведения во время выполнения программы, не меняя код клиента.

Применимость:
стратегия используется, когда необходимо переключаться между разными вариантами поведения (алгоритмами) во время выполнения программы.

Плюсы:
- позволяет легко добавлять новые стратегии без изменения основного кода
- упрощает тестирование, так как стратегии изолированы друг от друга
- делает код более гибким за счет инкапсуляции алгоритмов

Минусы:
- увеличивает количество кода за счет создания множества классов для разных стратегий
- усложняет понимание кода, если стратегий слишком много

Пример:
1. Сортировка данных — можно менять алгоритмы сортировки (быстрая сортировка, сортировка пузырьком) в зависимости от объема данных.
2. Оплата товаров — пользователь может выбрать разные способы оплаты (карта, PayPal, криптовалюта) без изменения логики магазина.
*/

import "fmt"

// Интерфейс стратегии сортировки
type SortingStrategy interface {
	Sort([]int) []int
	Name() string // метод получения имени стратегии
}

// Реализация сортировки пузырьком
type BubbleSort struct{}

func (b *BubbleSort) Sort(data []int) []int {
	n := len(data)
	arr := append([]int(nil), data...) // Копируем массив, чтобы не менять оригинал
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
			}
		}
	}
	return arr
}

func (b *BubbleSort) Name() string {
	return "Пузырьковая сортировка"
}

// Реализация быстрой сортировки
type QuickSort struct{}

func (q *QuickSort) Sort(data []int) []int {
	if len(data) < 2 {
		return data
	}
	arr := append([]int(nil), data...) // Копируем массив
	pivot := arr[0]
	var left, right []int
	for _, v := range arr[1:] {
		if v < pivot {
			left = append(left, v)
		} else {
			right = append(right, v)
		}
	}
	return append(append(q.Sort(left), pivot), q.Sort(right)...)
}

func (q *QuickSort) Name() string {
	return "Быстрая сортировка"
}

// Контекст — список чисел с возможностью выбора стратегии сортировки
type NumberList struct {
	numbers  []int
	strategy SortingStrategy
}

// Установка стратегии сортировки
func (n *NumberList) SetSortingStrategy(strategy SortingStrategy) {
	n.strategy = strategy
}

// Сортировка списка с использованием выбранной стратегии
func (n *NumberList) SortNumbers() {
	if n.strategy == nil {
		fmt.Println("Ошибка: стратегия сортировки не выбрана!")
		return
	}
	n.numbers = n.strategy.Sort(n.numbers)
	fmt.Printf("Метод: %s | Отсортированный список: %v\n", n.strategy.Name(), n.numbers)
}

func main() {
	data := []int{42, 10, 3, 7, 18, 25}
	list := &NumberList{numbers: data}

	// Пузырьковая сортировка
	list.SetSortingStrategy(&BubbleSort{})
	list.SortNumbers()

	// Быстрая сортировка
	list.SetSortingStrategy(&QuickSort{})
	list.SortNumbers()
}
