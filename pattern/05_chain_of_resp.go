package main

/*
	Реализовать паттерн «цепочка вызовов».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Chain-of-responsibility_pattern
*/

/*
Паттерн "Цепочка вызовов" (Chain of Responsibility) — это поведенческий паттерн проектирования, который позволяет передавать запрос по цепочке обработчиков.
Каждый обработчик решает, может ли он обработать запрос, или передаёт его дальше по цепочке.

Применимость:
передача запроса по цепочке полезна, когда:
- существует несколько объектов, которые могут обработать запрос, и нужно дать им возможность решить, какой объект будет его обрабатывать.
- обработка запроса должна быть гибкой и многократной, без явного указания порядка обработки.

Плюсы:
- уменьшает зависимость между отправителем и получателем.
- позволяет динамически добавлять или изменять обработчиков без изменения кода отправителя.
- может помочь в обработке запросов в иерархическом порядке.

Минусы:
- может возникнуть ситуация, когда запрос не будет обработан, если в цепочке нет подходящего обработчика.
- сложность в отладке, если цепочка слишком длинная.

Примеры:
1. В системах, где запросы обрабатываются несколькими уровнями (например, в валидации данных или обработке событий).
2. В web-приложениях, где запросы проходят через несколько уровней (например, через фильтры, которые могут выполнять разные действия, такие как аутентификация, авторизация, логирование).
*/

import "fmt"

// Интерфейс обработчика
type Handler interface {
	Handle(amount int) bool  // Обработать запрос
	SetNext(handler Handler) // Установить следующий обработчик
}

// Базовый обработчик
type BaseHandler struct {
	next Handler
}

// Передача запроса дальше, если обработчик не справился
func (b *BaseHandler) Handle(amount int) bool {
	if b.next != nil {
		return b.next.Handle(amount)
	}
	fmt.Println("Запрос отклонен: сумма слишком большая.")
	return false
}

// Установка следующего обработчика
func (b *BaseHandler) SetNext(handler Handler) {
	b.next = handler
}

// Менеджер может утверждать до 10 000
type Manager struct {
	BaseHandler
}

func (m *Manager) Handle(amount int) bool {
	if amount <= 10000 {
		fmt.Printf("Менеджер утвердил бюджет: %d\n", amount)
		return true
	}
	fmt.Println("Менеджер: сумма слишком большая, передаю дальше...")
	return m.BaseHandler.Handle(amount)
}

// Директор может утверждать до 50 000
type Director struct {
	BaseHandler
}

func (d *Director) Handle(amount int) bool {
	if amount <= 50000 {
		fmt.Printf("Директор утвердил бюджет: %d\n", amount)
		return true
	}
	fmt.Println("Директор: сумма слишком большая, передаю дальше...")
	return d.BaseHandler.Handle(amount)
}

// Совет директоров может утвердить любую сумму
type Board struct {
	BaseHandler
}

func (b *Board) Handle(amount int) bool {
	fmt.Printf("Совет директоров утвердил бюджет: %d\n", amount)
	return true
}

// Создание цепочки обработчиков
func createApprovalChain() Handler {
	manager := &Manager{}
	director := &Director{}
	board := &Board{}

	manager.SetNext(director)
	director.SetNext(board)

	return manager // Возвращаем первого в цепочке
}

func main() {
	approver := createApprovalChain() // настройка цепочки

	// Разные суммы запроса
	fmt.Println("Запрос на 5 000:")
	approver.Handle(5000)

	fmt.Println("\nЗапрос на 30 000:")
	approver.Handle(30000)

	fmt.Println("\nЗапрос на 100 000:")
	approver.Handle(100000)
}
