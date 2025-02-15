package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
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

type Config struct {
	Port int `json:"port"`
}

func LoadConfig(filename string) (Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Config{}, fmt.Errorf("error opening config file: %v", err)
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return Config{}, fmt.Errorf("error decoding config file: %v", err)
	}

	return config, nil
}

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

// GetEventsForDay получает события на день
func (c *Calendar) GetEventsForDay(date time.Time) []Event {
	var events []Event
	for _, event := range c.events {
		if event.Date.Equal(date) {
			events = append(events, event)
		}
	}
	return events
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
func writeJSONResponse(w http.ResponseWriter, status int, response interface{}) {
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
			// др ошибки, http.StatusInternalServerError = 500
			writeJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "method not allowed"})
			return
		}

		// Парсим данные из запроса
		event, err := parseCreateEventRequest(r)
		if err != nil {
			// ошибка входных данных, http.StatusBadRequest = 400
			writeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		// Добавляем событие в календарь
		if err := calendar.CreateEvent(event); err != nil {
			// ошибка бизнес-логики, http.StatusServiceUnavailable = 503
			writeJSONResponse(w, http.StatusServiceUnavailable, map[string]string{"error": err.Error()})
			return
		}

		// Успешный ответ
		writeJSONResponse(w, http.StatusOK, map[string]string{"result": "event created"})
	}
}

// updateEventHandler возвращает обработчик для обновления событий
func updateEventHandler(calendar *Calendar) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса (должен быть POST)
		if r.Method != http.MethodPost {
			// др ошибки, http.StatusInternalServerError = 500
			writeJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "method not allowed"})
			return
		}

		// Парсим данные из запроса
		event, err := parseCreateEventRequest(r)
		if err != nil {
			// ошибка входных данных, http.StatusBadRequest = 400
			writeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		// Проверяем, существует ли событие с указанным ID
		existingEvent, exists := calendar.events[event.ID]
		if !exists {
			// ошибка бизнес-логики, http.StatusServiceUnavailable = 503
			writeJSONResponse(w, http.StatusServiceUnavailable, map[string]string{"error": fmt.Sprintf("event with ID %d does not exist", event.ID)})
			return
		}

		// Обновляем данные события
		existingEvent.Title = event.Title
		existingEvent.Description = event.Description
		existingEvent.Date = event.Date
		calendar.events[event.ID] = existingEvent

		// Успешный ответ
		writeJSONResponse(w, http.StatusOK, map[string]string{"result": "event updated"})
	}
}

func deleteEventHandler(calendar *Calendar) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса (должен быть POST)
		if r.Method != http.MethodPost {
			// др ошибки, http.StatusInternalServerError = 500
			writeJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "method not allowed"})
			return
		}

		// Парсим данные из запроса (только ID требуется для удаления)
		if err := r.ParseForm(); err != nil {
			// ошибка входных данных, http.StatusBadRequest = 400
			writeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid request format"})
			return
		}

		// Извлекаем и валидируем ID события
		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			// ошибка входных данных, http.StatusBadRequest = 400
			writeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid event ID"})
			return
		}

		// Проверяем, существует ли событие с указанным ID
		if _, exists := calendar.events[id]; !exists {
			// ошибка бизнес-логики, http.StatusServiceUnavailable = 503
			writeJSONResponse(w, http.StatusServiceUnavailable, map[string]string{"error": fmt.Sprintf("event with ID %d does not exist", id)})
			return
		}

		// Удаляем событие
		delete(calendar.events, id)

		// Успешный ответ
		writeJSONResponse(w, http.StatusOK, map[string]string{"result": "event deleted"})
	}
}

func getEventsForDayHandler(calendar *Calendar) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			// др ошибки, http.StatusInternalServerError = 500
			writeJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "method not allowed"})
			return
		}

		// Получаем параметр date из query
		dateStr := r.URL.Query().Get("date")
		if dateStr == "" {
			// ошибка входных данных, http.StatusBadRequest = 400
			writeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "missing date parameter"})
			return
		}

		// Разбираем дату
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			// ошибка входных данных, http.StatusBadRequest = 400
			writeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid date format, expected YYYY-MM-DD"})
			return
		}

		// Получаем события
		events := calendar.GetEventsForDay(date)

		// Отправляем ответ
		writeJSONResponse(w, http.StatusOK, map[string]interface{}{"result": events})
	}
}

