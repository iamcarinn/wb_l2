package main

import (
	"testing"
	"time"
)

// Функция для создания сигнального канала, кот. закроется через заданное время.
func createSig(after time.Duration) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
	}()
	return c
}

// Тест с одним каналом
func TestOrSingleChannel(t *testing.T) {
	sig := createSig(1 * time.Second)
	start := time.Now()

	<-or(sig)

	elapsed := time.Since(start)
	if elapsed < 1*time.Second || elapsed > 1*time.Second+100*time.Millisecond {
		t.Errorf("Expected ~1s elapsed time, got %v", elapsed)
	}
}

// Тест с несколькими каналами
func TestOrMultipleChannels(t *testing.T) {
	start := time.Now()

	<-or(
		createSig(3*time.Second),
		createSig(1*time.Second),
		createSig(5*time.Second),
		createSig(2*time.Second),
	)

	elapsed := time.Since(start)
	if elapsed < 1*time.Second || elapsed > 1*time.Second+100*time.Millisecond {
		t.Errorf("Expected ~1s elapsed time, got %v", elapsed)
	}
}

// Тест без каналов
func TestOrNoChannels(t *testing.T) {
	start := time.Now()

	<-or()

	elapsed := time.Since(start)
	if elapsed > 100*time.Millisecond {
		t.Errorf("Expected immediate closure, but got %v", elapsed)
	}
}

// Тест с немедленным закрытием
func TestOrImmediateClose(t *testing.T) {
	immediate := make(chan interface{})
	close(immediate)

	start := time.Now()
	<-or(
		immediate,
		createSig(1*time.Second),
		createSig(2*time.Second),
	)

	elapsed := time.Since(start)
	if elapsed > 100*time.Millisecond {
		t.Errorf("Expected immediate closure, but got %v", elapsed)
	}
}
