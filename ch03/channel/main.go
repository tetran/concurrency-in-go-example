package main

import (
	"bytes"
	"fmt"
	"os"
	"sync"
)

func main() {
	fmt.Println("--- Running chan1 ---")
	chan1()

	fmt.Println("--- Running chan2 ---")
	chan2()

	fmt.Println("--- Running chan3 ---")
	chan3()

	fmt.Println("--- Running chan4 ---")
	chan4()

	fmt.Println("--- Running chan5 ---")
	chan5()

	fmt.Println("--- Running chan6 ---")
	chan6()

	fmt.Println("--- Running chan7 ---")
	chan7()

	// fmt.Println("--- Running chanErr2 ---")
	// chanErr2()
}

func chan1() {
	stringChan := make(chan string)
	go func() {
		stringChan <- "Hello channels!"
	}()
	fmt.Println(<-stringChan)
}

// channel return two values
func chan2() {
	stringChan := make(chan string)
	go func() {
		stringChan <- "Hello channels!"
	}()
	salutation, ok := <-stringChan
	fmt.Printf("(%v): %v\n", ok, salutation)
}

// read from closed channel
func chan3() {
	intChan := make(chan int)
	close(intChan)
	integer, ok := <-intChan
	fmt.Printf("(%v): %v\n", ok, integer)
}

// read with range
func chan4() {
	intChan := make(chan int)

	go func() {
		defer close(intChan)
		for i := 1; i <= 5; i++ {
			intChan <- i
		}
	}()

	for val := range intChan {
		fmt.Printf("%v ", val)
	}
	fmt.Println()
}

// release multiple goroutines at once
func chan5() {
	begin := make(chan any)
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			v, ok := <-begin
			fmt.Printf("(%v): %v\n", ok, v)
			fmt.Printf("%v has begun\n", i)
		}(i)
	}

	fmt.Println("Unblocking goroutines...")
	close(begin)
	wg.Wait()
}

// buffered channel
func chan6() {
	var stdoutBuff bytes.Buffer
	defer stdoutBuff.WriteTo(os.Stdout)

	intChan := make(chan int, 4)
	go func() {
		defer close(intChan)
		defer fmt.Fprintln(&stdoutBuff, "Producer Done.")
		for i := 0; i < 5; i++ {
			fmt.Fprintf(&stdoutBuff, "Sending: %d\n", i)
			intChan <- i
		}
	}()

	for val := range intChan {
		fmt.Fprintf(&stdoutBuff, "Received %v.\n", val)
	}
}

// channel ownership
func chan7() {
	chanOwner := func() <-chan int {
		resultChan := make(chan int, 5)
		go func() {
			defer close(resultChan)
			for i := 0; i <= 5; i++ {
				resultChan <- i
			}
		}()
		return resultChan
	}

	resultChan := chanOwner()
	for result := range resultChan {
		fmt.Printf("Received: %d\n", result)
	}
	fmt.Println("Done receiving!")
}

// // Compile error
// func chanErr1() {
// 	writeChan := make(chan<- any)
// 	readChan := make(<-chan any)
// 	// invalid operation: cannot receive from send-only channel writeChan (variable of type chan<- any)
// 	<-writeChan
// 	// invalid operation: cannot send to receive-only channel readChan (variable of type <-chan any)
// 	readChan <- struct{}{}
// }

// // Deadlock
// // fatal error: all goroutines are asleep - deadlock!
// func chanErr2() {
// 	stringChan := make(chan string)
// 	go func() {
// 		if 0 != 1 {
// 			return
// 		}
// 		stringChan <- "Hello chalnels!"
// 	}()
// 	fmt.Println(<-stringChan)
// }
