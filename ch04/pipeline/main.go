package main

import (
	"fmt"
	"math/rand"
)

func main() {
	fmt.Println("--- Running pipeline1 ---")
	pipeline1()

	fmt.Println("\n--- Running repeat_take ---")
	repeat_take()

	fmt.Println("\n--- Running repeatFn_take ---")
	repeatFn_take()
}

func pipeline1() {
	generator := func(done <-chan any, integers ...int) <-chan int {
		intChan := make(chan int, len(integers))
		go func() {
			defer close(intChan)
			for _, i := range integers {
				select {
				case <-done:
					return
				case intChan <- i:
				}
			}
		}()
		return intChan
	}

	multiply := func(
		done <-chan any,
		intChan <-chan int,
		multiplier int,
	) <-chan int {
		multipliedChan := make(chan int)
		go func() {
			defer close(multipliedChan)
			for i := range intChan {
				select {
				case <-done:
					return
				case multipliedChan <- i * multiplier:
				}
			}
		}()
		return multipliedChan
	}

	add := func(
		done <-chan any,
		intChan <-chan int,
		additive int,
	) <-chan int {
		addedChan := make(chan int)
		go func() {
			defer close(addedChan)
			for i := range intChan {
				select {
				case <-done:
					return
				case addedChan <- i + additive:
				}
			}
		}()
		return addedChan
	}

	done := make(chan any)
	defer close(done)

	intChan := generator(done, 1, 2, 3, 4)
	pipeline := multiply(done, add(done, multiply(done, intChan, 2), 1), 2)

	for v := range pipeline {
		fmt.Println(v)
	}
}

var (
	repeat = func(
		done <-chan any,
		values ...any,
	) <-chan any {
		valueChan := make(chan any)
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

	repeatFn = func(
		done <-chan any,
		fn func() any,
	) <-chan any {
		valueChan := make(chan any)
		go func() {
			defer close(valueChan)
			for {
				select {
				case <-done:
					return
				case valueChan <- fn():
				}
			}
		}()
		return valueChan
	}

	take = func(
		done <-chan any,
		valueChan <-chan any,
		num int,
	) <-chan any {
		takeChan := make(chan any)
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

	toString = func(
		done <-chan any,
		valueChan <-chan any,
	) <-chan string {
		stringChan := make(chan string)
		go func() {
			defer close(stringChan)
			for v := range valueChan {
				select {
				case <-done:
					return
				case stringChan <- v.(string):
				}
			}
		}()
		return stringChan
	}
)

func repeat_take() {
	done := make(chan any)
	defer close(done)

	for num := range take(done, repeat(done, 1, 2), 10) {
		fmt.Printf("%v ", num)
	}
	fmt.Println()
}

func repeatFn_take() {
	done := make(chan any)
	defer close(done)

	rand := func() any { return rand.Int() }

	for num := range take(done, repeatFn(done, rand), 10) {
		fmt.Println(num)
	}
}
