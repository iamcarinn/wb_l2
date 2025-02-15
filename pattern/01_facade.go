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
в библиотеке для работы с базой данных фасад может предоставить методы для работы с запросами,
скрывая сложности соединения с сервером и управлением транзакциями.
*/

import "fmt"

// Фасад умный дом
type SmartHomeFacade struct {
	lights     *Lights
	curtains   *Curtains
	thermostat *Thermostat
	alarm      *Alarm
}

// Конструктор фасада
func NewSmartHomeFacade() *SmartHomeFacade {
	return &SmartHomeFacade{
		lights:     &Lights{},
		curtains:   &Curtains{},
		thermostat: &Thermostat{},
		alarm:      &Alarm{},
	}
}

// Включение режима "Спокойной ночи"
func (s *SmartHomeFacade) GoodNightMode() {
	fmt.Println("Активация режима `Спокойной ночи`...")
	s.lights.Off()
	s.curtains.Close()
	s.thermostat.SetTemperature(20) // Устанавливаем комфортную температуру
	s.alarm.Activate()
	fmt.Println("Режим `Спокойной ночи` активирован!")
}

// Подсистема освещение
type Lights struct{}

func (l *Lights) Off() {
	fmt.Println("Свет: выключен")
}

// Подсистема шторы
type Curtains struct{}

func (c *Curtains) Close() {
	fmt.Println("Шторы: закрыты")
}

// Подсистема кондиционер
type Thermostat struct{}

func (t *Thermostat) SetTemperature(temp int) {
	fmt.Printf("Кондиционер: Установлена температура %d°C\n", temp)
}

// Подсистема сигнализация
type Alarm struct{}

func (a *Alarm) Activate() {
	fmt.Println("Сигнализация: активирована")
}

func main() {
	// фасад для умного дома
	smartHome := NewSmartHomeFacade()

	// режим "Спокойной ночи"
	smartHome.GoodNightMode()
}
