package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestCreateEvent(t *testing.T) {
	// Создаем новый календарь
	calendar := NewCalendar()

	// Создаем новый обработчик для создания события
	handler := createEventHandler(calendar)

	// Данные запроса
	data := "id=1&title=Meeting&description=Project discussion&date=2025-02-15"

	// Создаем новый POST-запрос
	req, err := http.NewRequest("POST", "/create_event", bytes.NewBufferString(data))
	if err != nil {
		t.Fatalf("Не удалось создать запрос: %v", err)
	}

	// Устанавливаем корректный Content-Type
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Создаем новый ResponseRecorder для записи ответа
	rr := httptest.NewRecorder()

	// Выполняем обработку запроса
	handler.ServeHTTP(rr, req)

	// Проверяем статус код
	if rr.Code != http.StatusOK {
		t.Errorf("Ожидался статус 200, получили %v", rr.Code)
	}

	// Разбираем JSON-ответ
	var response map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Ошибка при разборе JSON-ответа: %v", err)
	}

	// Ожидаемый результат
	expected := map[string]string{"result": "event created"}

	// Проверяем содержимое ответа
	if !reflect.DeepEqual(response, expected) {
		t.Errorf("Ожидался ответ %+v, получили %+v", expected, response)
	}
}

func TestUpdateEvent(t *testing.T) {
	calendar := NewCalendar()
	calendar.CreateEvent(Event{ID: 1, Title: "Old", Description: "Old desc", Date: time.Now()})
	handler := updateEventHandler(calendar)
	data := "id=1&title=New&description=New desc&date=2025-02-16"
	req, err := http.NewRequest("POST", "/update_event", bytes.NewBufferString(data))
	if err != nil {
		t.Fatalf("Не удалось создать запрос: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Ожидался статус 200, получили %v", rr.Code)
	}
	var response map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Ошибка при разборе JSON-ответа: %v", err)
	}
	expected := map[string]string{"result": "event updated"}
	if !reflect.DeepEqual(response, expected) {
		t.Errorf("Ожидался ответ %+v, получили %+v", expected, response)
	}
}

func TestDeleteEvent(t *testing.T) {
	calendar := NewCalendar()
	calendar.CreateEvent(Event{ID: 1, Title: "Meeting", Description: "Discussion", Date: time.Now()})
	handler := deleteEventHandler(calendar)
	data := "id=1"
	req, err := http.NewRequest("POST", "/delete_event", bytes.NewBufferString(data))
	if err != nil {
		t.Fatalf("Не удалось создать запрос: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Ожидался статус 200, получили %v", rr.Code)
	}
	var response map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Ошибка при разборе JSON-ответа: %v", err)
	}
	expected := map[string]string{"result": "event deleted"}
	if !reflect.DeepEqual(response, expected) {
		t.Errorf("Ожидался ответ %+v, получили %+v", expected, response)
	}
}

func TestGetEventsForDay(t *testing.T) {
	calendar := NewCalendar()
	eventDate := time.Date(2025, 2, 15, 0, 0, 0, 0, time.UTC)
	calendar.CreateEvent(Event{ID: 1, Title: "Meeting", Description: "Discussion", Date: eventDate})
	handler := getEventsForDayHandler(calendar)
	req, err := http.NewRequest("GET", "/events_for_day?date=2025-02-15", nil)
	if err != nil {
		t.Fatalf("Не удалось создать запрос: %v", err)
	}
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Ожидался статус 200, получили %v", rr.Code)
	}
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Ошибка при разборе JSON-ответа: %v", err)
	}
	if _, ok := response["result"]; !ok {
		t.Errorf("Ожидался ключ 'result' в ответе, получили %+v", response)
	}
}

func TestGetEventsForWeek(t *testing.T) {
	calendar := NewCalendar()

	// Добавляем события: два в одной неделе, одно в другой
	calendar.CreateEvent(Event{ID: 1, Title: "Meeting", Description: "Team sync", Date: time.Date(2025, 2, 10, 0, 0, 0, 0, time.UTC)}) // Понедельник
	calendar.CreateEvent(Event{ID: 2, Title: "Workshop", Description: "Go coding", Date: time.Date(2025, 2, 12, 0, 0, 0, 0, time.UTC)}) // Среда
	calendar.CreateEvent(Event{ID: 3, Title: "Conference", Description: "Tech talk", Date: time.Date(2025, 2, 20, 0, 0, 0, 0, time.UTC)}) // Другая неделя

	handler := getEventsForWeekHandler(calendar)

	// Запрос на получение событий недели, начиная с 10 февраля
	req, err := http.NewRequest("GET", "/events_for_week?date=2025-02-10", nil)
	if err != nil {
		t.Fatalf("Не удалось создать запрос: %v", err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Ожидался статус 200, получили %v", rr.Code)
	}

	// Разбираем JSON-ответ
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Ошибка при разборе JSON-ответа: %v", err)
	}

	// Проверяем, что есть ключ "result"
	events, ok := response["result"].([]interface{})
	if !ok {
		t.Fatalf("Ожидался массив событий в поле 'result', получили %+v", response)
	}

	// Должны быть только два события в этой неделе
	if len(events) != 2 {
		t.Errorf("Ожидалось 2 события, получено %d", len(events))
	}
}


func TestGetEventsForMonth(t *testing.T) {
	calendar := NewCalendar()

	// Добавляем события: два в одном месяце, одно в другом
	calendar.CreateEvent(Event{ID: 4, Title: "Summer Festival", Description: "Outdoor fun", Date: time.Date(2025, 6, 10, 0, 0, 0, 0, time.UTC)}) // 10 июня
	calendar.CreateEvent(Event{ID: 5, Title: "Beach Party", Description: "Sunset celebration", Date: time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC)}) // 15 июня
	calendar.CreateEvent(Event{ID: 6, Title: "Music Concert", Description: "Live performance", Date: time.Date(2025, 7, 5, 0, 0, 0, 0, time.UTC)}) // 5 июля

	handler := getEventsForMonthHandler(calendar)

	// Запрос на получение событий месяца, июня 2025 года
	req, err := http.NewRequest("GET", "/events_for_month?date=2025-06-01", nil)
	if err != nil {
		t.Fatalf("Не удалось создать запрос: %v", err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Ожидался статус 200, получили %v", rr.Code)
	}

	// Разбираем JSON-ответ
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Ошибка при разборе JSON-ответа: %v", err)
	}

	// Проверяем, что есть ключ "result"
	events, ok := response["result"].([]interface{})
	if !ok {
		t.Fatalf("Ожидался массив событий в поле 'result', получили %+v", response)
	}

	// Должны быть только два события в июне
	if len(events) != 2 {
		t.Errorf("Ожидалось 2 события в июне, получено %d", len(events))
	}
}
