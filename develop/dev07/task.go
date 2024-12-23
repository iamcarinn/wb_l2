package main

import (
	"fmt"
	"time"
)

/*
=== Or channel ===

Реализовать функцию, которая будет объединять один или более done каналов в single канал если один из его составляющих каналов закроется.
Одним из вариантов было бы очевидно написать выражение при помощи select, которое бы реализовывало эту связь,
однако иногда неизестно общее число done каналов, с которыми вы работаете в рантайме.
В этом случае удобнее использовать вызов единственной функции, которая, приняв на вход один или более or каналов, реализовывала весь функционал.

Определение функции:
var or func(channels ...<- chan interface{}) <- chan interface{}

Пример использования функции:
sig := func(after time.Duration) <- chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
}()
return c
}

start := time.Now()
<-or (
	sig(2*time.Hour),
	sig(5*time.Minute),
	sig(1*time.Second),
	sig(1*time.Hour),
	sig(1*time.Minute),
)

fmt.Printf(“fone after %v”, time.Since(start))
*/

// Функция объединяет несколько каналов в один. Возвращает канал, который закрывается, как только закрывается любой из входных каналов.
func or(channels ...<-chan interface{}) <-chan interface{} {
	// если каналы не переданы, возвращаем закрытый канал
	if len(channels) == 0 {
		closedChan := make(chan interface{})
		close(closedChan)
		return closedChan
	}

	// если передан один канал, возвращаем его
	if len(channels) == 1 {
		return channels[0]
	}

	// если передано больше одного канала, то рекурсивно объединяем
	orDone := make(chan interface{})
	go func() {
		defer close(orDone)

		// ожидаем закрытие любого канала
		switch len(channels) {
		case 2:
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		default:
			mid := len(channels) / 2
			select {
			case <-or(channels[:mid]...):
			case <-or(channels[mid:]...):
			}
		}
	}()
	return orDone
}

func main() {
	// Функция для сигнального канала, кот. закроется ч/з заданное время.
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	// or() с несколькими каналами.
	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)

	fmt.Printf("Done after %v\n", time.Since(start))
}
