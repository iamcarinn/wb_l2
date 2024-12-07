package main

import (
	"testing"
	"time"
	"github.com/beevik/ntp"
)

// Тест успешного получения времени с реальным запросом.
func TestGetTime_RealSuccess(t *testing.T) {
	actualTime, err := GetTime()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Проверяем, что время находится в пределах разумного диапазона.
	now := time.Now()
	if actualTime.Before(now.Add(-time.Hour)) || actualTime.After(now.Add(time.Hour)) {
		t.Errorf("expected time near now, got %v", actualTime)
	}
}

// Тест обработки ошибки.
func TestGetTime_RealError(t *testing.T) {
	_, err := ntp.Time("invalid.ntp.server")
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}
