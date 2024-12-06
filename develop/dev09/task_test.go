package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// Тестируем функцию validateURL
func TestValidateURL(t *testing.T) {
	tests := []struct {
		url     string
		wantErr bool
	}{
		{"http://example.com", false},  // Валидный URL
		{"https://example.com", false}, // Валидный URL
		{"ftp://example.com", true},    // Невалидный URL
		{"", true},                     // Пустой URL
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			err := validateURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateURL(%q) = %v, wantErr %v", tt.url, err, tt.wantErr)
			}
		})
	}
}

// Тестируем функцию fileNameParse
func TestFileNameParse(t *testing.T) {
	tests := []struct {
		site     string
		expected string
	}{
		{"http://example.com", "example.com.html"},
		{"https://www.example.com", "www.example.com.html"},
		{"http://example.com/page", "example.com.html"},
		{"", "index.html"},
	}

	for _, tt := range tests {
		t.Run(tt.site, func(t *testing.T) {
			got := fileNameParse(tt.site)
			if got != tt.expected {
				t.Errorf("fileNameParse(%q) = %q, want %q", tt.site, got, tt.expected)
			}
		})
	}
}

// Тестируем функцию downloadHTML
func TestDownloadHTML(t *testing.T) {
	// Запускаем HTTP-сервер для тестов
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<html><body>Test Page</body></html>"))
	}))
	defer server.Close()

	// Заменяем реальный сайт на адрес тестового сервера
	site := server.URL

	// Вызываем функцию для скачивания
	err := downloadHTML(site)
	if err != nil {
		t.Fatalf("downloadHTML(%q) failed: %v", site, err)
	}

	// Проверяем, что файл был создан
	fileName := fileNameParse(site)
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		t.Fatalf("Expected file %s to be created, but it does not exist", fileName)
	}

	// Удаляем файл после теста
	os.Remove(fileName)
}
