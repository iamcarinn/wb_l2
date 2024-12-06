package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

/*
=== HTTP server ===

Реализовать HTTP сервер для работы с календарем. В рамках задания необходимо работать строго со стандартной HTTP библиотекой.
В рамках задания необходимо:
	1. Реализовать вспомогательные функции для сериализации объектов доменной области в JSON.
	2. Реализовать вспомогательные функции для парсинга и валидации параметров методов /create_event и /update_event.
	3. Реализовать HTTP обработчики для каждого из методов API, используя вспомогательные функции и объекты доменной области.
	4. Реализовать middleware для логирования запросов
Методы API: POST /create_event POST /update_event POST /delete_event GET /events_for_day GET /events_for_week GET /events_for_month
Параметры передаются в виде www-url-form-encoded (т.е. обычные user_id=3&date=2019-09-09).
В GET методах параметры передаются через queryString, в POST через тело запроса.
В результате каждого запроса должен возвращаться JSON документ содержащий либо {"result": "..."} в случае успешного выполнения метода,
либо {"error": "..."} в случае ошибки бизнес-логики.

В рамках задачи необходимо:
	1. Реализовать все методы.
	2. Бизнес логика НЕ должна зависеть от кода HTTP сервера.
	3. В случае ошибки бизнес-логики сервер должен возвращать HTTP 503. В случае ошибки входных данных (невалидный int например) сервер должен возвращать HTTP 400. В случае остальных ошибок сервер должен возвращать HTTP 500. Web-сервер должен запускаться на порту указанном в конфиге и выводить в лог каждый обработанный запрос.
	4. Код должен проходить проверки go vet и golint.
*/

type Event struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
}

type Calendar struct {
	events map[int]Event
}

// NewCalendar создает новый календарь
func NewCalendar() *Calendar {
	return &Calendar{events: make(map[int]Event)}
}

// CreateEvent добавляет событие в календарь
func (c *Calendar) CreateEvent(event Event) error {
	if _, exists := c.events[event.ID]; exists {
		return fmt.Errorf("event with ID %d already exists", event.ID)
	}
	c.events[event.ID] = event
	return nil
}

// parseCreateEventRequest парсит данные из POST-запроса и возвращает объект Event
func parseCreateEventRequest(r *http.Request) (Event, error) {
	// Парсим тело запроса как форму
	if err := r.ParseForm(); err != nil {
		return Event{}, fmt.Errorf("невалидный формат запроса: %v", err)
	}

	// Извлекаем и валидируем ID события
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		return Event{}, fmt.Errorf("невалидный формат ID: %v", err)
	}

	// Извлекаем заголовок и описание события
	title := r.FormValue("title")
	description := r.FormValue("description")

	// Извлекаем и валидируем дату события
	dateStr := r.FormValue("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return Event{}, fmt.Errorf("невалидный формат даты: %v", err)
	}

	// Возвращаем объект Event
	return Event{
		ID:          id,
		Title:       title,
		Description: description,
		Date:        date,
	}, nil
}

// writeJSONResponse формирует и отправляет JSON-ответ
func writeJSONResponse(w http.ResponseWriter, status int, response map[string]string) {
	// Устанавливаем заголовки ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// Кодируем и отправляем JSON-ответ
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Ошибка при формировании ответа: %v", err)
	}
}

// createEventHandler возвращает обработчик для создания событий
func createEventHandler(calendar *Calendar) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса (должен быть POST)
		if r.Method != http.MethodPost {
			writeJSONResponse(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}

		// Парсим данные из запроса
		event, err := parseCreateEventRequest(r)
		if err != nil {
			writeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		// Добавляем событие в календарь
		if err := calendar.CreateEvent(event); err != nil {
			writeJSONResponse(w, http.StatusServiceUnavailable, map[string]string{"error": err.Error()})
			return
		}

		// Успешный ответ
		writeJSONResponse(w, http.StatusOK, map[string]string{"result": "event created"})
	}
}

func main() {
	// Создаем новый календарь
	calendar := NewCalendar()

	// Регистрируем обработчик для создания событий
	http.HandleFunc("/create_event", createEventHandler(calendar))

	// Запускаем HTTP сервер
	port := ":8080"
	log.Printf("Запуск HTTP сервера на порту %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
