package main

import "testing"

func BenchmarkGeneric(b *testing.B) {
	done := make(chan any)
	defer close(done)

	b.ResetTimer()
	for range toString(done, take(done, repeat(done, "a"), b.N)) {
	}
}

func BenchmarkTyped(b *testing.B) {
	repeat := func(
		done <-chan any,
		values ...string,
	) <-chan string {
		valueChan := make(chan string)
		go func() {
			defer close(valueChan)
			for {
				for _, v := range values {
					select {
					case <-done:
						return
					case valueChan <- v:
					}
				}
			}
		}()
		return valueChan
	}

	take := func(
		done <-chan any,
		valueChan <-chan string,
		num int,
	) <-chan string {
		takeChan := make(chan string)
		go func() {
			defer close(takeChan)
			for i := 0; i < num; i++ {
				select {
				case <-done:
					return
				case takeChan <- <-valueChan:
				}
			}
		}()
		return takeChan
	}

	done := make(chan any)
	defer close(done)

	b.ResetTimer()
	for range take(done, repeat(done, "a"), b.N) {
	}
}