func getEventsForWeekHandler(calendar *Calendar) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			// др ошибки, http.StatusInternalServerError = 500
			writeJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "method not allowed"})
			return
		}

		dateStr := r.URL.Query().Get("date")
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			// ошибка входных данных, http.StatusBadRequest = 400
			writeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid date format"})
			return
		}

		startOfWeek, endOfWeek := getWeekRange(date)

		var events []Event
		for d := startOfWeek; !d.After(endOfWeek); d = d.AddDate(0, 0, 1) {
			events = append(events, calendar.GetEventsForDay(d)...)	// ...распаковываем срез событий на день
		}

		writeJSONResponse(w, http.StatusOK, map[string]interface{}{"result": events})
	}
}
// getWeekRange возвращает начало (понедельник) и конец (воскресенье) недели для указанной даты.
func getWeekRange(date time.Time) (time.Time, time.Time) {
	// Определяем смещение от понедельника
	offset := int(date.Weekday()) - 1
	if offset < 0 { // Если день — воскресенье, делаем его седьмым днём недели
		offset = 6
	}

	// Определяем начало недели (понедельник)
	startOfWeek := date.AddDate(0, 0, -offset).Truncate(24 * time.Hour)

	// Конец недели — воскресенье (понедельник + 6 дней)
	endOfWeek := startOfWeek.AddDate(0, 0, 6)

	return startOfWeek, endOfWeek
}

func getEventsForMonthHandler(calendar *Calendar) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			// др ошибки, http.StatusInternalServerError = 500
			writeJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "method not allowed"})
			return
		}

		dateStr := r.URL.Query().Get("date")
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			// ошибка входных данных, http.StatusBadRequest = 400
			writeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid date format"})
			return
		}

		startOfMonth, endOfMonth := getMonthRange(date)

		var events []Event
		for d := startOfMonth; !d.After(endOfMonth); d = d.AddDate(0, 0, 1) {
			events = append(events, calendar.GetEventsForDay(d)...)	// ...распаковываем срез событий на день
		}

		writeJSONResponse(w, http.StatusOK, map[string]interface{}{"result": events})
	}
}

// getMonthRange возвращает начало (первый день) и конец (последний день) месяца для указанной даты.
func getMonthRange(date time.Time) (time.Time, time.Time) {
	// Определяем первый день месяца
	startOfMonth := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.UTC)

	// Определяем последний день месяца
	endOfMonth := startOfMonth.AddDate(0, 1, -1) // Следующий месяц - 1 день

	return startOfMonth, endOfMonth
}

// loggingMiddleware оборачивает обработчик и логирует детали каждого запроса.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// Логирование начала запроса
		log.Printf("Начало запроса: метод=%s, url=%s", r.Method, r.URL.String())

		// Вызов следующего обработчика
		next.ServeHTTP(w, r)

		// Логирование завершения запроса с длительностью обработки
		duration := time.Since(start)
		log.Printf("Завершение запроса: метод=%s, url=%s, длительность=%s", r.Method, r.URL.String(), duration)
	})
}


func main() {
	// Загружаем конфиг
	config, err := LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Преобразуем порт в строку
	portStr := strconv.Itoa(config.Port)

	// Проверяем, что порт начинается с ":"
	if !strings.HasPrefix(portStr, ":") {
		portStr = ":" + portStr
	}

	// Создаем новый календарь
	calendar := NewCalendar()

	// Создаем новый ServeMux и регистрируем обработчики
	mux := http.NewServeMux()
	mux.HandleFunc("/create_event", createEventHandler(calendar))
	mux.HandleFunc("/update_event", updateEventHandler(calendar))
	mux.HandleFunc("/delete_event", deleteEventHandler(calendar))
	mux.HandleFunc("/events_for_day", getEventsForDayHandler(calendar))
	mux.HandleFunc("/events_for_week", getEventsForWeekHandler(calendar))
	mux.HandleFunc("/events_for_month", getEventsForMonthHandler(calendar))

	// Оборачиваем мультиплексор в мидлвар для логирования запросов
	loggedMux := loggingMiddleware(mux)

	log.Printf("Запуск HTTP сервера на порту %s", portStr)
	if err := http.ListenAndServe(portStr, loggedMux); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
