package main

import (
	"fmt"
	"sync"
)

func main() {
	fmt.Println("--- Running once1 ---")
	once1()

	fmt.Println("--- Running once2 ---")
	once2()

	// fmt.Println("--- Running once3 ---")
	// once3()
}

func once1() {
	var count int
	increment := func() {
		count++
	}

	var once sync.Once
	var increments sync.WaitGroup
	increments.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer increments.Done()
			once.Do(increment)
		}()
	}

	increments.Wait()
	fmt.Printf("Count is %d\n", count)
}

func once2() {
	var count int
	increment := func() { count++ }
	decrement := func() { count-- }

	var once sync.Once
	once.Do(increment)
	once.Do(decrement)

	fmt.Printf("Count is %d\n", count)
}

// Deadlock!!!
func once3() {
	var onceA, onceB sync.Once
	var initB func()
	initA := func() { onceB.Do(initB) }
	initB = func() { onceA.Do(initA) }
	onceA.Do(initA)
}
