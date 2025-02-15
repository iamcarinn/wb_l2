package main

/*
	Реализовать паттерн «строитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Builder_pattern
*/

/*
Паттерн "Строитель" (Builder) — это паттерн проектирования, который позволяет создавать сложные объекты пошагово, предоставляя гибкость в процессе их сборки.

Применимость:
строитель используется, когда необходимо создать сложный объект, но есть необходимость в гибкости для различных конфигураций этого объекта,
а также когда объект имеет множество компонентов, которые могут быть добавлены в разном порядке или с различными параметрами.

Плюсы:
- позволяет разделить процесс создания объекта на независимые шаги
- улучшает читаемость кода, так как каждый этап создания объекта инкапсулирован в отдельном методе
- упрощает создание объектов с множеством параметров
- позволяет легко модифицировать процесс создания объекта без изменения самого объекта

Минусы:
- увеличивает количество кода из-за необходимости создания отдельных классов для каждого строителя
- может быть излишним для объектов с простыми конфигурациями, где порядок действий не имеет значения

Примеры:
1. Создание сложных документов (например, PDF-файлов), где необходимо поэтапно добавлять различные элементы.
2. Настройка приложений — создание настроек для приложений с множеством опций (например, для банковского ПО).

*/


import "fmt"

// Пицца - продукт, который мы будем строить
type Pizza struct {
	dough     string
	sauce     string
	cheese    string
	toppings  []string
}

// Строитель пиццы
type PizzaBuilder struct {
	pizza *Pizza
}

// Конструктор строитель
func NewPizzaBuilder() *PizzaBuilder {
	return &PizzaBuilder{
		pizza: &Pizza{},
	}
}

// Устанавливаем тесто
func (b *PizzaBuilder) SetDough(dough string) *PizzaBuilder {
	b.pizza.dough = dough
	return b
}

// Устанавливаем соус
func (b *PizzaBuilder) SetSauce(sauce string) *PizzaBuilder {
	b.pizza.sauce = sauce
	return b
}

// Устанавливаем сыр
func (b *PizzaBuilder) SetCheese(cheese string) *PizzaBuilder {
	b.pizza.cheese = cheese
	return b
}

// Добавляем ингредиенты
func (b *PizzaBuilder) AddTopping(topping string) *PizzaBuilder {
	b.pizza.toppings = append(b.pizza.toppings, topping)
	return b
}

// Получаем готовую пиццу
func (b *PizzaBuilder) Build() *Pizza {
	return b.pizza
}

func main() {
	// Строим пиццу поэтапно
	builder := NewPizzaBuilder()
	pizza := builder.SetDough("Тонкое тесто").
		SetSauce("Томатный соус").
		SetCheese("Моцарелла").
		AddTopping("Пепперони").
		AddTopping("Грибы").
		Build()

	// Выводим информацию о собранной пицце
	fmt.Printf("Пицца собрана:\nТесто: %s\nСоус: %s\nСыр: %s\nИнгредиенты: %v\n",
		pizza.dough, pizza.sauce, pizza.cheese, pizza.toppings)
}
