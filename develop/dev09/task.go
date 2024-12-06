package main

/*
=== Утилита wget ===

Реализовать утилиту wget с возможностью скачивать сайты целиком

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var site string // переменная для хранения URL

func init() {
	flag.StringVar(&site, "s", "", "Site to download") // Флаг -s для URL
}

// Функция проверяет URL
func validateURL(url string) error {
	if url == "" {
		return fmt.Errorf("please provide a site URL using -s flag") // Если URL не указан, завершаем с ошибкой
	}

	if !strings.HasPrefix(url, "http") {
		return fmt.Errorf("invalid URL: must start with http or https") // Проверяем, что URL начинается с http или https
	}
	return nil
}

// Функция возвращает имя файла для сохранения HTML
func fileNameParse(site string) string {
	parts := strings.Split(site, "/")
	if len(parts) > 2 {
		return parts[2] + ".html" // исп-ем домен как имя
	}
	return "index.html" // если домен не найден, то исп-ем "index.html"
}

// Функция скачивает HTML страницу
func downloadHTML(site string) error {
	// выполняем http GET-запрос для получения страницы
	resp, err := http.Get(site)
	if err != nil {
		return fmt.Errorf("error downloading site: %v", err)
	}
	defer resp.Body.Close() // закрываем тело ответа

	// создаем файл для сохранения HTML
	fileName := fileNameParse(site)
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	// копируем содержимое страницы в файл
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("error saving HTML: %v", err)
	}

	log.Printf("HTML saved as %s\n", fileName)
	return nil
}

// run запускает выполнение программы
func run() {
	flag.Parse()
	// Проверяем URL
	if err := validateURL(site); err != nil {
		log.Fatal(err)
	}

	// Скачиваем HTML
	if err := downloadHTML(site); err != nil {
		log.Fatal(err)
	}
}

func main() {
	run()
}
