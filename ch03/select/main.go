package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("--- Running select1 ---")
	select1()

	fmt.Println("--- Running select2 ---")
	select2()

	fmt.Println("--- Running select3 ---")
	select3()

	fmt.Println("--- Running select4 ---")
	select4()

	fmt.Println("--- Running select5 ---")
	select5()
}

func select1() {
	start := time.Now()
	c := make(chan any)
	go func() {
		time.Sleep(1 * time.Second)
		close(c)
	}()

	fmt.Println("Blocking on read...")
	select {
	case <-c:
		fmt.Printf("Unblocked %v later.\n", time.Since(start))
	}
}

func select2() {
	c1 := make(chan any)
	close(c1)
	c2 := make(chan any)
	close(c2)

	var c1Count, c2Count int
	for i := 1000; i >= 0; i-- {
		select {
		case <-c1:
			c1Count++
		case <-c2:
			c2Count++
		}
	}

	fmt.Printf("c1Count: %d\nc2Count: %d\n", c1Count, c2Count)
}

func select3() {
	var c <-chan int
	select {
	case <-c:
	case <-time.After(1 * time.Second):
		fmt.Println("Timed out.")
	}
}

func select4() {
	start := time.Now()
	var c1, c2 <-chan int
	select {
	case <-c1:
	case <-c2:
	default:
		fmt.Printf("In default after %v\n\n", time.Since(start))
	}
}

func select5() {
	done := make(chan any)
	go func() {
		time.Sleep(5 * time.Second)
		close(done)
	}()

	workCounter := 0
loop:
	for {
		select {
		case <-done:
			break loop
		default:
		}

		workCounter++
		time.Sleep(1 * time.Second)
	}

	fmt.Printf("Achieved %v cycles of work before signalled to stop.\n", workCounter)
}
