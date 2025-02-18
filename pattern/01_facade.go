package main

/*
	Реализовать паттерн «фасад».
Объяснить применимость паттерна, его плюсы и минусы,а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Facade_pattern
*/

/*
Паттерн "Фасад" (Facade) — это паттерн проектирования, который предоставляет простой интерфейс для работы с сложной системой.

Применимость:
фасад используется, когда система имеет сложную структуру с множеством подсистем,
и необходимо скрыть её детали, предоставив простой и понятный интерфейс.

Плюсы:
- упрощение взаимодействия с системой
- сокрытие сложности системы от пользователя

Минусы:
- может скрыть важные детали, что делает систему менее гибкой для расширения в будущем

Пример:
1. Умные дома — фасад скрывает сложность управления освещением, отоплением и безопасностью через единый интерфейс.
2. Системы платежей — фасад упрощает взаимодействие с различными платежными шлюзами, скрывая детали обработки транзакций.
*/

import "fmt"

// Фасад системы заказа в ресторане
type RestaurantFacade struct {
	kitchen    *Kitchen
	serving    *Serving
	payment    *Payment
	notification *Notification
}

// Конструктор фасада
func NewRestaurantFacade() *RestaurantFacade {
	return &RestaurantFacade{
		kitchen:      &Kitchen{},
		serving:      &Serving{},
		payment:      &Payment{},
		notification: &Notification{},
	}
}

// Оформление заказа
func (r *RestaurantFacade) PlaceOrder(dish string) {
	fmt.Println("Обработка заказа...")
	r.kitchen.PrepareDish(dish)
	r.serving.ServeDish(dish)
	r.payment.ProcessPayment()
	r.notification.SendConfirmation()
	fmt.Println("Заказ оформлен успешно!")
}

// Подсистема кухня
type Kitchen struct{}

func (k *Kitchen) PrepareDish(dish string) {
	fmt.Printf("Кухня: приготовлено блюдо \"%s\"\n", dish)
}

// Подсистема подача
type Serving struct{}

func (s *Serving) ServeDish(dish string) {
	fmt.Printf("Подача: подано блюдо \"%s\"\n", dish)
}

// Подсистема оплата
type Payment struct{}

func (p *Payment) ProcessPayment() {
	fmt.Println("Оплата: обработана транзакция")
}

// Подсистема уведомления
type Notification struct{}

func (n *Notification) SendConfirmation() {
	fmt.Println("Уведомление: отправлено подтверждение заказа")
}

func main() {
	// фасад для ресторана
	restaurant := NewRestaurantFacade()

	// оформление заказа
	restaurant.PlaceOrder("Паста карбонара")
}
