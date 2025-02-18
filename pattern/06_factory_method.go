package main

/*
	Реализовать паттерн «фабричный метод».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Factory_method_pattern
*/

/*
Паттерн "Фабричный метод" (Factory Method) — это паттерн проектирования, который определяет общий интерфейс для создания объектов,
но делегирует создание конкретных экземпляров подклассам.

Применимость:
Фабричный метод используется, когда:
- код должен работать с разными типами объектов, но заранее неизвестно, какие именно понадобятся
- важно централизовать процесс создания объектов
- создание объектов связано с логикой, которую не стоит дублировать в разных местах кода

Плюсы:
- изолирует код от конкретных реализаций объектов
- упрощает поддержку и расширение системы
- централизует создание объектов, облегчая управление зависимостями

Минусы:
- увеличивает сложность кода из-за дополнительных классов и интерфейсов
- иногда избыточен, если создание объекта простое

Примеры:
1. Работа с базами данных: фабрика может создавать соединения к MySQL, PostgreSQL, SQLite в зависимости от настроек.
2. Система уведомлений: в зависимости от типа (SMS, Email, Push) создается нужный сервис.
*/


import "fmt"

// Интерфейс транспортного средства
type Transport interface {
	Deliver()
}

// Конкретная реализация - Грузовик
type Truck struct{}

func (t *Truck) Deliver() {
	fmt.Println("Доставка по земле грузовиком.")
}

// Конкретная реализация - Корабль
type Ship struct{}

func (s *Ship) Deliver() {
	fmt.Println("Доставка по морю кораблем.")
}

// Фабричный метод (интерфейс)
type Logistics interface {
	CreateTransport() Transport
}

// Реализация фабрики для грузовиков
type RoadLogistics struct{}

func (r *RoadLogistics) CreateTransport() Transport {
	return &Truck{}
}

// Реализация фабрики для кораблей
type SeaLogistics struct{}

func (s *SeaLogistics) CreateTransport() Transport {
	return &Ship{}
}

func main() {
	// Доставка по суше
	var logistics Logistics = &RoadLogistics{}
	transport := logistics.CreateTransport()
	transport.Deliver()

	// Доставка по морю
	logistics = &SeaLogistics{}
	transport = logistics.CreateTransport()
	transport.Deliver()
}
