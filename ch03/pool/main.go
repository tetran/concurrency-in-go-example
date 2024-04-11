package main

import (
	"fmt"
	"sync"
)

func main() {
	fmt.Println("--- Running pool1 ---")
	pool1()

	fmt.Println("--- Running pool2 ---")
	pool2()
}

func pool1() {
	myPool := &sync.Pool{
		New: func() any {
			fmt.Println("Creating new instance.")
			return struct{}{}
		},
	}

	myPool.Get()
	instance := myPool.Get()
	myPool.Put(instance)
	myPool.Get()
}

func pool2() {
	var numCalcsCreated int
	calcPool := &sync.Pool{
		New: func() any {
			numCalcsCreated += 1
			mem := make([]byte, 1024)
			return &mem
		},
	}

	// Keep 4KB in the pool
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())

	const numWorkers = 1024 * 1024
	var wg sync.WaitGroup
	wg.Add(numWorkers)
	for i := numWorkers; i > 0; i-- {
		go func() {
			defer wg.Done()

			mem := calcPool.Get().(*[]byte)
			defer calcPool.Put(mem)

			// Do something interesting
		}()
	}

	wg.Wait()
	fmt.Printf("%d calculators were created.", numCalcsCreated)
}
