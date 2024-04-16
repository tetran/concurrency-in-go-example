package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("--- Running send ---")
	send()

	fmt.Println("--- Running stop ---")
	stop()
}

func send() {
	done := make(chan any)
	go func() {
		time.Sleep(1 * time.Millisecond)
		close(done)
	}()

	intChan := make(chan int)
	go func() {
		for i := range intChan {
			fmt.Println(i)
		}
	}()

	for i := 0; ; i++ {
		select {
		case <-done:
			fmt.Println("Done")
			return
		case intChan <- i:
		}
	}
}

func stop() {
	done := make(chan any)
	go func() {
		time.Sleep(1 * time.Millisecond)
		close(done)
	}()

	for i := 0; ; i++ {
		select {
		case <-done:
			fmt.Printf("\nDone: %d\n", i)
			return
		default:
			fmt.Print(".")
		}
	}
}
