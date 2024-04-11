package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Println("--- Running wg1 ---")
	wg1()

	fmt.Println("--- Running wg2 ---")
	wg2()
}

func wg1() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("1st goroutine sleeping...")
		time.Sleep(1 * time.Nanosecond)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("2nd goroutine sleeping...")
		time.Sleep(2 * time.Nanosecond)
	}()

	wg.Wait()
	fmt.Println("All goroutines complete.")
}

func wg2() {
	hello := func(wg *sync.WaitGroup, id int) {
		defer wg.Done()
		fmt.Printf("Hello from %v\n", id)
	}

	const numGreeters = 5
	var wg sync.WaitGroup
	wg.Add(numGreeters)
	for i := 0; i < numGreeters; i++ {
		go hello(&wg, i+1)
	}
	wg.Wait()
}
