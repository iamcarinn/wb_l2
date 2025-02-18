package main

/*
	Реализовать паттерн «комманда».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Command_pattern
*/

/*
Паттерн "Команда" (Command) — это паттерн проектирования, который превращает запросы в объекты, позволяя передавать
их как параметры, ставить их в очередь, логировать их и поддерживать отмену операций.

Применимость:
передача запроса как объекта полезна, когда:
- необходимо параметризовать объекты действиями (например, передавать команды в пульт управления).
- нужно реализовать систему, где команды можно отменить или повторить.
- задачи требуют очереди или логирования.

Плюсы:
- позволяет отменить/повторить действия.
- упрощает реализацию операций, требующих последовательного выполнения.
- позволяет легко добавлять новые команды без изменения существующего кода.

Минусы:
- может увеличить сложность программы, добавив новые классы для команд.
- может затруднить отладку, так как нужно отслеживать больше объектов команд.

Примеры:
1. В GUI-программах: кнопка "Отменить" или "Повторить" — команды, которые можно отменить или повторить.
2. В системах с большим количеством операций, где команды должны быть сохранены или обработаны в очереди.
*/


import "fmt"

// Получатель (ToDoList)
type ToDoList struct {
	tasks []string
}

func (t *ToDoList) Add(task string) {
	t.tasks = append(t.tasks, task)
	fmt.Println("Задача добавлена:", task)
}

func (t *ToDoList) Remove(task string) {
	for i, tsk := range t.tasks {
		if tsk == task {
			t.tasks = append(t.tasks[:i], t.tasks[i+1:]...)
			fmt.Println("Задача удалена:", task)
			return
		}
	}
}

func (t *ToDoList) List() {
	fmt.Println("Список задач:")
	for _, task := range t.tasks {
		fmt.Println("-", task)
	}
}

// Интерфейс команды
type Command interface {
	Execute()
}

// Конкретная команда: Добавление задачи
type AddTaskCommand struct {
	list *ToDoList
	task string
}

func (c *AddTaskCommand) Execute() {
	c.list.Add(c.task)
}

// Конкретная команда: Удаление задачи
type RemoveTaskCommand struct {
	list *ToDoList
	task string
}

func (c *RemoveTaskCommand) Execute() {
	c.list.Remove(c.task)
}

// Пульт управления
type TaskManager struct {
	command Command
}

func (tm *TaskManager) SetCommand(cmd Command) {
	tm.command = cmd
}

func (tm *TaskManager) PressButton() {
	tm.command.Execute()
}

func main() {
	// Создаём список дел и пульт
	list := &ToDoList{}
	manager := &TaskManager{}

	// Добавляем задачи
	manager.SetCommand(&AddTaskCommand{list, "Купить продукты"})
	manager.PressButton()
	manager.SetCommand(&AddTaskCommand{list, "Убрать квартиру"})
	manager.PressButton()
	manager.SetCommand(&AddTaskCommand{list, "Погулять с собакой"})
	manager.PressButton()
	manager.SetCommand(&AddTaskCommand{list, "Сделать домашку по Go"})
	manager.PressButton()

	// Выводим текущий список задач
	list.List()

	// Удаляем одну задачу
	manager.SetCommand(&RemoveTaskCommand{list, "Погулять с собакой"})
	manager.PressButton()

	// Выводим обновлённый список задач
	list.List()
}

